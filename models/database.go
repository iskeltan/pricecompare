package models

import (
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // for gorm
)

// Database method will be export
func Database() (Db *gorm.DB, err error) {
	Db, err = gorm.Open("postgres", "dbname=pricecompare sslmode=disable")
	//defer db.Close()
	if err != nil {
		panic(err)
	}
	log.Print("[+] has been connect to database server.")

	Db.LogMode(true)

	return
}

// Migrate for make database migrations
func Migrate() {

	log.Print("[+] Migration started.")

	Db, _ := Database()

	Db.AutoMigrate(&User{})
	Db.AutoMigrate(&Token{})

	log.Print("[+] Migration has been completed.")
}
