package user

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/thalaivar-subu/paylaterservice/structs"
)

func CreateUser(name string, email string, amount string, db *gorm.DB) (string, error) {
	if Exists("", email, db) {
		errMsg := errors.New("User Already Exists")
		return "", errMsg
	}
	floatAmount, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		errMsg := errors.New(err.Error())
		return "", errMsg
	}
	err = db.Create(&structs.User{Name: name, Email: email, CreditLimit: floatAmount, Dues: 0}).Error
	if err != nil {
		errMsg := errors.New(err.Error())
		return "", errMsg
	}
	return name + "(" + amount + ")", nil
}

func GetUsersAtCredLimit(db *gorm.DB) (string, error) {
	rows, err := db.Where("credit_limit = ?", 0).Find(&[]structs.User{}).Rows()
	if err != nil {
		errMsg := errors.New(err.Error())
		return "", errMsg
	}
	defer rows.Close()
	usersAtCredLimit := make([]string, 0)
	for rows.Next() {
		var user structs.User
		db.ScanRows(rows, &user)
		usersAtCredLimit = append(usersAtCredLimit, user.Name)
	}
	return strings.Join(usersAtCredLimit, "\n"), nil
}

func GetTotalDues(db *gorm.DB) (string, error) {
	rows, err := db.Find(&[]structs.User{}).Rows()
	if err != nil {
		errMsg := errors.New(err.Error())
		return "", errMsg
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
	return strings.Join(duesSlice, "\n"), nil
}

func GetUserDues(userName string, db *gorm.DB) (string, error) {
	rows, err := db.Where("name=?", userName).Find(&structs.User{}).Rows()
	if err != nil {
		errMsg := errors.New(err.Error())
		return "", errMsg
	}
	defer rows.Close()
	duesString := ""
	for rows.Next() {
		var user structs.User
		db.ScanRows(rows, &user)
		duesString = user.Name + "(" + fmt.Sprintf("%.2f", user.Dues) + ")"
	}
	return duesString, nil
}

func PayBack(userName string, amount string, db *gorm.DB) (string, error) {
	User, err := FindOne(userName, "", db)
	if err != nil {
		errMsg := errors.New("Not a valid User")
		return "", errMsg
	}
	floatAmount, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		errMsg := errors.New(err.Error())
		return "", errMsg
	}
	dues := User.Dues - floatAmount
	if floatAmount > User.Dues {
		dues = 0
	}
	err = db.Model(&structs.User{}).Where("id = ?", User.ID).Updates(map[string]interface{}{"credit_limit": User.CreditLimit + floatAmount, "dues": dues}).Error
	if err != nil {
		errMsg := errors.New(err.Error())
		return "", errMsg
	}
	resultMsg := "user3(dues: " + fmt.Sprintf("%.2f", User.Dues) + ")"
	if floatAmount > User.Dues {
		resultMsg += " , Credit limit is increased to " + fmt.Sprintf("%.2f", User.CreditLimit+(floatAmount-User.Dues))
	}
	return resultMsg, nil
}
