package controller

import (
	"net/http"
	"strings"
	"time"

	"github.com/dimasyudhana/Qoin-Digital-Indonesia/app/middlewares"
	"github.com/dimasyudhana/Qoin-Digital-Indonesia/features/transaction"
	"github.com/dimasyudhana/Qoin-Digital-Indonesia/utils/response"
	echo "github.com/labstack/echo/v4"
)

var log = middlewares.Log()

type Controller struct {
	service transaction.UseCase
}

func New(us transaction.UseCase) transaction.Controller {
	return &Controller{
		service: us,
	}
}

// Carts implements transaction.Controller.
func (tc *Controller) Carts() echo.HandlerFunc {
	return func(c echo.Context) error {
		req := struct {
			Transaction          TransactionRequest          `json:"transaction"`
			Product_Transactions []ProductTransactionRequest `json:"product_transactions"`
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

		transactionCore := req.Transaction.ToCore()
		productTransactionsCore := make([]transaction.Product_TransactionsCore, len(req.Product_Transactions))
		for i, ptr := range req.Product_Transactions {
			productTransactionsCore[i] = ptr.ToCore()
		}

		_, err := tc.service.Carts(userId, transactionCore, productTransactionsCore...)
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

// Invoice implements transaction.Controller.
func (tc *Controller) Invoice() echo.HandlerFunc {
	return func(c echo.Context) error {
		userId, errToken := middlewares.ExtractToken(c)
		if errToken != nil {
			log.Error("missing or malformed JWT")
			return response.UnauthorizedError(c, "Missing or malformed JWT")
		}

		transactionId := c.Param("transaction_id")
		result, err := tc.service.Invoice(userId, transactionId)
		if err != nil {
			if strings.Contains(err.Error(), "invoice record not found") {
				log.Error("invoice record not found")
				return response.NotFoundError(c, "The requested resource was not found")
			} else {
				log.Error("internal server error")
				return response.InternalServerError(c, "Internal server error")
			}
		}

		return c.JSON(http.StatusOK, response.ResponseFormat(http.StatusOK, "Successful Operation", invoice(result), nil))
	}
}

// Earnings implements transaction.Controller.
func (tc *Controller) Earnings() echo.HandlerFunc {
	return func(c echo.Context) error {
		userId, errToken := middlewares.ExtractToken(c)
		if errToken != nil {
			log.Error("missing or malformed JWT")
			return response.UnauthorizedError(c, "Missing or malformed JWT")
		}

		PurchaseStartDateStr := c.QueryParam("start_date")
		PurchaseEndDateStr := c.QueryParam("end_date")
		PurchaseStartDate, err := time.Parse("2006-01-02 15:04:05", PurchaseStartDateStr)
		if err != nil {
			log.Error("failed to parse start_date")
			return response.BadRequestError(c, "Invalid value for start_date")
		}

		log.Sugar().Info(PurchaseStartDate)

		PurchaseEndDate, err := time.Parse("2006-01-02 15:04:05", PurchaseEndDateStr)
		if err != nil {
			log.Error("failed to parse end_date")
			return response.BadRequestError(c, "Invalid value for end_date")
		}

		log.Sugar().Info(PurchaseEndDate)

		result, err := tc.service.Earnings(userId, PurchaseStartDate, PurchaseEndDate)
		if err != nil {
			if strings.Contains(err.Error(), "list reservations record not found") {
				log.Error("list reservations record not found")
				return response.NotFoundError(c, "The requested resource was not found")
			} else {
				log.Error("internal server error")
				return response.InternalServerError(c, "Internal server error")
			}
		}
		return c.JSON(http.StatusOK, response.ResponseFormat(http.StatusOK, "Successful Operation", earnings(result), nil))
	}
}
