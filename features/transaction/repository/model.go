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

type Invoice struct {
	ProductTransactionID string
	ProductProductID     string
	TransactionID        string
	RestaurantID         string
	RestaurantName       string
	Owner                string
	UserID               string
	Buyer                string
	ProductName          string
	Subtotal             float64
	Quantity             float64
	Stock                float64
	Invoice              string
	Grandtotal           float64
	PaymentStatus        string
	PaymentMethod        string
	PaymentType          string
	PaymentCode          string
	PurchaseStartDate    time.Time
	PurchaseEndDate      time.Time
}

type Earnings struct {
	Username string
	Earnings float64
}

func invoiceEntities(i Invoice) transaction.InvoiceCore {
	return transaction.InvoiceCore{
		ProductTransactionID: i.ProductTransactionID,
		ProductProductID:     i.ProductProductID,
		TransactionID:        i.TransactionID,
		RestaurantID:         i.RestaurantID,
		RestaurantName:       i.RestaurantName,
		Owner:                i.Owner,
		UserID:               i.UserID,
		Buyer:                i.Buyer,
		ProductName:          i.ProductName,
		Subtotal:             i.Subtotal,
		Quantity:             i.Quantity,
		Stock:                i.Stock,
		Invoice:              i.Invoice,
		Grandtotal:           i.Grandtotal,
		PaymentStatus:        i.PaymentStatus,
		PaymentMethod:        i.PaymentMethod,
		PaymentType:          i.PaymentType,
		PaymentCode:          i.PaymentCode,
		PurchaseStartDate:    i.PurchaseStartDate,
		PurchaseEndDate:      i.PurchaseEndDate,
	}
}

func earningsModels(e Earnings) transaction.EarningsCore {
	return transaction.EarningsCore{
		Username: e.Username,
		Earnings: e.Earnings,
	}
}

// Map TransactionCore to Transaction model
func transactionModels(c transaction.TransactionCore) (Transaction, error) {
	transactionID, err := identity.GenerateID()
	if err != nil {
		return Transaction{}, err
	}

	invoice, err := identity.GenerateID()
	if err != nil {
		return Transaction{}, err
	}

	grandtotal := 0.0
	for _, pt := range c.Product_Transactions {
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
		RestaurantID:      c.RestaurantID,
		UserID:            c.UserID,
		Invoice:           invoice,
		Grandtotal:        grandtotal,
		PaymentStatus:     c.PaymentStatus,
		PaymentMethod:     c.PaymentMethod,
		PaymentType:       c.PaymentType,
		PaymentCode:       payment_code,
		PurchaseStartDate: purchaseStartDate,
		PurchaseEndDate:   purchaseEndDate,
		CreatedAt:         c.CreatedAt,
		UpdatedAt:         c.UpdatedAt,
		DeletedAt:         c.DeletedAt,
	}, nil
}

// Product_TransactionCore slice to Product_Transactions model slice
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
			CreatedAt:            c.CreatedAt,
			UpdatedAt:            c.UpdatedAt,
			DeletedAt:            c.DeletedAt,
		}
	}
	return models, nil
}

// Map Transaction model to TransactionCore
func transactionEntities(m Transaction) transaction.TransactionCore {
	return transaction.TransactionCore{
		TransactionID:     m.TransactionID,
		RestaurantID:      m.RestaurantID,
		UserID:            m.UserID,
		Invoice:           m.Invoice,
		Grandtotal:        m.Grandtotal,
		PaymentStatus:     m.PaymentStatus,
		PaymentMethod:     m.PaymentMethod,
		PaymentType:       m.PaymentType,
		PaymentCode:       m.PaymentCode,
		PurchaseStartDate: m.PurchaseStartDate,
		PurchaseEndDate:   m.PurchaseEndDate,
		CreatedAt:         m.CreatedAt,
		UpdatedAt:         m.UpdatedAt,
		DeletedAt:         m.DeletedAt,
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
			CreatedAt:            m.CreatedAt,
			UpdatedAt:            m.UpdatedAt,
			DeletedAt:            m.DeletedAt,
		}
	}
	return cores
}
