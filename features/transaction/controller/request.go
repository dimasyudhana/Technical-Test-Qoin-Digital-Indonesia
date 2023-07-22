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
	ProductProductID string `json:"product_product_id"`
	TransactionID    string `json:"transaction_id"`
	Subtotal         string `json:"subtotal"`
	Quantity         string `json:"quantity"`
}

func parseFloat64(str string) (float64, error) {
	return strconv.ParseFloat(str, 64)
}

func (tr *TransactionRequest) ToCore() transaction.TransactionCore {
	return transaction.TransactionCore{
		RestaurantID:  tr.RestaurantID,
		PaymentMethod: tr.PaymentMethod,
		PaymentType:   tr.PaymentType,
	}
}

func (ptr *ProductTransactionRequest) ToCore() transaction.Product_TransactionsCore {
	subtotal, err := parseFloat64(ptr.Subtotal)
	if err != nil {
		log.Error("error while parsing grandtotal to float64")
		return transaction.Product_TransactionsCore{}
	}

	quantity, err := parseFloat64(ptr.Quantity)
	if err != nil {
		log.Error("error while parsing grandtotal to float64")
		return transaction.Product_TransactionsCore{}
	}

	return transaction.Product_TransactionsCore{
		ProductProductID: ptr.ProductProductID,
		TransactionID:    ptr.TransactionID,
		Subtotal:         subtotal,
		Quantity:         quantity,
	}
}
