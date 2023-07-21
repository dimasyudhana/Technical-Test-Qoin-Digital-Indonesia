package database

import (
	product "github.com/dimasyudhana/Qoin-Digital-Indonesia/features/product/repository"
	transaction "github.com/dimasyudhana/Qoin-Digital-Indonesia/features/transaction/repository"
	user "github.com/dimasyudhana/Qoin-Digital-Indonesia/features/user/repository"
	"gorm.io/gorm"
)

func InitMigration(db *gorm.DB) error {
	err := db.SetupJoinTable(&transaction.Transaction{}, "Product_Transactions", &transaction.Product_Transactions{})
	if err != nil {
		log.Sugar().Error("setup err ", err)
		panic(err.Error())
	}

	err = db.AutoMigrate(
		&user.User{},
		&product.Restaurant{},
		&product.Product{},
		&transaction.Transaction{},
	)

	if err != nil {
		log.Fatal(err.Error())
	}

	log.Info("success migrated to database")

	return err
}
