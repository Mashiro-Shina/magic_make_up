package dto

import (
	"magicMakeup/common"
	"magicMakeup/entities"
)

type StarDTO struct {
	ID int
	Content string
	Images string
	PublishTime int64
	LikeNum int
	CommentNum int
	ForwardNum int
	IsForward int
	PreviousID int
}

func ToStarDTO(star *entities.Star) *StarDTO {
	base64Image, _ := common.ReadAndEncodingImage(star.Images)

	return &StarDTO{
		ID:          star.ID,
		Content:     star.Content,
		Images:      base64Image,
		PublishTime: star.PublishTime,
		LikeNum:     star.LikeNum,
		CommentNum:  star.CommentNum,
		ForwardNum:  star.ForwardNum,
		IsForward:   star.IsForward,
		PreviousID:  star.PreviousID,
	}
}
