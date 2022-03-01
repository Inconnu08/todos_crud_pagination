package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"path"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"os"
)

var DB *gorm.DB

type Database struct {
	*gorm.DB
}

func OpenDbConnection() *gorm.DB {
	dialect := os.Getenv("DB_DIALECT")
	username := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	host := os.Getenv("DB_HOST")
	var db *gorm.DB
	var err error
	if dialect == "sqlite3" {
		db, err = gorm.Open("sqlite3", path.Join(".", "app.db"))
	} else {
		databaseUrl := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable ", host, username, password, dbName)
		db, err = gorm.Open(dialect, databaseUrl)
	}

	if err != nil {
		fmt.Println("db err: ", err)
		os.Exit(-1)
	}

	db.DB().SetMaxIdleConns(10)
	db.LogMode(true)

	DB = db
	return DB
}

// Delete the database after running testing cases.
func DeleteDb(db *gorm.DB) error {
	dbName := os.Getenv("DB_NAME")
	_ = db.Close()
	err := os.Remove(path.Join(".", dbName))
	return err
}

// Using this function to get a connection, you can create your connection pool here.
func GetDb() *gorm.DB {
	return DB
}
