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

	err = tx.Commit().Error
	if err != nil {
		log.Error("failed to commit database transaction")
		return restaurant, err
	}

	return restaurant, nil
}

// Stocks implements product.Repository.
func (pq *Query) Stocks(userId string, productId string) (product.StockCore, error) {
	stockCore := Stock{}

	// Use raw SQL to fetch the data
	err := pq.db.Raw(`
		SELECT restaurants.restaurant_name, products.product_name, products.product_quantity
		FROM products
		JOIN restaurants ON products.restaurant_id = restaurants.restaurant_id
		WHERE restaurants.user_id = ? AND products.product_id = ?
	`, userId, productId).Scan(&stockCore).Error
	if err != nil {
		log.Error("failed to get stock data")
		return product.StockCore{}, err
	}

	log.Sugar().Infof("%+v", stockCore)

	return StockModels(stockCore), nil
}

func StockModels(s Stock) product.StockCore {
	return product.StockCore{
		RestaurantName:  s.RestaurantName,
		ProductName:     s.ProductName,
		ProductQuantity: s.ProductQuantity,
	}
}
