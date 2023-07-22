package repository

import (
	"github.com/dimasyudhana/Qoin-Digital-Indonesia/app/middlewares"
	"github.com/dimasyudhana/Qoin-Digital-Indonesia/features/product"
	"github.com/dimasyudhana/Qoin-Digital-Indonesia/utils/identity"
	"gorm.io/gorm"
)

var log = middlewares.Log()

type Query struct {
	db *gorm.DB
}

func New(db *gorm.DB) product.Repository {
	return &Query{
		db: db,
	}
}

// RegisterRestaurantAndProducts implements product.Repository.
func (pq *Query) RegisterRestaurantAndProducts(userId string, request product.RestaurantCore) (restaurant product.RestaurantCore, err error) {
	// Start a transaction
	tx := pq.db.Begin()
	if tx.Error != nil {
		log.Error("failed to start database transaction")
		return restaurant, tx.Error
	}

	request.UserID = userId
	restaurantModel := restaurantEntities(request)
	restaurantModel.RestaurantID, _ = identity.GenerateID()
	if err := tx.Create(&restaurantModel).Error; err != nil {
		tx.Rollback()
		log.Error("failed to create restaurant")
		return restaurant, err
	}

	// Create the products associated with the restaurant
	for i := range request.Products {
		request.Products[i].RestaurantID = restaurantModel.RestaurantID
		request.Products[i].ProductID, _ = identity.GenerateID()
		productModel := productEntities(request.Products[i])
		if err := tx.Create(&productModel).Error; err != nil {
			tx.Rollback()
			log.Error("failed to create product")
			return restaurant, err
		}
	}

	// Commit the transaction
	err = tx.Commit().Error
	if err != nil {
		log.Error("failed to commit database transaction")
		return restaurant, err
	}

	return restaurant, nil
}
