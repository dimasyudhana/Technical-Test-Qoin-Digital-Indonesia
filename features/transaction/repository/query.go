package repository

import (
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
