package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string `json:"name" gorm:"type:varchar(20);not null"`
	Password string `json:"password" gorm:"varchar(11);not null;unique"`
	Type     int32  `json:"type" gorm:"size:255;not null"`
}
