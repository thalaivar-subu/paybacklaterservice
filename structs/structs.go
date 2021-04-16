package structs

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Name        string  `json:"name"`
	Email       string  `json:"email"`
	CreditLimit float64 `json:"credit_limit"`
	Dues        float64 `json:"dues"`
}

type Merchant struct {
	gorm.Model
	Name            string  `json:"name"`
	Email           string  `json:"email"`
	DiscountPercent float64 `json:"discount_percent"`
}

type Transaction struct {
	gorm.Model
	Amount          float64 `json:"amount"`
	UserID          uint
	MerchantID      uint
	AmountToService float64 `json:"amount_to_service"`
}
