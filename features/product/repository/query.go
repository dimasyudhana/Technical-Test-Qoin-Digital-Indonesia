package repository

import (
	"errors"

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

	// Generate a new restaurant ID
	restaurantId, err := identity.GenerateID()
	if err != nil {
		tx.Rollback()
		log.Error("error while creating id for restaurant")
		return restaurant, errors.New("error while creating id for restaurant")
	}

	// Create the restaurant record
	restaurant = product.RestaurantCore{
		RestaurantID:   restaurantId,
		UserID:         userId,
		RestaurantName: request.RestaurantName,
	}
	if err := tx.Create(&restaurant).Error; err != nil {
		tx.Rollback()
		log.Error("failed to create restaurant")
		return restaurant, err
	}

	// Create the products associated with the restaurant
	for i := range request.Products {
		request.Products[i].RestaurantID = restaurant.RestaurantID
		if err := tx.Create(&request.Products[i]).Error; err != nil {
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
