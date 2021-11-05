package controller

import (
	"fmt"
	"ginExample/common"
	"ginExample/model"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Register(c *gin.Context) {
	db := common.GetDB()

	name := c.PostForm("name")
	password := c.PostForm("password")

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
	hasedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"code": 500, "msg": "加密出错"})
		return
	}
	newUser := model.User{
		Name:     name,
		Password: string(hasedPassword),
	}
	db.Create(&newUser)

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "注册成功",
	})
}

func isNameExist(db *gorm.DB, name string) bool {
	var user model.User
	db.Where("name=?", name).First(&user)
	if user.ID != 0 {
		return true
	}
	return false
}

func Login(c *gin.Context) {
	db := common.GetDB()

	name := c.PostForm("name")
	password := c.PostForm("password")

	//check password
	var user model.User
	db.Where("name = ?", name).First(&user)
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"code": 400, "msg": "密码错误"})
		return
	}

	//token
	token, err := common.ReleaseToken(user)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"code": 500, "msg": "系统异常！"})
		log.Printf("token generate error: ", err)
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"data": gin.H{"token": token},
		"msg":  "登录成功",
	})
}

func Info(c *gin.Context) {
	user, _ := c.Get("user")
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"user": user}})
}
