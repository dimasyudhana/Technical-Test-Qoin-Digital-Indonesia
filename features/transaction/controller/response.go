package controller

import (
	"fmt"

	"github.com/dimasyudhana/Qoin-Digital-Indonesia/features/transaction"
)

type invoiceResponse struct {
	TransactionID     string     `json:"transaction_id,omitempty"`
	Invoice           string     `json:"invoice,omitempty"`
	Grandtotal        string     `json:"grandtotal,omitempty"`
	PurchaseStartDate string     `json:"purchase_start_date,omitempty"`
	PurchaseEndDate   string     `json:"purchase_end_date,omitempty"`
	Products          []Products `json:"products,omitempty"`
}

type Products struct {
	ProductName string `json:"product_name,omitempty"`
}

type earningsResponse struct {
	Username string `json:"username,omitempty"`
	Earnings string `json:"earnings,omitempty"`
}

func invoice(r transaction.Product_TransactionsCore) invoiceResponse {
	response := invoiceResponse{
		TransactionID:     r.TransactionID,
		Invoice:           r.Transaction.Invoice,
		Grandtotal:        fmt.Sprintf("%.2f", r.Transaction.Grandtotal),
		PurchaseStartDate: r.Transaction.PurchaseStartDate.Format("2006-01-02 15:04:05"),
		PurchaseEndDate:   r.Transaction.PurchaseEndDate.Format("2006-01-02 15:04:05"),
		Products: []Products{
			{
				ProductName: r.Product.ProductName,
			},
		},
	}

	return response
}

func earnings(r transaction.EarningsCore) earningsResponse {
	response := earningsResponse{
		Username: r.Username,
		Earnings: fmt.Sprintf("%.2f", r.Earnings),
	}

	return response
}
