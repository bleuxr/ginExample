package model

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type Post struct {
	ID         string `json:"id" gorm:"type:char(36);primary_key"`
	UserId     uint   `json:"user_id" gorm:"not null"`
	CategoryId uint   `json:"category_id" gorm:"not null"`
	//自动关联外键
	Category  *Category
	Title     string `json:"title" gorm:"type:varchar(50);not null"`
	HeadImg   string `json:"head_img"`
	Content   string `json:"content" gorm:"type:text;not null"`
	CreatedAt Time   `json:"created_at" gorm:"type:timestamp"`
	UpdatedAt Time   `json:"updated_at" gorm:"type:timestamp"`
}

func (post *Post) BeforeCreate(tx *gorm.DB) error {
	post.ID = uuid.NewV4().String()
	return nil
}
