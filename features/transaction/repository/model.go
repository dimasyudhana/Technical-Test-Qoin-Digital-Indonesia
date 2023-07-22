package repository

import (
	"time"

	"github.com/dimasyudhana/Qoin-Digital-Indonesia/features/transaction"
	"github.com/dimasyudhana/Qoin-Digital-Indonesia/utils/identity"
	"gorm.io/gorm"
)

type Transaction struct {
	TransactionID     string         `gorm:"primaryKey;type:varchar(45)"`
	RestaurantID      string         `gorm:"foreignKey:RestaurantID;type:varchar(45)"`
	UserID            string         `gorm:"foreignKey:UserID;type:varchar(45)"`
	Invoice           string         `gorm:"type:varchar(45);not null"`
	Grandtotal        float64        `gorm:"type:decimal(10,2);"`
	PaymentStatus     string         `gorm:"type:enum('pending','success','cancel','expire');default:'pending'"`
	PaymentMethod     string         `gorm:"type:text;not null"`
	PaymentType       string         `gorm:"type:text;not null"`
	PaymentCode       string         `gorm:"type:text;not null"`
	PurchaseStartDate time.Time      `gorm:"type:datetime"`
	PurchaseEndDate   time.Time      `gorm:"type:datetime"`
	CreatedAt         time.Time      `gorm:"type:datetime"`
	UpdatedAt         time.Time      `gorm:"type:datetime"`
	DeletedAt         gorm.DeletedAt `gorm:"index"`
	User              User           `gorm:"references:UserID"`
	Restaurant        Restaurant     `gorm:"references:RestaurantID"`
	Products          []Product      `gorm:"many2many:product_transactions;foreignKey:TransactionID;joinForeignKey:TransactionID"`
}

type Product struct {
	ProductID       string     `gorm:"primaryKey;type:varchar(45)"`
	RestaurantID    string     `gorm:"foreignKey:RestaurantID;type:varchar(45)"`
	ProductName     string     `gorm:"type:text;not null"`
	Description     string     `gorm:"type:text;not null"`
	ProductImage    string     `gorm:"type:text;not null"`
	ProductCategory string     `gorm:"type:text;not null"`
	ProductPrice    float64    `gorm:"type:decimal(10,2);"`
	ProductQuantity float64    `gorm:"type:decimal(10,2);"`
	CreatedAt       time.Time  `gorm:"type:datetime"`
	UpdatedAt       time.Time  `gorm:"type:datetime"`
	IsDeleted       bool       `gorm:"type:boolean"`
	Restaurant      Restaurant `gorm:"references:RestaurantID"`
}

type Restaurant struct {
	RestaurantID      string        `gorm:"primaryKey;type:varchar(45)"`
	UserID            string        `gorm:"foreignKey:UserID;type:varchar(45)"`
	RestaurantName    string        `gorm:"type:text;not null"`
	Description       string        `gorm:"type:text;not null"`
	Status            string        `gorm:"type:text;not null"`
	RestaurantProfile string        `gorm:"type:text;not null"`
	CreatedAt         time.Time     `gorm:"type:datetime"`
	UpdatedAt         time.Time     `gorm:"type:datetime"`
	IsDeleted         bool          `gorm:"type:boolean"`
	User              User          `gorm:"references:UserID"`
	Products          []Product     `gorm:"foreignKey:RestaurantID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	Transactions      []Transaction `gorm:"foreignKey:RestaurantID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

type Product_Transactions struct {
	ProductTransactionID string         `gorm:"primaryKey;type:varchar(45)"`
	ProductProductID     string         `gorm:"foreignKey:ProductID;type:varchar(45)"`
	TransactionID        string         `gorm:"foreignKey:TransactionID;type:varchar(45)"`
	Subtotal             float64        `gorm:"type:decimal(10,2);"`
	Quantity             float64        `gorm:"type:decimal(10,2);"`
	Stock                float64        `gorm:"type:decimal(10,2);"`
	CreatedAt            time.Time      `gorm:"type:datetime"`
	UpdatedAt            time.Time      `gorm:"type:datetime"`
	DeletedAt            gorm.DeletedAt `gorm:"index"`
	Product              Product        `gorm:"foreignKey:ProductProductID"`
	Transaction          Transaction    `gorm:"references:TransactionID"`
}

type User struct {
	UserID         string        `gorm:"primaryKey;type:varchar(45)"`
	Username       string        `gorm:"type:varchar(225);not null"`
	Email          string        `gorm:"type:varchar(225);not null;unique"`
	Password       string        `gorm:"type:text;not null"`
	Role           string        `gorm:"type:enum('user', 'owner');default:'user'"`
	Status         string        `gorm:"type:enum('verified', 'unverified');default:'unverified'"`
	ProfilePicture string        `gorm:"type:text"`
	CreatedAt      time.Time     `gorm:"type:datetime"`
	UpdatedAt      time.Time     `gorm:"type:datetime"`
	IsDeleted      bool          `gorm:"type:boolean"`
	Restaurant     []Restaurant  `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	Transactions   []Transaction `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

// Map TransactionCore to Transaction model
func transactionModels(core transaction.TransactionCore) (Transaction, error) {
	transactionID, err := identity.GenerateID()
	if err != nil {
		return Transaction{}, err
	}

	invoice, err := identity.GenerateID()
	if err != nil {
		return Transaction{}, err
	}

	grandtotal := 0.0
	for _, pt := range core.Product_Transactions {
		grandtotal += pt.Subtotal
	}

	payment_code, err := identity.GenerateID()
	if err != nil {
		return Transaction{}, err
	}

	purchaseStartDate := time.Now()
	purchaseEndDate := purchaseStartDate.Add(24 * time.Hour)

	return Transaction{
		TransactionID:     transactionID,
		RestaurantID:      core.RestaurantID,
		UserID:            core.UserID,
		Invoice:           invoice,
		Grandtotal:        grandtotal,
		PaymentStatus:     core.PaymentStatus,
		PaymentMethod:     core.PaymentMethod,
		PaymentType:       core.PaymentType,
		PaymentCode:       payment_code,
		PurchaseStartDate: purchaseStartDate,
		PurchaseEndDate:   purchaseEndDate,
		CreatedAt:         core.CreatedAt,
		UpdatedAt:         core.UpdatedAt,
		DeletedAt:         core.DeletedAt,
	}, nil
}

// Map Product_TransactionCore slice to Product_Transactions model slice
func productTransactionsModels(transactionID string, cores ...transaction.Product_TransactionsCore) ([]Product_Transactions, error) {
	models := make([]Product_Transactions, len(cores))
	for i, c := range cores {
		productTransactionID, err := identity.GenerateID()
		if err != nil {
			return nil, err
		}

		models[i] = Product_Transactions{
			ProductTransactionID: productTransactionID,
			ProductProductID:     c.ProductProductID,
			TransactionID:        transactionID,
			Subtotal:             c.Subtotal,
			Quantity:             c.Quantity,
			Stock:                c.Stock,
			CreatedAt:            c.CreatedAt,
			UpdatedAt:            c.UpdatedAt,
			DeletedAt:            c.DeletedAt,
		}
	}
	return models, nil
}

// Map Transaction model to TransactionCore
func transactionEntities(model Transaction) transaction.TransactionCore {
	return transaction.TransactionCore{
		TransactionID:     model.TransactionID,
		RestaurantID:      model.RestaurantID,
		UserID:            model.UserID,
		Invoice:           model.Invoice,
		Grandtotal:        model.Grandtotal,
		PaymentStatus:     model.PaymentStatus,
		PaymentMethod:     model.PaymentMethod,
		PaymentType:       model.PaymentType,
		PaymentCode:       model.PaymentCode,
		PurchaseStartDate: model.PurchaseStartDate,
		PurchaseEndDate:   model.PurchaseEndDate,
		CreatedAt:         model.CreatedAt,
		UpdatedAt:         model.UpdatedAt,
		DeletedAt:         model.DeletedAt,
		// Product_Transactions: productTransactionsEntities(Product_TransactionsCore),
	}
}

// Map Product_Transactions model slice to Product_TransactionCore slice
func productTransactionsEntities(models []Product_Transactions) []transaction.Product_TransactionsCore {
	cores := make([]transaction.Product_TransactionsCore, len(models))
	for i, m := range models {
		cores[i] = transaction.Product_TransactionsCore{
			ProductTransactionID: m.ProductTransactionID,
			ProductProductID:     m.ProductProductID,
			TransactionID:        m.TransactionID,
			Subtotal:             m.Subtotal,
			Quantity:             m.Quantity,
			Stock:                m.Stock,
			CreatedAt:            m.CreatedAt,
			UpdatedAt:            m.UpdatedAt,
			DeletedAt:            m.DeletedAt,
		}
	}
	return cores
}
