package merchant

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/jinzhu/gorm"
	"github.com/thalaivar-subu/paylaterservice/helper"
	"github.com/thalaivar-subu/paylaterservice/structs"
)

func UpdateMerchantDiscount(userName string, amount string, db *gorm.DB) (string, error) {
	Merchant, err := FindOne(userName, "", db)
	if err != nil {
		errMsg := errors.New("Not a valid Merchant")
		return "", errMsg
	}
	floatAmount, err := strconv.ParseFloat(helper.TrimSuffix(amount, "%"), 64)
	if err != nil {
		errMsg := errors.New(err.Error())
		return "", errMsg
	}
	err = db.Model(&structs.Merchant{}).Where("id = ?", Merchant.ID).Updates(map[string]interface{}{"discount_percent": floatAmount}).Error
	if err != nil {
		errMsg := errors.New(err.Error())
		return "", errMsg
	}
	return Merchant.Name + "(" + fmt.Sprintf("%.2f", floatAmount) + "%)", nil
}

func GetDiscount(merchantName string, db *gorm.DB) (string, error) {
	rows, err := db.Where("name=?", merchantName).Find(&structs.Merchant{}).Rows()
	if err != nil {
		errMsg := errors.New(err.Error())
		return "", errMsg
	}
	defer rows.Close()
	discountString := ""
	for rows.Next() {
		var merchant structs.Merchant
		db.ScanRows(rows, &merchant)
		discountString = merchant.Name + "(" + fmt.Sprintf("%.2f", merchant.DiscountPercent) + "%)"
	}
	return discountString, nil
}

func CreateMerchant(name string, email string, discount string, db *gorm.DB) (string, error) {
	if Exists("", email, db) {
		errMsg := errors.New("Merchant Already Exists")
		return "", errMsg
	}
	floatDiscount, err := strconv.ParseFloat(helper.TrimSuffix(discount, "%"), 64)
	if err != nil {
		errMsg := errors.New(err.Error())
		return "", errMsg
	}
	err = db.Create(&structs.Merchant{Name: name, Email: email, DiscountPercent: floatDiscount}).Error
	if err != nil {
		errMsg := errors.New(err.Error())
		return "", errMsg
	}
	return name + "(" + discount + ")", nil
}
