package config

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB
var err error

func Connect() {
	DB, err = gorm.Open(mysql.Open("root:@tcp(127.0.0.1:3306)/belajargolang"))
	if err != nil {
		log.Println("Connection Failed", err)
	} else {
		log.Println("Connected")
	}

}
