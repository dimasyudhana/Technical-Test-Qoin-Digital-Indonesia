package product

import (
	"time"

	"github.com/labstack/echo/v4"
)

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
	DeletedAt            time.Time
	User                 UserCore
	Restaurant           RestaurantCore
	Product_Transactions []Product_TransactionCore
}

type Product_TransactionCore struct {
	ProductTransactionID string
	ProductProductID     string
	TransactionID        string
	Subtotal             float64
	Quantity             float64
	Stock                float64
	CreatedAt            time.Time
	UpdatedAt            time.Time
	DeletedAt            time.Time
}

type StockCore struct {
	RestaurantName  string
	ProductName     string
	ProductQuantity float64
}

type Controller interface {
	RegisterRestaurantAndProducts() echo.HandlerFunc
	Stocks() echo.HandlerFunc
}

type UseCase interface {
	RegisterRestaurantAndProducts(userId string, request RestaurantCore) (RestaurantCore, error)
	Stocks(userId string, productId string) (StockCore, error)
}

type Repository interface {
	RegisterRestaurantAndProducts(userId string, request RestaurantCore) (RestaurantCore, error)
	Stocks(userId string, productId string) (StockCore, error)
}
