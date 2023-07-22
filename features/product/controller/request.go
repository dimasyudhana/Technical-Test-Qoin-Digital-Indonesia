package controller

import (
	"strconv"

	"github.com/dimasyudhana/Qoin-Digital-Indonesia/features/product"
)

type RegisterRestaurantRequest struct {
	RestaurantName string                   `json:"restaurant_name" form:"restaurant_name"`
	Products       []RegisterProductRequest `json:"products" form:"products"`
}

type RegisterProductRequest struct {
	ProductName     string `json:"product_name" form:"product_name"`
	Description     string `json:"description" form:"description"`
	ProductImage    string `json:"product_image" form:"product_image"`
	ProductCategory string `json:"product_category" form:"product_category"`
	ProductPrice    string `json:"product_price" form:"product_price"`
	ProductQuantity string `json:"product_quantity" form:"product_quantity"`
}

func (r RegisterProductRequest) registerProduct() product.ProductCore {
	productPrice, err := strconv.ParseFloat(r.ProductPrice, 64)
	if err != nil {
		log.Error("error while parsing price to float64")
		return product.ProductCore{}
	}
	productQuantity, err := strconv.ParseFloat(r.ProductQuantity, 64)
	if err != nil {
		log.Error("error while parsing quantity to float64")
		return product.ProductCore{}
	}
	return product.ProductCore{
		ProductName:     r.ProductName,
		Description:     r.Description,
		ProductImage:    r.ProductImage,
		ProductCategory: r.ProductCategory,
		ProductPrice:    productPrice,
		ProductQuantity: productQuantity,
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
