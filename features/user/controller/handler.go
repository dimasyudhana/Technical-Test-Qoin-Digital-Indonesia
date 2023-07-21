package controller

import (
	"context"
	"net/http"
	"strings"

	"github.com/dimasyudhana/Qoin-Digital-Indonesia/app/middlewares"
	"github.com/dimasyudhana/Qoin-Digital-Indonesia/features/user"
	"github.com/dimasyudhana/Qoin-Digital-Indonesia/utils/response"
	"github.com/go-playground/mold/modifiers"
	"github.com/labstack/echo/v4"
)

var log = middlewares.Log()

type Controller struct {
	service user.UseCase
}

func New(us user.UseCase) user.Controller {
	return &Controller{
		service: us,
	}
}

// Register implements user.Controller.
func (uc *Controller) Register() echo.HandlerFunc {
	return func(c echo.Context) error {
		request := registerRequest{}
		conform := modifiers.New()

		err := c.Bind(&request)
		if err != nil {
			log.Error("error on bind input")
			return response.BadRequestError(c, "Bad request")
		}

		err = conform.Struct(context.Background(), &request)
		if err != nil {
			log.Error("validation failed")
			return response.BadRequestError(c, "Bad request, validation failed")
		}

		result, err := uc.service.Register(RequestToCore(request))
		if err != nil {
			var message string
			switch {
			case strings.Contains(err.Error(), "request cannot be empty"):
				log.Error("request cannot be empty")
				message = "Bad request, request cannot be empty"
			case strings.Contains(err.Error(), "error insert data, duplicate input"):
				log.Error("error insert data, duplicate input")
				message = "Bad request, duplicate input"
			case strings.Contains(err.Error(), "no row affected"):
				log.Error("no row affected")
				message = "Bad request, duplicate entry"
			case strings.Contains(err.Error(), "error while creating id for user"):
				log.Error("error while creating id for user")
				message = "Internal server error"
			case strings.Contains(err.Error(), "error while hashing password"):
				log.Error("error while hashing password")
				message = "Internal server error"
			default:
				log.Error("internal server error")
				message = "Internal server error"
			}
			return response.BadRequestError(c, message)
		}

		return c.JSON(http.StatusCreated, response.ResponseFormat(http.StatusCreated, "Successfully registered", register(result), nil))
	}
}

// Login implements user.Controller.
func (uc *Controller) Login() echo.HandlerFunc {
	return func(c echo.Context) error {
		request := loginRequest{}
		conform := modifiers.New()

		err := c.Bind(&request)
		if err != nil {
			log.Error("error on bind input")
			return response.BadRequestError(c, "Bad request")
		}

		err = conform.Struct(context.Background(), &request)
		if err != nil {
			log.Error("validation failed")
			return response.BadRequestError(c, "Bad request, validation failed")
		}

		result, token, err := uc.service.Login(RequestToCore(request))
		if err != nil {
			var message string
			switch {
			case strings.Contains(err.Error(), "invalid email format"):
				log.Error("bad request, invalid email format")
				message = "Bad request, invalid email format"
			case strings.Contains(err.Error(), "password cannot be empty"):
				log.Error("bad request, password cannot be empty")
				message = "Bad request, password cannot be empty"
			case strings.Contains(err.Error(), "invalid email and password"):
				log.Error("bad request, invalid email and password")
				message = "Bad request, invalid email and password"
			case strings.Contains(err.Error(), "password does not match"):
				log.Error("bad request, password does not match")
				message = "Bad request, password does not match"
			case strings.Contains(err.Error(), "no row affected"):
				log.Error("no row affected")
				message = "The requested resource was not found"
			case strings.Contains(err.Error(), "error while creating jwt token"):
				log.Error("internal server error, error while creating jwt token")
				message = "Internal server error"
			default:
				log.Error("internal server error")
				message = "Internal server error"
			}
			return response.BadRequestError(c, message)
		}

		return c.JSON(http.StatusOK, response.ResponseFormat(http.StatusOK, "Successful login", loginResponse{
			UserID: result.UserID, Token: token,
		}, nil))
	}
}

// ProfileUser implements user.Controller.
func (uc *Controller) Profile() echo.HandlerFunc {
	return func(c echo.Context) error {
		userId, err := middlewares.ExtractToken(c)
		if err != nil {
			log.Error("missing or malformed JWT")
			return response.UnauthorizedError(c, "Missing or malformed JWT")
		}

		profile, err := uc.service.Profile(userId)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				log.Error("venues not found")
				return response.NotFoundError(c, "The requested resource was not found")
			}
			log.Error("internal server error")
			return response.InternalServerError(c, "Internal server error")
		}

		resp := profileUser(profile)
		return c.JSON(http.StatusOK, response.ResponseFormat(http.StatusOK, "Successfully operation.", resp, nil))
	}
}
