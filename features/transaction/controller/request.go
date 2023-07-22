package controller

import (
	"strconv"

	"github.com/dimasyudhana/Qoin-Digital-Indonesia/features/transaction"
)

type TransactionRequest struct {
	RestaurantID  string `json:"restaurant_id"`
	PaymentMethod string `json:"payment_method"`
	PaymentType   string `json:"payment_type"`
}

type ProductTransactionRequest struct {
	ProductProductID string `json:"product_id"`
	TransactionID    string `json:"transaction_id"`
	Subtotal         string `json:"subtotal"`
	Quantity         string `json:"quantity"`
}

func (tr *TransactionRequest) Carts() transaction.TransactionCore {
	return transaction.TransactionCore{
		RestaurantID:  tr.RestaurantID,
		PaymentMethod: tr.PaymentMethod,
		PaymentType:   tr.PaymentType,
	}
}

func (ptr *ProductTransactionRequest) Carts() transaction.Product_TransactionsCore {
	subtotal, err := strconv.ParseFloat(ptr.Subtotal, 64)
	if err != nil {
		log.Error("error while parsing grandtotal to float64")
		return transaction.Product_TransactionsCore{}
	}

	quantity, err := strconv.ParseFloat(ptr.Quantity, 64)
	if err != nil {
		log.Error("error while parsing quantity to float64")
		return transaction.Product_TransactionsCore{}
	}

	return transaction.Product_TransactionsCore{
		ProductProductID: ptr.ProductProductID,
		TransactionID:    ptr.TransactionID,
		Subtotal:         subtotal,
		Quantity:         quantity,
	}
}
