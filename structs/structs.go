package structs

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Name        string  `json:"name"`
	Email       string  `json:"email"`
	CreditLimit float64 `json:"credit_limit"`
}

type Merchant struct {
	gorm.Model
	Name            string  `json:"name"`
	Email           string  `json:"email"`
	DiscountPercent float64 `json:"discount_percent"`
}

type Transaction struct {
	gorm.Model
	Amount     float64 `json:"amount"`
	UserID     int
	MerchantID int
}
