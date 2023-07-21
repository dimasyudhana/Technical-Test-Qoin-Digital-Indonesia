package router

import (
	"github.com/dimasyudhana/Qoin-Digital-Indonesia/app/middlewares"
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
}

func initUserRouter(db *gorm.DB, e *echo.Echo) {
	userData := ur.New(db)
	userService := uu.New(userData)
	userHandler := uc.New(userService)

	e.POST("/register", userHandler.Register())
	e.POST("/login", userHandler.Login())
	e.GET("/users", userHandler.Profile(), middlewares.JWTMiddleware())
}
