package controller

import (
	"ginExample/model"
	"ginExample/repository"
	"ginExample/response"
	"ginExample/vo"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ICategoryController interface {
	RestController
}

type CategoryController struct {
	Repository repository.CategoryRepository
}

func NewCategoryController() ICategoryController {
	repository := repository.NewCategoryRepository()
	repository.DB.AutoMigrate(model.Category{})
	return CategoryController{Repository: repository}
}

func (ca CategoryController) Create(c *gin.Context) {
	var requestCategory vo.CreateCategoryRequest
	if err := c.ShouldBind(&requestCategory); err != nil {
		response.Fail(c, nil, "数据验证错误，分类名称必填")
		return
	}

	category, err := ca.Repository.Create(requestCategory.Name)
	if err != nil {
		// response.Fail(c, nil, "创建失败")
		panic(err)
		return
	}
	response.Success(c, gin.H{"category": category}, "")
}

func (ca CategoryController) Update(c *gin.Context) {
	//绑定body中的参数
	var requestCategory vo.CreateCategoryRequest
	if err := c.ShouldBind(&requestCategory); err != nil {
		response.Fail(c, nil, "数据验证错误，分类名称必填")
		return
	}
	//获取path中的参数
	categoryId, _ := strconv.Atoi(c.Params.ByName("id"))

	// var updateCategory model.Category
	// if err := ca.DB.First(&updateCategory, categoryId).Error; err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
	// 	response.Fail(c, nil, "分类不存在")
	// 	// fmt.Println("updateCategory:", updateCategory)
	// 	return
	// }

	// //更新分类
	// ca.DB.Model(&updateCategory).Update("name", requestCategory.Name)

	updateCategory, err := ca.Repository.SelectById(categoryId)
	if err != nil {
		response.Fail(c, nil, "分类不存在")
		return
	}

	category, err := ca.Repository.Update(*updateCategory, requestCategory.Name)
	if err != nil {
		response.Fail(c, nil, "更新出错")
		return
	}

	response.Success(c, gin.H{"category": category}, "更新分类成功")
}

func (ca CategoryController) Show(c *gin.Context) {
	// 获取path中的参数,将字符串强转成int类型
	categoryId, _ := strconv.Atoi(c.Params.ByName("id"))
	// var category model.Category
	// if err := ca.DB.First(&category, categoryId).Error; err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
	// 	response.Fail(c, nil, "分类不存在")
	// 	return
	// }
	category, err := ca.Repository.SelectById(categoryId)
	if err != nil {
		response.Fail(c, nil, "分类不存在")
		return
	}

	response.Success(c, gin.H{"category": category}, "")
}
func (ca CategoryController) Delete(c *gin.Context) {
	// 获取path中的参数,将字符串强转成int类型
	categoryId, _ := strconv.Atoi(c.Params.ByName("id"))
	// if err := ca.DB.Delete(model.Category{}, categoryId).Error; err != nil {
	// 	response.Fail(c, nil, "删除失败，请重试")
	// 	return
	// }
	err := ca.Repository.DeleteById(categoryId)
	if err != nil {
		response.Fail(c, nil, "删除失败")
		return
	}
	response.Success(c, nil, "删除成功")
}
