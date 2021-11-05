package controller

import (
	"errors"
	"ginExample/common"
	"ginExample/model"
	"ginExample/response"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ICategoryController interface {
	RestController
}

type CategoryController struct {
	DB *gorm.DB
}

func NewCategoryController() ICategoryController {
	db := common.GetDB()
	db.AutoMigrate(model.Category{})
	return CategoryController{DB: db}
}

func (ca CategoryController) Create(c *gin.Context) {
	var requestCategory model.Category
	c.Bind(&requestCategory)
	if requestCategory.Name == "" {
		response.Fail(c, nil, "数据验证错误，分类名称必填")
		return
	}
	ca.DB.Create(&requestCategory)
	response.Success(c, gin.H{"category": requestCategory}, "")
}

func (ca CategoryController) Update(c *gin.Context) {
	//绑定body中的参数
	var requestCategory model.Category
	c.Bind(&requestCategory)
	if requestCategory.Name == "" {
		response.Fail(c, nil, "数据验证错误，分类名称必填")
		return
	}
	//获取path中的参数
	categoryId, _ := strconv.Atoi(c.Params.ByName("id"))

	var updateCategory model.Category
	if err := ca.DB.First(&updateCategory, categoryId).Error; err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		response.Fail(c, nil, "分类不存在")
		// fmt.Println("updateCategory:", updateCategory)
		return
	}

	//更新分类
	ca.DB.Model(&updateCategory).Update("name", requestCategory.Name)
	response.Success(c, nil, "更新分类成功")
}

func (ca CategoryController) Show(c *gin.Context) {
	// 获取path中的参数,将字符串强转成int类型
	categoryId, _ := strconv.Atoi(c.Params.ByName("id"))
	var category model.Category
	if err := ca.DB.First(&category, categoryId).Error; err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		response.Fail(c, nil, "分类不存在")
		return
	}

	response.Success(c, gin.H{"category": category}, "")
}
func (ca CategoryController) Delete(c *gin.Context) {
	// 获取path中的参数,将字符串强转成int类型
	categoryId, _ := strconv.Atoi(c.Params.ByName("id"))
	if err := ca.DB.Delete(model.Category{}, categoryId).Error; err != nil {
		response.Fail(c, nil, "删除失败，请重试")
		return
	}
	response.Success(c, nil, "删除成功")
}
