package common

import (
	"fmt"
	"ginExample/model"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() *gorm.DB {
	// driverName := viper.GetString("datasource.driverName")
	host := viper.GetString("datasource.host")
	port := viper.GetString("datasource.port")
	database := viper.GetString("datasource.database")
	username := viper.GetString("datasource.username")
	password := viper.GetString("datasource.password")
	charset := viper.GetString("datasource.charset")
	loc := viper.GetString("datasource.loc")
	args := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=%s",
		username,
		password,
		host,
		port,
		database,
		charset,
		loc)
	fmt.Printf(args)
	// dsn := "root:MyNewPass4!@tcp(127.0.0.1:3306)/mydb?charset=utf8mb4&parseTime=True&loc=Local"
	// db, err := gorm.Open(driverName.Open(args), &gorm.Config{})
	db, err := gorm.Open(mysql.Open(args), &gorm.Config{})
	if err != nil {
		fmt.Printf("fail to connect database! err:" + err.Error())
	}
	//create table
	db.AutoMigrate(&model.User{})
	DB = db
	return db
}

func GetDB() *gorm.DB {
	return DB
}
