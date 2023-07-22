package controller

import (
	"net/http"
	"strings"

	"github.com/dimasyudhana/Qoin-Digital-Indonesia/app/middlewares"
	"github.com/dimasyudhana/Qoin-Digital-Indonesia/features/product"
	"github.com/dimasyudhana/Qoin-Digital-Indonesia/utils/response"
	echo "github.com/labstack/echo/v4"
)

var log = middlewares.Log()

type Controller struct {
	service product.UseCase
}

func New(us product.UseCase) product.Controller {
	return &Controller{
		service: us,
	}
}

// RegisterProduct implements product.Controller.
func (rh *Controller) RegisterRestaurantAndProducts() echo.HandlerFunc {
	return func(c echo.Context) error {
		req := struct {
			Restaurant RegisterRestaurantRequest `json:"restaurant"`
			Products   []RegisterProductRequest  `json:"products"`
		}{}

		userId, errToken := middlewares.ExtractToken(c)
		if errToken != nil {
			log.Error("missing or malformed JWT")
			return response.UnauthorizedError(c, "Missing or malformed JWT")
		}

		errBind := c.Bind(&req)
		if errBind != nil {
			log.Error("error on bind request")
			return response.BadRequestError(c, "Bad request")
		}

		restaurantCore := req.Restaurant.registerRestaurant()
		for _, p := range req.Products {
			productCore := p.registerProduct()
			restaurantCore.Products = append(restaurantCore.Products, productCore)
		}

		_, err := rh.service.RegisterRestaurantAndProducts(userId, restaurantCore)
		if err != nil {
			switch {
			case strings.Contains(err.Error(), "empty"):
				log.Error("bad request, request cannot be empty")
				return response.BadRequestError(c, "Bad request")
			case strings.Contains(err.Error(), "unregistered user"):
				log.Error("unregistered user")
				return response.BadRequestError(c, "Bad request")
			default:
				log.Error("internal server error")
				return response.InternalServerError(c, "Internal server error")
			}
		}
		return c.JSON(http.StatusCreated, response.ResponseFormat(http.StatusCreated, "Successfully operation", nil, nil))
	}
}
