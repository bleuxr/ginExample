package controller

import (
	"fmt"
	"ginExample/common"
	"ginExample/dto"
	"ginExample/model"
	"ginExample/response"
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
		response.Response(c, http.StatusUnprocessableEntity, 422, nil, "用户名太短")
		return
	}

	//query & insert
	if isNameExist(db, name) {
		response.Response(c, http.StatusUnprocessableEntity, 422, nil, "用户已经存在")
		return
	}
	hasedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		response.Response(c, http.StatusInternalServerError, 500, nil, "加密出错")
		return
	}
	newUser := model.User{
		Name:     name,
		Password: string(hasedPassword),
	}
	db.Create(&newUser)
	response.Success(c, nil, "注册成功")
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

	// name := c.PostForm("name")
	// password := c.PostForm("password")
	var requestUser = model.User{}
	c.Bind(&requestUser)

	//check password
	var user model.User
	db.Where("name = ?", requestUser.Name).First(&user)
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(requestUser.Password)); err != nil {
		response.Response(c, http.StatusUnprocessableEntity, 400, nil, "密码错误")
		return
	}

	//token
	token, err := common.ReleaseToken(user)
	if err != nil {
		response.Response(c, http.StatusUnprocessableEntity, 500, nil, "系统异常！")
		log.Printf("token generate error: ", err)
		return
	}
	response.Success(c, gin.H{"token": token}, "登录成功")
}

func Info(c *gin.Context) {
	user, _ := c.Get("user")
	response.Success(c, gin.H{"user": dto.ToUserDto(user.(model.User))}, "获取用户信息成功")
}
