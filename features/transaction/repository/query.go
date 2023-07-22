package repository

import (
	"errors"
	"time"

	"github.com/dimasyudhana/Qoin-Digital-Indonesia/app/middlewares"
	"github.com/dimasyudhana/Qoin-Digital-Indonesia/features/transaction"
	"gorm.io/gorm"
)

var log = middlewares.Log()

type Query struct {
	db *gorm.DB
}

func New(db *gorm.DB) transaction.Repository {
	return &Query{
		db: db,
	}
}

// Carts implements transaction.Repository.
func (tq *Query) Carts(userId string, tr transaction.TransactionCore, ptr ...transaction.Product_TransactionsCore) (transaction.TransactionCore, error) {
	tx := tq.db.Begin()
	if tx.Error != nil {
		log.Error("failed to start database transaction")
		return transaction.TransactionCore{}, tx.Error
	}

	tr.UserID = userId
	transactionModel, err := transactionModels(tr)
	if err != nil {
		tx.Rollback()
		log.Error("failed to map transaction")
		return transaction.TransactionCore{}, err
	}

	if err := tx.Create(&transactionModel).Error; err != nil {
		tx.Rollback()
		log.Error("failed to create transaction")
		return transaction.TransactionCore{}, err
	}

	request := transactionEntities(transactionModel)

	productTransactionsModel, err := productTransactionsModels(request.TransactionID, ptr...)
	if err != nil {
		tx.Rollback()
		log.Error("failed to map product transactions")
		return transaction.TransactionCore{}, err
	}

	grandtotal := 0.0
	for _, pt := range productTransactionsModel {
		grandtotal += pt.Subtotal
	}

	transactionModel.Grandtotal = grandtotal
	if err := tx.Save(&transactionModel).Error; err != nil {
		tx.Rollback()
		log.Error("failed to update grandtotal")
		return transaction.TransactionCore{}, err
	}

	for i := range productTransactionsModel {
		if err := tx.Create(&productTransactionsModel[i]).Error; err != nil {
			tx.Rollback()
			log.Error("failed to create product transaction")
			return transaction.TransactionCore{}, err
		}
	}

	if err := tx.Commit().Error; err != nil {
		log.Error("failed to commit database transaction")
		return transaction.TransactionCore{}, err
	}

	return request, nil
}

// Invoice implements transaction.Repository.
func (tq *Query) Invoice(userId string, transactionId string) ([]transaction.InvoiceCore, error) {
	result := []Invoice{}
	err := tq.db.Raw(`SELECT *, products.product_name, buyer.username as buyer, restaurants.restaurant_name, owner.username as owner
		FROM product_transactions
		JOIN transactions ON product_transactions.transaction_id = transactions.transaction_id
		JOIN products ON product_transactions.product_product_id = products.product_id
		JOIN restaurants ON products.restaurant_id = restaurants.restaurant_id
		JOIN users as owner ON restaurants.user_id = owner.user_id
		JOIN users as buyer ON transactions.user_id = buyer.user_id
		WHERE product_transactions.transaction_id = ?
	`, transactionId).Scan(&result).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Error("invoice data not found")
			return []transaction.InvoiceCore{}, errors.New("invoice data not found")
		}
		log.Sugar().Error("error executing invoice query:", err)
		return []transaction.InvoiceCore{}, err
	}

	var invoice []transaction.InvoiceCore
	for _, result := range result {
		invoice = append(invoice, invoiceEntities(result))
	}

	return invoice, nil
}

// Earnings implements transaction.Repository.
func (tq *Query) Earnings(userId string, PurchaseStartDate time.Time, PurchaseEndDate time.Time) (transaction.EarningsCore, error) {
	result := Earnings{}
	query := tq.db.Raw(`
		SELECT transactions.user_id, SUM(transactions.grandtotal) AS earnings
		FROM transactions 
		WHERE transactions.user_id = ?
		AND ((transactions.purchase_start_date BETWEEN ? AND ?) OR (transactions.purchase_end_date BETWEEN ? AND ?))
		AND transactions.payment_status = "success"
		GROUP BY transactions.user_id
	`, userId, PurchaseStartDate, PurchaseEndDate, PurchaseStartDate, PurchaseEndDate).Scan(&result)
	if errors.Is(query.Error, gorm.ErrRecordNotFound) {
		log.Error("earnings record not found")
		return transaction.EarningsCore{}, errors.New("earnings record not found")
	} else if query.Error != nil {
		log.Sugar().Error("error executing earnings query:", query.Error)
		return transaction.EarningsCore{}, query.Error
	} else {
		log.Info("earnings data found in the database")
	}

	return earningsModels(result), nil
}
