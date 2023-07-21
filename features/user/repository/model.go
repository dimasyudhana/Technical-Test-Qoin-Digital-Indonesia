package repository

import (
	"time"

	"github.com/dimasyudhana/Qoin-Digital-Indonesia/features/user"
	"gorm.io/gorm"
)

type User struct {
	UserID         string        `gorm:"primaryKey;type:varchar(45)"`
	Username       string        `gorm:"type:varchar(225);not null"`
	Email          string        `gorm:"type:varchar(225);not null;unique"`
	Password       string        `gorm:"type:text;not null"`
	Role           string        `gorm:"type:enum('user', 'owner');default:'user'"`
	Status         string        `gorm:"type:enum('verified', 'unverified');default:'unverified'"`
	ProfilePicture string        `gorm:"type:text"`
	CreatedAt      time.Time     `gorm:"type:datetime"`
	UpdatedAt      time.Time     `gorm:"type:datetime"`
	IsDeleted      bool          `gorm:"type:boolean"`
	Restaurant     []Restaurant  `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	Transactions   []Transaction `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

type Restaurant struct {
	RestaurantID      string        `gorm:"primaryKey;type:varchar(45)"`
	UserID            string        `gorm:"foreignKey:UserID;type:varchar(45)"`
	RestaurantName    string        `gorm:"type:text;not null"`
	Description       string        `gorm:"type:text;not null"`
	Status            string        `gorm:"type:text;not null"`
	RestaurantProfile string        `gorm:"type:text;not null"`
	CreatedAt         time.Time     `gorm:"type:datetime"`
	UpdatedAt         time.Time     `gorm:"type:datetime"`
	IsDeleted         bool          `gorm:"type:boolean"`
	User              User          `gorm:"references:UserID"`
	Products          []Product     `gorm:"foreignKey:RestaurantID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	Transactions      []Transaction `gorm:"foreignKey:RestaurantID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

type Product struct {
	ProductID       string     `gorm:"primaryKey;type:varchar(45)"`
	RestaurantID    string     `gorm:"foreignKey:RestaurantID;type:varchar(45)"`
	ProductName     string     `gorm:"type:text;not null"`
	Description     string     `gorm:"type:text;not null"`
	ProductImage    string     `gorm:"type:text;not null"`
	ProductCategory string     `gorm:"type:text;not null"`
	ProductPrice    float64    `gorm:"type:decimal(10,2);"`
	ProductQuantity float64    `gorm:"type:decimal(10,2);"`
	CreatedAt       time.Time  `gorm:"type:datetime"`
	UpdatedAt       time.Time  `gorm:"type:datetime"`
	IsDeleted       bool       `gorm:"type:boolean"`
	Restaurant      Restaurant `gorm:"references:RestaurantID"`
}

type Transaction struct {
	TransactionID        string         `gorm:"primaryKey;type:varchar(45)"`
	RestaurantID         string         `gorm:"foreignKey:RestaurantID;type:varchar(45)"`
	UserID               string         `gorm:"foreignKey:UserID;type:varchar(45)"`
	Invoice              string         `gorm:"type:varchar(45);not null"`
	Grandtotal           float64        `gorm:"type:decimal(10,2);"`
	PaymentStatus        string         `gorm:"type:enum('pending','success','cancel','expire');default:'pending'"`
	PaymentMethod        string         `gorm:"type:text;not null"`
	PaymentType          string         `gorm:"type:text;not null"`
	PaymentCode          string         `gorm:"type:text;not null"`
	PurchaseStartDate    time.Time      `gorm:"type:datetime"`
	PurchaseEndDate      time.Time      `gorm:"type:datetime"`
	CreatedAt            time.Time      `gorm:"type:datetime"`
	UpdatedAt            time.Time      `gorm:"type:datetime"`
	DeletedAt            gorm.DeletedAt `gorm:"index"`
	User                 User           `gorm:"references:UserID"`
	Restaurant           Restaurant     `gorm:"references:RestaurantID"`
	Product_Transactions []Product      `gorm:"many2many:product_transactions;foreignKey:TransactionID;joinForeignKey:TransactionID"`
}

// User-model to user-core
func userModels(u User) user.UserCore {
	return user.UserCore{
		UserID:         u.UserID,
		Username:       u.Username,
		Email:          u.Email,
		Password:       u.Password,
		Role:           u.Role,
		Status:         u.Status,
		ProfilePicture: u.ProfilePicture,
		CreatedAt:      u.CreatedAt,
		UpdatedAt:      u.UpdatedAt,
		IsDeleted:      u.IsDeleted,
	}
}

// User-core to user-model
func userEntities(u user.UserCore) User {
	return User{
		UserID:         u.UserID,
		Username:       u.Username,
		Email:          u.Email,
		Password:       u.Password,
		Role:           u.Role,
		Status:         u.Status,
		ProfilePicture: u.ProfilePicture,
		CreatedAt:      u.CreatedAt,
		UpdatedAt:      u.UpdatedAt,
		IsDeleted:      u.IsDeleted,
	}
}
