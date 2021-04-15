package crud

import (
	"errors"
	"strconv"

	"github.com/jinzhu/gorm"
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
	floatDiscount, err := strconv.ParseFloat(discount, 64)
	if err != nil {
		errMsg := errors.New(err.Error())
		return false, "", errMsg
	}
	err = db.Create(&structs.Merchant{Name: name, Email: email, DiscountPercent: floatDiscount}).Error
	if err != nil {
		errMsg := errors.New(err.Error())
		return false, "", errMsg
	}
	return false, name + "(" + discount + "%)", nil
}
