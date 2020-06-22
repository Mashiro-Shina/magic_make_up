package repositories

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"magicMakeup/common"
	"magicMakeup/dto"
	"strings"
)

func GetCommentNotifications(userID int, userType string) ([][]*dto.UserDTO, []*dto.CommentDTO, error) {
	keys, err := redis.Values(redisConn.Do("keys", fmt.Sprintf("notice:user%d:comment*:%s",
		userID, userType)))
	if err != nil {
		return nil, nil, err
	}

	comments := make([]*dto.CommentDTO, len(keys))
	users := make([][]*dto.UserDTO, len(keys))
	for i, key := range keys {
		items, err := redis.Values(redisConn.Do("hkeys", string(key.([]byte))))
		if err != nil {
			return nil, nil, err
		}
		var row []*dto.UserDTO
		for _, item := range items {
			id := common.Int(string(item.([]byte)))
			user, _ := SearchUserByID(id)
			row = append(row, dto.ToUserDTO(user))
		}
		users[i] = row

		commentID := common.Int(strings.Split(string(key.([]byte)), ":")[2][7:])
		comment, _ := SearchCommentByID(commentID)
		comments[i] = dto.ToCommentDTO(comment)

		// 删除 hashmap
		_, _ = redisConn.Do("del", string(key.([]byte)))
	}

	return users, comments, nil
}

func GetStarNotifications(userID int, userType string) ([][]*dto.UserDTO, []*dto.StarDTO, error) {
	keys, err := redis.Values(redisConn.Do("keys", fmt.Sprintf("notice:user%d:star*:%s",
		userID, userType)))
	if err != nil {
		return nil, nil, err
	}

	stars := make([]*dto.StarDTO, len(keys))
	users := make([][]*dto.UserDTO, len(keys))
	for i, key := range keys {
		items, err := redis.Values(redisConn.Do("hkeys", string(key.([]byte))))
		if err != nil {
			return nil, nil, err
		}

		var row []*dto.UserDTO
		for _, item := range items {
			id := common.Int(string(item.([]byte)))
			user, _ := SearchUserByID(id)
			row = append(row, dto.ToUserDTO(user))
		}
		users[i] = row

		starID := common.Int(strings.Split(string(key.([]byte)), ":")[2][4:])
		star, _ := SearchStarByID(starID)
		stars[i] = dto.ToStarDTO(star)

		// 删除 hashmap
		_, _ = redisConn.Do("del", string(key.([]byte)))
	}

	return users, stars, nil
}

func GetCommentLikeNotifications(userID int) ([][]*dto.UserDTO, []*dto.CommentDTO, error) {
	return GetCommentNotifications(userID, "like_users")
}

func GetCommentReplyNotifications(userID int) ([][]*dto.UserDTO, []*dto.CommentDTO, error)  {
	return GetCommentNotifications(userID, "reply_users")
}

func GetStarLikeNotifications(userID int) ([][]*dto.UserDTO, []*dto.StarDTO, error)  {
	return GetStarNotifications(userID, "like_users")
}

func GetStarForwardNotifications(userID int) ([][]*dto.UserDTO, []*dto.StarDTO, error) {
	return GetStarNotifications(userID, "forward_users")
}
