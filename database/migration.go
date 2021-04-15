package database

import (
	"github.com/jinzhu/gorm"
	"github.com/thalaivar-subu/paylaterservice/structs"
)

// Migrate -> Migrates table
func Migrate(db *gorm.DB) {
	db.AutoMigrate(&structs.User{}, &structs.Transaction{}, &structs.Merchant{})
}
