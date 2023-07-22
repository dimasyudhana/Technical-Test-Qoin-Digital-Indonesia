package main

import (
	"github.com/dimasyudhana/Qoin-Digital-Indonesia/app/config"
	"github.com/dimasyudhana/Qoin-Digital-Indonesia/app/database"
	"github.com/dimasyudhana/Qoin-Digital-Indonesia/app/router"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	cfg := config.InitConfig()
	db := database.InitDatabase(cfg)
	database.InitMigration(db)
	router.InitRouter(db, e)
	e.Logger.Fatal(e.Start(":8181"))
}
