package entities

import (
	"strconv"
)

type Star struct {
	ID int
	UserID int
	Content string
	Images string
	PublishTime int64
	LikeNum int
	CommentNum int
	ForwardNum int
	IsForward int
	PreviousID int
}

func NewStar(userID int, content string, images string, publishTime int64, isForward string) *Star {
	state, _ := strconv.Atoi(isForward)
	return &Star{
		UserID: userID,
		Content: content,
		Images: images,
		PublishTime: publishTime,
		IsForward: state,
	}
}
