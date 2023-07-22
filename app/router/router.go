package router

import (
	"github.com/dimasyudhana/Qoin-Digital-Indonesia/app/middlewares"
	pc "github.com/dimasyudhana/Qoin-Digital-Indonesia/features/product/controller"
	pr "github.com/dimasyudhana/Qoin-Digital-Indonesia/features/product/repository"
	pu "github.com/dimasyudhana/Qoin-Digital-Indonesia/features/product/usecase"
	tc "github.com/dimasyudhana/Qoin-Digital-Indonesia/features/transaction/controller"
	tr "github.com/dimasyudhana/Qoin-Digital-Indonesia/features/transaction/repository"
	tu "github.com/dimasyudhana/Qoin-Digital-Indonesia/features/transaction/usecase"
	uc "github.com/dimasyudhana/Qoin-Digital-Indonesia/features/user/controller"
	ur "github.com/dimasyudhana/Qoin-Digital-Indonesia/features/user/repository"
	uu "github.com/dimasyudhana/Qoin-Digital-Indonesia/features/user/usecase"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/gorm"
)

func InitRouter(db *gorm.DB, e *echo.Echo) {
	e.Use(middleware.CORS())
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `[${time_rfc3339}] ${status} ${method} ${host}${path} ${latency_human}` + "\n",
	}))

	initUserRouter(db, e)
	initProductRouter(db, e)
	initTransactionRouter(db, e)
}

func initUserRouter(db *gorm.DB, e *echo.Echo) {
	userData := ur.New(db)
	userService := uu.New(userData)
	userHandler := uc.New(userService)

	e.POST("/register", userHandler.Register())
	e.POST("/login", userHandler.Login())
	e.GET("/users", userHandler.Profile(), middlewares.JWTMiddleware())
}

func initProductRouter(db *gorm.DB, e *echo.Echo) {
	productData := pr.New(db)
	productService := pu.New(productData)
	productHandler := pc.New(productService)

	e.POST("/products", productHandler.RegisterRestaurantAndProducts(), middlewares.JWTMiddleware())
}

func initTransactionRouter(db *gorm.DB, e *echo.Echo) {
	transactionData := tr.New(db)
	transactionService := tu.New(transactionData)
	transactionHandler := tc.New(transactionService)

	e.POST("/transactions", transactionHandler.Carts(), middlewares.JWTMiddleware())
	e.GET("/transactions/:transaction_id", transactionHandler.Invoice(), middlewares.JWTMiddleware())
}
