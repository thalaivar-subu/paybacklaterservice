package txn

import (
	"errors"
	"strconv"

	"github.com/jinzhu/gorm"
	"github.com/thalaivar-subu/paylaterservice/cmd/merchant"
	user "github.com/thalaivar-subu/paylaterservice/cmd/user"
	"github.com/thalaivar-subu/paylaterservice/structs"
)

func CreateTxn(userName string, merchantName string, amount string, db *gorm.DB) (string, error) {
	User, err := user.FindOne(userName, "", db)
	if err != nil {
		errMsg := errors.New("Not a valid User")
		return "", errMsg
	}
	Merchant, err := merchant.FindOne(merchantName, "", db)
	if err != nil {
		errMsg := errors.New("Not a valid merchant")
		return "", errMsg
	}
	floatAmount, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		errMsg := errors.New(err.Error())
		return "", errMsg
	}
	if User.CreditLimit < floatAmount {
		errMsg := errors.New("reason: credit limit")
		return "", errMsg
	}
	amountToService := floatAmount * (Merchant.DiscountPercent / 100)
	err = db.Create(&structs.Transaction{UserID: User.ID,
		MerchantID:      Merchant.ID,
		Amount:          floatAmount,
		AmountToService: amountToService,
	}).Error
	if err != nil {
		errMsg := errors.New(err.Error())
		return "", errMsg
	}
	newCreditLimit := User.CreditLimit - floatAmount
	err = db.Model(&structs.User{}).Where("id = ?", User.ID).Updates(map[string]interface{}{"credit_limit": newCreditLimit, "dues": floatAmount}).Error
	if err != nil {
		errMsg := errors.New(err.Error())
		return "", errMsg
	}
	return "success!", nil

}
