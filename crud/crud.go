package crud

import (
	"errors"
	"strconv"

	"github.com/jinzhu/gorm"
	"github.com/thalaivar-subu/paylaterservice/helper"
	"github.com/thalaivar-subu/paylaterservice/structs"
)

func CreateUser(name string, email string, amount string, db *gorm.DB) (bool, string, error) {
	records := structs.User{}
	if db.Where("email = ?", email).Find(&records).RowsAffected != 0 {
		errMsg := errors.New("User Already Exists")
		return false, "", errMsg
	}
	floatAmount, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		errMsg := errors.New(err.Error())
		return false, "", errMsg
	}
	err = db.Create(&structs.User{Name: name, Email: email, CreditLimit: floatAmount}).Error
	if err != nil {
		errMsg := errors.New(err.Error())
		return false, "", errMsg
	}
	return true, name + "(" + amount + ")", nil
}

func CreateMerchant(name string, email string, discount string, db *gorm.DB) (bool, string, error) {
	records := structs.User{}
	if db.Where("email = ?", email).Find(&records).RowsAffected != 0 {
		errMsg := errors.New("Merchant Already Exists")
		return false, "", errMsg
	}
	floatDiscount, err := strconv.ParseFloat(helper.TrimSuffix(discount, "%"), 64)
	if err != nil {
		errMsg := errors.New(err.Error())
		return false, "", errMsg
	}
	err = db.Create(&structs.Merchant{Name: name, Email: email, DiscountPercent: floatDiscount}).Error
	if err != nil {
		errMsg := errors.New(err.Error())
		return false, "", errMsg
	}
	return true, name + "(" + discount + ")", nil
}

func CreateTxn(user string, merchant string, amount string, db *gorm.DB) (bool, string, error) {
	User := structs.User{}
	if db.Where("name = ?", user).Find(&User).Error != nil {
		errMsg := errors.New("Not a valid User")
		return false, "", errMsg
	}
	Merchant := structs.Merchant{}
	if db.Where("name = ?", merchant).Find(&Merchant).Error != nil {
		errMsg := errors.New("Not a valid merchant")
		return false, "", errMsg
	}
	floatAmount, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		errMsg := errors.New(err.Error())
		return false, "", errMsg
	}
	if User.CreditLimit < floatAmount {
		errMsg := errors.New("(reason: credit limit)")
		return false, "", errMsg
	}
	amountToService := floatAmount * (Merchant.DiscountPercent / 100)
	err = db.Create(&structs.Transaction{UserID: User.ID,
		MerchantID:      Merchant.ID,
		Amount:          floatAmount,
		AmountToService: amountToService,
	}).Error
	if err != nil {
		errMsg := errors.New(err.Error())
		return false, "", errMsg
	}
	newCreditLimit := User.CreditLimit - floatAmount
	err = db.Model(&structs.User{}).Where("id = ?", User.ID).Updates(map[string]interface{}{"credit_limit": newCreditLimit}).Error
	if err != nil {
		errMsg := errors.New(err.Error())
		return false, "", errMsg
	}
	return true, "success!", nil

}
