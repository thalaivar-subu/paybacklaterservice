package merchant

import (
	"github.com/jinzhu/gorm"
	"github.com/thalaivar-subu/paylaterservice/structs"
)

func Exists(name string, email string, db *gorm.DB) bool {
	if email != "" {
		if db.Where("email = ?", email).Find(&structs.Merchant{}).Error != nil {
			return false
		}
	} else {
		if db.Where("name=?", name).First(&structs.Merchant{}).Error != nil {
			return false
		}
	}
	return true
}

func FindOne(name string, email string, db *gorm.DB) (structs.Merchant, error) {
	Merchant := structs.Merchant{}
	err := db.Where("name=?", name).First(&Merchant).Error
	if err != nil {
		return structs.Merchant{}, err
	}
	return Merchant, nil
}
