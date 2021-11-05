package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string
	Password string
	Type     int32
}

func main() {
	db := InitDB()
	// db.Close() is removed from v2
	// defer db.Close()

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.POST("/api/auth/register", func(c *gin.Context) {
		name := c.PostForm("name")
		password := c.PostForm(("password"))

		fmt.Println(name + " " + password)

		if len(name) < 3 {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "name too short"})
			return
		}

		//query & insert
		if isNameExist(db, name) {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "用户已经存在"})
			return
		}
		newUser := User{
			Name:     name,
			Password: password,
		}
		db.Create(&newUser)

		c.JSON(200, gin.H{
			"msg": "注册成功",
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080
}

func InitDB() *gorm.DB {
	dsn := "root:MyNewPass4!@tcp(127.0.0.1:3306)/mydb?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Printf("fail to connect database! err:" + err.Error())
	}
	//create table
	db.AutoMigrate(&User{})
	return db
}

func isNameExist(db *gorm.DB, name string) bool {
	var user User
	db.Where("name=?", name).First(&user)
	if user.ID != 0 {
		return true
	}
	return false
}
