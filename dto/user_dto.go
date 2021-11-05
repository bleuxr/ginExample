package dto

import "ginExample/model"

type UserDto struct {
	Name string `jason:"name"`
	Type int32
}

func ToUserDto(user model.User) UserDto {
	return UserDto{
		Name: user.Name,
		Type: user.Type,
	}
}
