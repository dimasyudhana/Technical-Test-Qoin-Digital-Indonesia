package controller

import (
	"fmt"

	"github.com/dimasyudhana/Qoin-Digital-Indonesia/features/product"
)

type RegisterRestaurantAndProductsResponse struct {
	Restaurant RegisteredRestaurant `json:"restaurant"`
	Products   []RegisteredProduct  `json:"products"`
}

type RegisteredRestaurant struct {
	RestaurantID string `json:"restaurant_id"`
}

type RegisteredProduct struct {
	ProductID string `json:"product_id"`
}

func register(result product.RestaurantCore) RegisterRestaurantAndProductsResponse {
	response := RegisterRestaurantAndProductsResponse{
		Restaurant: RegisteredRestaurant{
			RestaurantID: result.RestaurantID,
		},
		Products: make([]RegisteredProduct, len(result.Products)),
	}

	for i, p := range result.Products {
		response.Products[i] = RegisteredProduct{
			ProductID: p.ProductID,
		}
	}

	return response
}

type StockResponse struct {
	RestaurantName  string `json:"restaurant_name"`
	ProductName     string `json:"product_name"`
	ProductQuantity string `json:"stock"`
}

func stock(r product.StockCore) StockResponse {
	response := StockResponse{
		RestaurantName:  r.RestaurantName,
		ProductName:     r.ProductName,
		ProductQuantity: fmt.Sprintf("%.2f", r.ProductQuantity),
	}

	return response
}
