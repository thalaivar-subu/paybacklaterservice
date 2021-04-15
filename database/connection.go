package database

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/thalaivar-subu/paylaterservice/config"

	// Gorm Blank Import
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var Db *gorm.DB

// ConnectMysql -> Connects Mysql
func ConnectMysql() {
	db, err := gorm.Open("mysql", config.Config["mysql"])
	if err != nil {
		fmt.Println("Connection Failed to Open", err)
	} else {
		fmt.Println("Mysql Connection Established")
		Migrate(db)
	}
	Db = db
}
