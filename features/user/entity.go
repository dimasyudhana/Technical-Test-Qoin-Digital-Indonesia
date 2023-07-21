package user

import (
	"time"

	"github.com/labstack/echo/v4"
)

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
	Product_Transactions []ProductCore
}

type Controller interface {
	Register() echo.HandlerFunc
	Login() echo.HandlerFunc
	Profile() echo.HandlerFunc
}

type UseCase interface {
	Register(request UserCore) (UserCore, error)
	Login(request UserCore) (UserCore, string, error)
	Profile(userId string) (UserCore, error)
}

type Repository interface {
	Register(request UserCore) (UserCore, error)
	Login(request UserCore) (UserCore, string, error)
	Profile(userId string) (UserCore, error)
}
