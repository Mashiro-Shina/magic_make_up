package dto

import "magicMakeup/entities"

type CommentDTO struct {
	ID int
	Content string
	PublishTime int64
	LikeNum int
	ReplyNum int
}

func ToCommentDTO(comment *entities.Comment) *CommentDTO {
	return &CommentDTO{
		ID: comment.ID,
		Content: comment.Content,
		PublishTime: comment.PublishTime,
		LikeNum: comment.LikeNum,
		ReplyNum: comment.ReplyNum,
	}
}
