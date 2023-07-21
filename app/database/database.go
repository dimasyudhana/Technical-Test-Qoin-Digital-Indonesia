package database

import (
	"fmt"

	"github.com/dimasyudhana/Qoin-Digital-Indonesia/app/config"
	"github.com/dimasyudhana/Qoin-Digital-Indonesia/app/middlewares"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var log = middlewares.Log()

func InitDatabase(c *config.AppConfig) *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		c.DBUSER, c.DBPASSWORD, c.DBHOST, c.DBPORT, c.DBNAME,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err.Error())
	}

	log.Info("success connected to database")

	return db
}
