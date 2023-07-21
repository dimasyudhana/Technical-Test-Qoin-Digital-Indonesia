package usecase

import (
	"errors"
	"strings"

	"github.com/dimasyudhana/Qoin-Digital-Indonesia/app/middlewares"
	"github.com/dimasyudhana/Qoin-Digital-Indonesia/features/product"
)

var log = middlewares.Log()

type Service struct {
	query product.Repository
}

func New(pd product.Repository) product.UseCase {
	return &Service{
		query: pd,
	}
}

// RegisterRestaurantAndProducts implements product.UseCase.
func (ps *Service) RegisterRestaurantAndProducts(userId string, request product.RestaurantCore) (product.RestaurantCore, error) {
	if request.RestaurantName == "" {
		return product.RestaurantCore{}, errors.New("restaurant name cannot be empty")
	}

	if len(request.Products) == 0 {
		return product.RestaurantCore{}, errors.New("at least one product is required")
	}

	result, err := ps.query.RegisterRestaurantAndProducts(userId, request)
	if err != nil {
		if strings.Contains(err.Error(), "products record not found") {
			log.Error("products record not found")
			return product.RestaurantCore{}, errors.New("products record not found")
		} else {
			log.Error("internal server error")
			return product.RestaurantCore{}, errors.New("internal server error")
		}
	}

	return result, err
}
