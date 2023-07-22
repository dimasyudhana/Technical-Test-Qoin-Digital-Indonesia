package repository

import (
	"time"

	"github.com/dimasyudhana/Qoin-Digital-Indonesia/features/product"
	transaction "github.com/dimasyudhana/Qoin-Digital-Indonesia/features/transaction/repository"
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
	User              User                      `gorm:"foreignKey:UserID"`
	Products          []Product                 `gorm:"foreignKey:RestaurantID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	Transactions      []transaction.Transaction `gorm:"foreignKey:RestaurantID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

type User struct {
	UserID         string       `gorm:"primaryKey;type:varchar(45)"`
	Username       string       `gorm:"type:varchar(225);not null"`
	Email          string       `gorm:"type:varchar(225);not null;unique"`
	Password       string       `gorm:"type:text;not null"`
	Role           string       `gorm:"type:enum('user', 'owner');default:'user'"`
	Status         string       `gorm:"type:enum('verified', 'unverified');default:'unverified'"`
	ProfilePicture string       `gorm:"type:text"`
	CreatedAt      time.Time    `gorm:"type:datetime"`
	UpdatedAt      time.Time    `gorm:"type:datetime"`
	IsDeleted      bool         `gorm:"type:boolean"`
	Restaurant     []Restaurant `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

func productEntities(core product.ProductCore) Product {
	return Product{
		ProductID:       core.ProductID,
		RestaurantID:    core.RestaurantID,
		ProductName:     core.ProductName,
		Description:     core.Description,
		ProductImage:    core.ProductImage,
		ProductCategory: core.ProductCategory,
		ProductPrice:    core.ProductPrice,
		ProductQuantity: core.ProductQuantity,
		CreatedAt:       core.CreatedAt,
		UpdatedAt:       core.UpdatedAt,
		IsDeleted:       core.IsDeleted,
	}
}

func restaurantEntities(core product.RestaurantCore) Restaurant {
	products := make([]Product, len(core.Products))

	for i, p := range core.Products {
		products[i] = productEntities(p)
	}

	return Restaurant{
		RestaurantID:      core.RestaurantID,
		UserID:            core.UserID,
		RestaurantName:    core.RestaurantName,
		Description:       core.Description,
		Status:            core.Status,
		RestaurantProfile: core.RestaurantProfile,
		CreatedAt:         core.CreatedAt,
		UpdatedAt:         core.UpdatedAt,
		IsDeleted:         core.IsDeleted,
		Products:          products,
	}
}

func productModels(model Product) product.ProductCore {
	return product.ProductCore{
		ProductID:       model.ProductID,
		RestaurantID:    model.RestaurantID,
		ProductName:     model.ProductName,
		Description:     model.Description,
		ProductImage:    model.ProductImage,
		ProductCategory: model.ProductCategory,
		ProductPrice:    model.ProductPrice,
		ProductQuantity: model.ProductQuantity,
		CreatedAt:       model.CreatedAt,
		UpdatedAt:       model.UpdatedAt,
		IsDeleted:       model.IsDeleted,
	}
}

func restaurantModels(model Restaurant) product.RestaurantCore {
	products := make([]product.ProductCore, len(model.Products))

	for i, p := range model.Products {
		products[i] = productModels(p)
	}

	return product.RestaurantCore{
		RestaurantID:      model.RestaurantID,
		UserID:            model.UserID,
		RestaurantName:    model.RestaurantName,
		Description:       model.Description,
		Status:            model.Status,
		RestaurantProfile: model.RestaurantProfile,
		CreatedAt:         model.CreatedAt,
		UpdatedAt:         model.UpdatedAt,
		IsDeleted:         model.IsDeleted,
		Products:          products,
	}
}
