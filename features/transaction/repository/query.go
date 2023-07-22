package repository

import (
	"errors"

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

	// log.Sugar().Infof("Transaction: %+v", request)
	// log.Sugar().Infof("Product_Transactions: %+v", request.Product_Transactions)
	return request, nil
}

// Invoice implements transaction.Repository.
func (tq *Query) Invoice(userId string, transactionId string) (transaction.Product_TransactionsCore, error) {
	result := Product_Transactions{}
	err := tq.db.Preload("Product").Preload("Transaction").First(&result, "transaction_id = ?", transactionId).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Error("invoice data not found")
			return transaction.Product_TransactionsCore{}, errors.New("invoice data not found")
		}
		log.Sugar().Error("error executing invoice query:", err)
		return transaction.Product_TransactionsCore{}, err
	}

	response, err := invoiceModels(result)
	if err != nil {
		return transaction.Product_TransactionsCore{}, err
	}

	// log.Sugar().Infof("%+v", result)
	return response, nil
}

func invoiceModels(model Product_Transactions) (transaction.Product_TransactionsCore, error) {
	return transaction.Product_TransactionsCore{
		ProductTransactionID: model.ProductTransactionID,
		ProductProductID:     model.ProductProductID,
		TransactionID:        model.TransactionID,
		Subtotal:             model.Subtotal,
		Quantity:             model.Quantity,
		Product: transaction.ProductCore{
			ProductName: model.Product.ProductName,
		},
		Transaction: transaction.TransactionCore{
			Invoice:           model.Transaction.Invoice,
			Grandtotal:        model.Transaction.Grandtotal,
			PurchaseStartDate: model.Transaction.PurchaseStartDate,
			PurchaseEndDate:   model.Transaction.PurchaseEndDate,
		},
	}, nil
}
