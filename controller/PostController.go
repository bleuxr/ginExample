package controller

import (
	"errors"
	"ginExample/common"
	"ginExample/model"
	"ginExample/response"
	"ginExample/vo"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type IPostController interface {
	RestController
	PageList(c *gin.Context)
}

type PostController struct {
	DB *gorm.DB
}

func (p PostController) Create(c *gin.Context) {
	var requestPost vo.CreatePostRequest
	// 数据验证
	if err := c.ShouldBind(&requestPost); err != nil {
		log.Print(err.Error())
		response.Fail(c, nil, "数据验证错误")
		return
	}

	// 获取登录用户
	user, _ := c.Get("user")

	// 创建post
	post := model.Post{
		UserId:     user.(model.User).ID,
		CategoryId: requestPost.CategoryId,
		Title:      requestPost.Title,
		HeadImg:    requestPost.HeadImg,
		Content:    requestPost.Content,
	}

	// 插入数据
	if err := p.DB.Create(&post).Error; err != nil {
		panic(err)
		// return
	}

	// 成功
	response.Success(c, nil, "创建成功")
}

func (p PostController) Update(c *gin.Context) {
	var requestPost vo.CreatePostRequest
	// 数据验证
	if err := c.ShouldBind(&requestPost); err != nil {
		log.Print(err.Error())
		response.Fail(c, nil, "数据验证错误")
		return
	}

	// 获取path中的id
	postId := c.Params.ByName("id")

	var post model.Post
	if err := p.DB.Where("id = ?", postId).First(&post).Error; err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		response.Fail(c, nil, "文章不存在")
		return
	}

	// // 判断当前用户是否为文章的作者
	// // 获取登录用户
	user, _ := c.Get("user")
	userId := user.(model.User).ID
	if userId != post.UserId {
		response.Fail(c, nil, "文章不属于您，请勿非法操作")
		return
	}

	// // // 更新文章
	if err := p.DB.Model(model.Post{}).Where("id = ?", postId).Updates(requestPost).Error; err != nil {
		response.Fail(c, nil, "更新失败")
		return
	}

	response.Success(c, gin.H{"post": post}, "更新成功")
}

func (p PostController) Show(c *gin.Context) {
	// 获取path中的id
	postId := c.Params.ByName("id")

	var post model.Post
	if err := p.DB.Preload("Category").Where("id = ?", postId).First(&post).Error; err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		response.Fail(c, nil, "文章不存在")
		return
	}

	response.Success(c, gin.H{"post": post}, "成功")
}

func (p PostController) Delete(c *gin.Context) {
	// 获取path中的id
	postId := c.Params.ByName("id")

	var post model.Post
	if err := p.DB.Where("id = ?", postId).First(&post).Error; err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		response.Fail(c, nil, "文章不存在")
		return
	}

	// 判断当前用户是否为文章的作者
	// 获取登录用户
	user, _ := c.Get("user")
	userId := user.(model.User).ID
	if userId != post.UserId {
		response.Fail(c, nil, "文章属于您，请勿非法操作")
		return
	}

	p.DB.Delete(&post)

	response.Success(c, gin.H{"post": post}, "成功")
}

func (p PostController) PageList(c *gin.Context) {
	pageNum, err := strconv.Atoi(c.DefaultQuery("pageNum", "1"))
	if err != nil {
		response.Fail(c, nil, "pageNum参数错误")
	}
	pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	if err != nil {
		response.Fail(c, nil, "pageSize参数错误")
	}
	// fmt.Println(pageNum, " ", pageSize)
	var posts []model.Post
	p.DB.Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&posts)

	var total int64
	p.DB.Model(model.Post{}).Count(&total)

	response.Success(c, gin.H{"data": posts, "total": total}, "成功")
}

func NewPostController() IPostController {
	db := common.GetDB()
	db.AutoMigrate(model.Post{})
	return PostController{DB: db}
}
