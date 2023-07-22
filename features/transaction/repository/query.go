package repository

import (
	"errors"
	"time"

	"github.com/dimasyudhana/Qoin-Digital-Indonesia/app/middlewares"
	"github.com/dimasyudhana/Qoin-Digital-Indonesia/features/transaction"
	"github.com/dimasyudhana/Qoin-Digital-Indonesia/utils/identity"
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

	if transactionModel.PaymentStatus == "" {
		transactionModel.PaymentStatus = "pending"
	}

	// Use raw SQL for inserting the data
	err = tx.Exec(`
		INSERT INTO transactions (
			transaction_id, 
			restaurant_id, 
			user_id, 
			invoice, 
			grandtotal, 
			payment_status, 
			payment_method, 
			payment_type, 
			payment_code, 
			purchase_start_date, 
			purchase_end_date,
			created_at,
			updated_at
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`, transactionModel.TransactionID,
		transactionModel.RestaurantID,
		transactionModel.UserID,
		transactionModel.Invoice,
		transactionModel.Grandtotal,
		transactionModel.PaymentStatus,
		transactionModel.PaymentMethod,
		transactionModel.PaymentType,
		transactionModel.PaymentCode,
		time.Now(),
		time.Now().Add(24*time.Hour),
		time.Now(),
		time.Now()).Error
	if err != nil {
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
	err = tx.Exec(`
		UPDATE transactions
		SET
			restaurant_id = ?,
			user_id = ?,
			invoice = ?,
			grandtotal = ?,
			payment_status = ?,
			payment_method = ?,
			payment_type = ?,
			payment_code = ?,
			purchase_start_date = ?,
			purchase_end_date = ?,
			updated_at = ?
		WHERE
			transactions.deleted_at IS NULL AND transaction_id = ?
		`,
		transactionModel.RestaurantID,
		transactionModel.UserID,
		transactionModel.Invoice,
		transactionModel.Grandtotal,
		transactionModel.PaymentStatus,
		transactionModel.PaymentMethod,
		transactionModel.PaymentType,
		transactionModel.PaymentCode,
		transactionModel.PurchaseStartDate,
		transactionModel.PurchaseEndDate,
		time.Now(), // Set updated_at to the current time
		transactionModel.TransactionID).Error
	if err != nil {
		tx.Rollback()
		log.Error("failed to update grandtotal")
		return transaction.TransactionCore{}, err
	}

	for i := range productTransactionsModel {
		productTransactionID, err := identity.GenerateID()
		if err != nil {
			tx.Rollback()
			log.Error("failed to generate product transaction ID")
			return transaction.TransactionCore{}, err
		}

		err = tx.Exec(`INSERT INTO product_transactions (
			product_transaction_id,
			transaction_id,
			product_product_id,
			quantity,
			subtotal
			) 
			VALUES (?, ?, ?, ?, ?)
			`, productTransactionID,
			productTransactionsModel[i].TransactionID,
			productTransactionsModel[i].ProductProductID,
			productTransactionsModel[i].Quantity,
			productTransactionsModel[i].Subtotal).Error
		if err != nil {
			tx.Rollback()
			log.Error("failed to create product transaction")
			return transaction.TransactionCore{}, err
		}

		productID := productTransactionsModel[i].ProductProductID
		quantity := productTransactionsModel[i].Quantity

		var product Product
		err = tx.Raw("SELECT * FROM products WHERE product_id = ?", productID).Scan(&product).Error
		if err != nil {
			tx.Rollback()
			log.Error("failed to get product from database")
			return transaction.TransactionCore{}, err
		}

		product.ProductQuantity -= quantity
		if product.ProductQuantity == 0.0 {
			tx.Rollback()
			log.Error("unavailable product stock")
			return transaction.TransactionCore{}, err
		}

		err = tx.Exec("UPDATE products SET product_quantity = ? WHERE product_id = ?", product.ProductQuantity, productID).Error
		if err != nil {
			tx.Rollback()
			log.Error("failed to update product stock")
			return transaction.TransactionCore{}, err
		}
	}

	err = tx.Commit().Error
	if err != nil {
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
