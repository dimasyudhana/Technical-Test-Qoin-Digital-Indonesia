package usecase

import (
	"errors"
	"strings"

	"github.com/dimasyudhana/Qoin-Digital-Indonesia/app/middlewares"
	"github.com/dimasyudhana/Qoin-Digital-Indonesia/features/transaction"
)

var log = middlewares.Log()

type Service struct {
	query transaction.Repository
}

func New(ud transaction.Repository) transaction.UseCase {
	return &Service{
		query: ud,
	}
}

// Carts implements transaction.UseCase.
func (ts *Service) Carts(userId string, tr transaction.TransactionCore, ptr ...transaction.Product_TransactionsCore) (transaction.TransactionCore, error) {
	result, err := ts.query.Carts(userId, tr, ptr...)
	if err != nil {
		if strings.Contains(err.Error(), "products record not found") {
			log.Error("products record not found")
			return transaction.TransactionCore{}, errors.New("products record not found")
		} else {
			log.Error("internal server error")
			return transaction.TransactionCore{}, errors.New("internal server error")
		}
	}

	return result, err
}

// Invoice implements transaction.UseCase.
func (ts *Service) Invoice(userId string, transactionId string) (transaction.Product_TransactionsCore, error) {
	result, err := ts.query.Invoice(userId, transactionId)
	if err != nil {
		if strings.Contains(err.Error(), "transactions record not found") {
			log.Error("transactions record not found")
			return transaction.Product_TransactionsCore{}, errors.New("transactions record not found")
		} else {
			log.Error("internal server error")
			return transaction.Product_TransactionsCore{}, errors.New("internal server error")
		}
	}

	return result, err
}
