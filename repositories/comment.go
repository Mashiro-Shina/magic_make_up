package repositories

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"magicMakeup/common"
	"magicMakeup/entities"
	"time"
)

func HasKey(key string) (bool, error) {
	res, err := redis.Bool(redisConn.Do("exists", key))
	if err != nil {
		return false, err
	}
	return res, nil
}

func InitCommentID(key string) error {
	_, err := redisConn.Do("set", key, 1)
	return err
}

func SetCommentID(key string, value int) error {
	_, err := redisConn.Do("set", key, value)
	return err
}

func GetCommentIDFromRedis(key string) (int, error) {
	res, err := redis.Int(redisConn.Do("get", key))
	return res, err
}

func InsertComment(comment *entities.Comment) (*entities.Comment, error) {
	res := mysqlConn.Create(comment)
	if res.Error != nil {
		return nil, err
	}
	return comment, nil
}

func ReplyComment(comment *entities.Comment) error {
	originalComment, _ := SearchCommentByID(comment.ReplyID)
	originalComment.ReplyNum++
	_, _ = UpdateComment(originalComment)

	// 记录回复评论的 id
	_, err := redisConn.Do("hset", fmt.Sprintf("comment%d:reply_comments", originalComment.ID), comment.ID, time.Now().Unix())
	return err
}

func UpdateComment(comment *entities.Comment) (*entities.Comment, error) {
	res := mysqlConn.Save(comment)
	if res.Error != nil {
		return nil, res.Error
	}
	return comment, nil
}

func DeleteComment(commentID int) error {
	comment, _ := SearchCommentByID(commentID)
	// 删除的评论如果是一条回复，则更新原评论的回复列表，并递减回复数
	if comment.ReplyID != -1 {
		originalComment, _ := SearchCommentByID(comment.ReplyID)
		originalComment.ReplyNum--
		_, _ = UpdateComment(originalComment)
		_, err := redisConn.Do("hdel", fmt.Sprintf("comment%d:reply_comments", comment.ReplyID), comment.ID)
		if err != nil {
			return err
		}
	}
	res := mysqlConn.Delete(comment)
	return res.Error
}

func SearchCommentByID(commentID int) (*entities.Comment, error) {
	comment := &entities.Comment{}
	res := mysqlConn.First(comment, commentID)
	if res.Error != nil {
		return nil, res.Error
	}
	return comment, nil
}

func GetReplyCommentList(commentID int) ([]*entities.Comment, error) {
	var comments []*entities.Comment
	res := mysqlConn.Where("reply_id=?", commentID).Find(&comments)
	if res.Error != nil {
		return nil, res.Error
	}
	return comments, nil
}

func LikeComment(userID int, commentID int) error {
	comment, _ := SearchCommentByID(commentID)
	comment.LikeNum++
	_, _ = UpdateComment(comment)

	_, err := redisConn.Do("hset", fmt.Sprintf("comment%d:like_users", comment.ID), userID, time.Now().Unix())
	return err
}

func CancelLikeComment(userID int, commentID int) error {
	comment, _ := SearchCommentByID(commentID)
	comment.LikeNum--
	_, _ = UpdateComment(comment)

	_, err := redisConn.Do("hdel", fmt.Sprintf("comment%d:like_users", comment.ID), userID)
	return err
}

func GetCommentList(starID int) ([]*entities.Comment, error) {
	var comments []*entities.Comment
	res := mysqlConn.Where("star_id=? and reply_id=-1", starID).Find(&comments)
	if res.Error != nil {
		return nil, res.Error
	}
	return comments, nil
}

func GetLikeCommentList(commentID int) ([]*entities.User, error) {
	items, err := redis.Values(redisConn.Do("hkeys", fmt.Sprintf("comment%d:like_users", commentID)))
	if err != nil {
		return nil, err
	}

	users := make([]*entities.User, len(items))
	for i, item := range items {
		userID := common.Int(string(item.([]byte)))
		user, _ := SearchUserByID(userID)
		users[i] = user
	}

	return users, nil
}
