package user

import (
	"github.com/jinzhu/gorm"
	"github.com/thalaivar-subu/paylaterservice/structs"
)

func Exists(name string, email string, db *gorm.DB) bool {
	if email != "" {
		if db.Where("email = ?", email).Find(&structs.User{}).Error != nil {
			return false
		}
	} else {
		if db.Where("name=?", name).First(&structs.User{}).Error != nil {
			return false
		}
	}
	return true
}

func FindOne(name string, email string, db *gorm.DB) (structs.User, error) {
	User := structs.User{}
	err := db.Where("name=?", name).First(&User).Error
	if err != nil {
		return structs.User{}, err
	}
	return User, nil
}
