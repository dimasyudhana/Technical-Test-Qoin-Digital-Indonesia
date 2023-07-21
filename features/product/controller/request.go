package controller

import "github.com/dimasyudhana/Qoin-Digital-Indonesia/features/product"

type RegisterProductRequest struct {
	ProductName string `json:"product_name" form:"product_name"`
}

type RegisterRestaurantRequest struct {
	RestaurantName string                   `json:"restaurant_name" form:"restaurant_name"`
	Products       []RegisterProductRequest `json:"products" form:"products"`
}

func (r RegisterProductRequest) registerProduct() product.ProductCore {
	return product.ProductCore{
		ProductName: r.ProductName,
	}
}

func (r RegisterRestaurantRequest) registerRestaurant() product.RestaurantCore {
	restaurantCore := product.RestaurantCore{
		RestaurantName: r.RestaurantName,
	}

	for _, p := range r.Products {
		restaurantCore.Products = append(restaurantCore.Products, p.registerProduct())
	}

	return restaurantCore
}
