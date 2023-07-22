package controller

import (
	"fmt"
	"time"

	"github.com/dimasyudhana/Qoin-Digital-Indonesia/features/transaction"
)

type invoiceResponse struct {
	RestaurantName    string     `json:"restaurant_name,omitempty"`
	Owner             string     `json:"owner"`
	Buyer             string     `json:"buyer"`
	Invoice           string     `json:"invoice,omitempty"`
	Grandtotal        string     `json:"grandtotal,omitempty"`
	PurchaseStartDate string     `json:"purchase_start_date,omitempty"`
	PurchaseEndDate   string     `json:"purchase_end_date,omitempty"`
	PaymentStatus     string     `json:"payment_status,omitempty"`
	PaymentMethod     string     `json:"payment_method,omitempty"`
	PaymentType       string     `json:"payment_type,omitempty"`
	PaymentCode       string     `json:"payment_code,omitempty"`
	Products          []Products `json:"products,omitempty"`
}

type Products struct {
	ProductTransactionID string `json:"product_transaction_id,omitempty"`
	ProductProductID     string `json:"product_id,omitempty"`
	ProductName          string `json:"product_name,omitempty"`
	Subtotal             string `json:"subtotal,omitempty"`
	Quantity             string `json:"quantity,omitempty"`
}

type earningsResponse struct {
	Username string `json:"username,omitempty"`
	Earnings string `json:"earnings,omitempty"`
}

func invoice(transactions []transaction.InvoiceCore) map[string]*invoiceResponse {
	invoice := make(map[string]*invoiceResponse)

	for _, t := range transactions {
		key := t.TransactionID
		if _, ok := invoice[key]; !ok {
			invoice[key] = &invoiceResponse{
				RestaurantName:    t.RestaurantName,
				Owner:             t.Owner,
				Buyer:             t.Buyer,
				Invoice:           t.Invoice,
				Grandtotal:        fmt.Sprintf("%.2f", t.Grandtotal),
				PurchaseStartDate: parseTimeToString(t.PurchaseStartDate),
				PurchaseEndDate:   parseTimeToString(t.PurchaseEndDate),
				PaymentStatus:     t.PaymentStatus,
				PaymentMethod:     t.PaymentMethod,
				PaymentType:       t.PaymentType,
				PaymentCode:       t.PaymentCode,
				Products:          make([]Products, 0), // Initialize the products array
			}
		}

		// Add transaction data to the Products array
		invoice[key].Products = append(invoice[key].Products, Products{
			ProductTransactionID: t.ProductTransactionID,
			ProductProductID:     t.ProductProductID,
			ProductName:          t.ProductName,
			Subtotal:             fmt.Sprintf("%.2f", t.Subtotal),
			Quantity:             fmt.Sprintf("%.2f", t.Quantity),
		})
	}

	return invoice
}

func earnings(r transaction.EarningsCore) earningsResponse {
	response := earningsResponse{
		Username: r.Username,
		Earnings: fmt.Sprintf("%.2f", r.Earnings),
	}

	return response
}

func parseTimeToString(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}
