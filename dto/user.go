package dto

import (
	"log"
	"magicMakeup/common"
	"magicMakeup/entities"
)

type UserDTO struct {
	UserID int
	UserName string
	Signature string
	Phone string
	Avatar string
	FollowingNum int
	FollowersNum int
}

func NewUserDTO(id int, name string, signature string, phone string, avatar string, followingNum int, followersNum int) *UserDTO {
	return &UserDTO{
		UserID: id,
		UserName: name,
		Signature: signature,
		Phone: phone,
		Avatar: avatar,
		FollowingNum: followingNum,
		FollowersNum: followersNum,
	}
}

func ToUserDTO(user *entities.User) *UserDTO {
	avatar, err := common.ReadAndEncodingImage(user.Avatar)
	if err != nil {
		log.Printf("读取头像失败: %v", err)
		return nil
	}
	return NewUserDTO(user.ID, user.Name, user.Signature, user.Phone, avatar, user.FollowingNum, user.FollowersNum)
}
