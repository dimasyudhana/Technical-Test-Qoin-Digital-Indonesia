package repository

import (
	"time"

	transaction "github.com/dimasyudhana/Qoin-Digital-Indonesia/features/transaction/repository"
	user "github.com/dimasyudhana/Qoin-Digital-Indonesia/features/user/repository"
)

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

type Restaurant struct {
	RestaurantID      string                    `gorm:"primaryKey;type:varchar(45)"`
	UserID            string                    `gorm:"foreignKey:UserID;type:varchar(45)"`
	RestaurantName    string                    `gorm:"type:text;not null"`
	Description       string                    `gorm:"type:text;not null"`
	Status            string                    `gorm:"type:text;not null"`
	RestaurantProfile string                    `gorm:"type:text;not null"`
	CreatedAt         time.Time                 `gorm:"type:datetime"`
	UpdatedAt         time.Time                 `gorm:"type:datetime"`
	IsDeleted         bool                      `gorm:"type:boolean"`
	User              user.User                 `gorm:"references:UserID"`
	Products          []Product                 `gorm:"foreignKey:RestaurantID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	Transactions      []transaction.Transaction `gorm:"foreignKey:RestaurantID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}
