package repository

import (
	"time"

	"gorm.io/gorm"
)

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

type User struct {
	UserID         string        `gorm:"primaryKey;type:varchar(45)"`
	Username       string        `gorm:"type:varchar(225);not null"`
	Email          string        `gorm:"type:varchar(225);not null;unique"`
	Password       string        `gorm:"type:text;not null"`
	Role           string        `gorm:"type:enum('user', 'owner');default:'user'"`
	Status         string        `gorm:"type:enum('verified', 'unverified');default:'unverified'"`
	ProfilePicture string        `gorm:"type:varchar(255);default:'https://cdn.pixabay.com/photo/2015/10/05/22/37/blank-profile-picture-973460_1280.png'"`
	CreatedAt      time.Time     `gorm:"type:datetime"`
	UpdatedAt      time.Time     `gorm:"type:datetime"`
	IsDeleted      bool          `gorm:"type:boolean"`
	Restaurant     []Restaurant  `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	Transactions   []Transaction `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
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
