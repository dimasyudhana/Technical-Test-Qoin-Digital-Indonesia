package usecase

import (
	"errors"
	"strings"
	"time"

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
func (ts *Service) Invoice(userId string, transactionId string) ([]transaction.InvoiceCore, error) {
	result, err := ts.query.Invoice(userId, transactionId)
	if err != nil {
		if strings.Contains(err.Error(), "invoice record not found") {
			log.Error("invoice record not found")
			return []transaction.InvoiceCore{}, errors.New("invoice record not found")
		} else {
			log.Error("internal server error")
			return []transaction.InvoiceCore{}, errors.New("internal server error")
		}
	}

	return result, err
}

// Earnings implements transaction.UseCase.
func (ts *Service) Earnings(userId string, PurchaseStartDate time.Time, PurchaseEndDate time.Time) (transaction.EarningsCore, error) {
	result, err := ts.query.Earnings(userId, PurchaseStartDate, PurchaseEndDate)
	if err != nil {
		if strings.Contains(err.Error(), "invoice record not found") {
			log.Error("invoice record not found")
			return transaction.EarningsCore{}, errors.New("invoice record not found")
		} else {
			log.Error("internal server error")
			return transaction.EarningsCore{}, errors.New("internal server error")
		}
	}

	return result, err
}
