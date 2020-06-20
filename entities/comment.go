package entities

type Comment struct {
	ID int	`json:"id"`
	ReplyID int	`json:"reply_id"`
	UserID int	`json:"user_id"`
	StarID int	`json:"star_id"`
	Content string	`json:"content"`
	PublishTime int64	`json:"publish_time"`
	LikeNum int `json:"like_num"`
	ReplyNum int `json:"reply_num"`
}

func NewComment(id int, userID int, starID int, replyID int, content string, publishTime int64) *Comment {
	return &Comment{
		ID: id,
		UserID: userID,
		StarID: starID,
		ReplyID: replyID,
		Content: content,
		PublishTime: publishTime,
	}
}
