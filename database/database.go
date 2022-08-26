package database

import (
	"fmt"
	"log"

	"git.xenonstack.com/xs-onboarding/document-manage/config"
	"github.com/jinzhu/gorm"

	// postgres dialects
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

//CreateDatabaseTables funtion for create Table
func CreateDatabaseTables() {
	// connecting db using connection string
	db := config.DB
	if !(db.HasTable(UploadLinks{})) {
		db.CreateTable(UploadLinks{})
	}
	if !(db.HasTable(PolicyLinks{})) {
		db.CreateTable(PolicyLinks{})
	}

	// Database migration
	db.AutoMigrate(&UploadLinks{})
	db.AutoMigrate(&PolicyLinks{})
}

//CreateDatabase :
func CreateDatabase() error {
	// connecting with postgres database root db
	db, err := gorm.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Conf.Database.Host,
		config.Conf.Database.Port,
		config.Conf.Database.User,
		config.Conf.Database.Pass,
		"postgres",
		config.Conf.Database.Ssl))
	if err != nil {
		log.Print(err)
		return err
	}
	defer db.Close()
	// executing create database query.
	db.Exec(fmt.Sprintf("create database %s;", config.Conf.Database.Name))
	return nil
}
