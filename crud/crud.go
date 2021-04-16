package crud

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

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
	err = db.Create(&structs.User{Name: name, Email: email, CreditLimit: floatAmount, Dues: 0}).Error
	if err != nil {
		errMsg := errors.New(err.Error())
		return false, "", errMsg
	}
	return true, name + "(" + amount + ")", nil
}

func CreateMerchant(name string, email string, discount string, db *gorm.DB) (bool, string, error) {
	records := structs.Merchant{}
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

func CreateTxn(userName string, merchantName string, amount string, db *gorm.DB) (bool, string, error) {
	User := structs.User{}
	if db.Where("name = ?", userName).First(&User).Error != nil {
		errMsg := errors.New("Not a valid User")
		return false, "", errMsg
	}
	Merchant := structs.Merchant{}
	if db.Where("name = ?", merchantName).First(&Merchant).Error != nil {
		errMsg := errors.New("Not a valid merchant")
		return false, "", errMsg
	}
	floatAmount, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		errMsg := errors.New(err.Error())
		return false, "", errMsg
	}
	if User.CreditLimit < floatAmount {
		errMsg := errors.New("reason: credit limit")
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
	err = db.Model(&structs.User{}).Where("id = ?", User.ID).Updates(map[string]interface{}{"credit_limit": newCreditLimit, "dues": floatAmount}).Error
	if err != nil {
		errMsg := errors.New(err.Error())
		return false, "", errMsg
	}
	return true, "success!", nil

}

func GetUsersAtCredLimit(db *gorm.DB) (bool, string, error) {
	rows, err := db.Where("credit_limit = ?", 0).Find(&[]structs.User{}).Rows()
	fmt.Println(rows)
	if err != nil {
		errMsg := errors.New(err.Error())
		return false, "", errMsg
	}
	defer rows.Close()
	usersAtCredLimit := make([]string, 0)
	for rows.Next() {
		var user structs.User
		db.ScanRows(rows, &user)
		usersAtCredLimit = append(usersAtCredLimit, user.Name)
	}
	return true, strings.Join(usersAtCredLimit, "\n"), nil
}

func GetTotalDues(db *gorm.DB) (bool, string, error) {
	rows, err := db.Find(&[]structs.User{}).Rows()
	if err != nil {
		errMsg := errors.New(err.Error())
		return false, "", errMsg
	}
	defer rows.Close()
	duesSlice := make([]string, 0)
	total := 0.0
	for rows.Next() {
		var user structs.User
		db.ScanRows(rows, &user)
		total += user.Dues
		duesSlice = append(duesSlice, user.Name+":"+fmt.Sprintf("%.2f", user.Dues))
	}
	duesSlice = append(duesSlice, "total:"+fmt.Sprintf("%.2f", total))
	return true, strings.Join(duesSlice, "\n"), nil
}

func GetDiscount(merchantName string, db *gorm.DB) (bool, string, error) {
	rows, err := db.Where("name=?", merchantName).Find(&structs.Merchant{}).Rows()
	if err != nil {
		errMsg := errors.New(err.Error())
		return false, "", errMsg
	}
	defer rows.Close()
	discountString := ""
	for rows.Next() {
		var merchant structs.Merchant
		db.ScanRows(rows, &merchant)
		discountString = merchant.Name + "(" + fmt.Sprintf("%.2f", merchant.DiscountPercent) + "%)"
	}
	return true, discountString, nil
}

func GetUserDues(userName string, db *gorm.DB) (bool, string, error) {
	rows, err := db.Where("name=?", userName).Find(&structs.User{}).Rows()
	if err != nil {
		errMsg := errors.New(err.Error())
		return false, "", errMsg
	}
	defer rows.Close()
	duesString := ""
	for rows.Next() {
		var user structs.User
		db.ScanRows(rows, &user)
		duesString = user.Name + "(" + fmt.Sprintf("%.2f", user.Dues) + ")"
	}
	return true, duesString, nil
}

func PayBack(userName string, amount string, db *gorm.DB) (bool, string, error) {
	User := structs.User{}
	if db.Where("name=?", userName).First(&User).Error != nil {
		errMsg := errors.New("Not a valid User")
		return false, "", errMsg
	}
	floatAmount, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		errMsg := errors.New(err.Error())
		return false, "", errMsg
	}
	dues := User.Dues - floatAmount
	if floatAmount > User.Dues {
		dues = 0
	}
	err = db.Model(&structs.User{}).Where("id = ?", User.ID).Updates(map[string]interface{}{"credit_limit": User.CreditLimit + floatAmount, "dues": dues}).Error
	if err != nil {
		errMsg := errors.New(err.Error())
		return false, "", errMsg
	}
	resultMsg := "user3(dues: " + fmt.Sprintf("%.2f", User.Dues) + ")"
	if floatAmount > User.Dues {
		resultMsg += " , Credit limit is increased to " + fmt.Sprintf("%.2f", User.CreditLimit+(floatAmount-User.Dues))
	}
	return true, resultMsg, nil
}

func UpdateMerchantDiscount(userName string, amount string, db *gorm.DB) (bool, string, error) {
	Merchant := structs.Merchant{}
	if db.Where("name=?", userName).First(&Merchant).Error != nil {
		errMsg := errors.New("Not a valid Merchant")
		return false, "", errMsg
	}
	floatAmount, err := strconv.ParseFloat(helper.TrimSuffix(amount, "%"), 64)
	if err != nil {
		errMsg := errors.New(err.Error())
		return false, "", errMsg
	}
	err = db.Model(&structs.Merchant{}).Where("id = ?", Merchant.ID).Updates(map[string]interface{}{"discount_percent": floatAmount}).Error
	if err != nil {
		errMsg := errors.New(err.Error())
		return false, "", errMsg
	}
	return true, Merchant.Name + "(" + fmt.Sprintf("%.2f", floatAmount) + "%)", nil
}
