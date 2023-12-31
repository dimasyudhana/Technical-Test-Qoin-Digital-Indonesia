package transaction

import (
	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type TransactionCore struct {
	TransactionID        string
	RestaurantID         string
	UserID               string
	Invoice              string
	Grandtotal           float64
	PaymentStatus        string
	PaymentMethod        string
	PaymentType          string
	PaymentCode          string
	PurchaseStartDate    time.Time
	PurchaseEndDate      time.Time
	CreatedAt            time.Time
	UpdatedAt            time.Time
	DeletedAt            gorm.DeletedAt
	User                 UserCore
	Restaurant           RestaurantCore
	Products             []ProductCore
	Product_Transactions []Product_TransactionsCore
	Earnings             float64
}

type RestaurantCore struct {
	RestaurantID      string
	UserID            string
	RestaurantName    string
	Description       string
	Status            string
	RestaurantProfile string
	CreatedAt         time.Time
	UpdatedAt         time.Time
	IsDeleted         bool
	User              UserCore
	Products          []ProductCore
	Transactions      []TransactionCore
}

type ProductCore struct {
	ProductID       string
	RestaurantID    string
	ProductName     string
	Description     string
	ProductImage    string
	ProductCategory string
	ProductPrice    float64
	ProductQuantity float64
	CreatedAt       time.Time
	UpdatedAt       time.Time
	IsDeleted       bool
	Restaurant      RestaurantCore
}

type Product_TransactionsCore struct {
	ProductTransactionID string
	ProductProductID     string
	TransactionID        string
	Subtotal             float64
	Quantity             float64
	CreatedAt            time.Time
	UpdatedAt            time.Time
	DeletedAt            gorm.DeletedAt
	Product              ProductCore
	Transaction          TransactionCore
}

type UserCore struct {
	UserID         string
	Username       string
	Email          string
	Password       string
	Role           string
	Status         string
	ProfilePicture string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	IsDeleted      bool
	Restaurant     []RestaurantCore
	Transactions   []TransactionCore
}

type InvoiceCore struct {
	ProductTransactionID string
	ProductProductID     string
	TransactionID        string
	RestaurantID         string
	RestaurantName       string
	Owner                string
	UserID               string
	Buyer                string
	ProductName          string
	Subtotal             float64
	Quantity             float64
	Stock                float64
	Invoice              string
	Grandtotal           float64
	PaymentStatus        string
	PaymentMethod        string
	PaymentType          string
	PaymentCode          string
	PurchaseStartDate    time.Time
	PurchaseEndDate      time.Time
}

type EarningsCore struct {
	Username string
	Earnings float64
}

type Controller interface {
	Carts() echo.HandlerFunc
	Invoice() echo.HandlerFunc
	Earnings() echo.HandlerFunc
}

type UseCase interface {
	Carts(userId string, tr TransactionCore, ptr ...Product_TransactionsCore) (TransactionCore, error)
	Invoice(userId string, transactionId string) ([]InvoiceCore, error)
	Earnings(userId string, PurchaseStartDate time.Time, PurchaseEndDate time.Time) (EarningsCore, error)
}

type Repository interface {
	Carts(userId string, tr TransactionCore, ptr ...Product_TransactionsCore) (TransactionCore, error)
	Invoice(userId string, transactionId string) ([]InvoiceCore, error)
	Earnings(userId string, PurchaseStartDate time.Time, PurchaseEndDate time.Time) (EarningsCore, error)
}
