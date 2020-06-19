package repositories

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"log"
	"magicMakeup/common"
	"magicMakeup/entities"
	"time"
)

func InsertStar(star *entities.Star) (*entities.Star, error) {
	res := mysqlConn.Create(star)

	if res.Error != nil {
		return nil, err
	}

	return star, nil
}

func UpdateStar(star *entities.Star) (*entities.Star, error) {
	res := mysqlConn.Save(star)

	if res.Error != nil {
		return nil, err
	}

	return star, nil
}

func SearchStarByID(id int) (*entities.Star, error) {
	star := &entities.Star{}
	res := mysqlConn.First(star, id)

	if res.Error != nil {
		return nil, err
	}

	return star, nil
}

func ForwardStar(star *entities.Star) (*entities.Star, error) {
	star, _ = UpdateStar(star)
	// 递增转发数
	originStar, _ := SearchStarByID(star.PreviousID)
	originStar.ForwardNum++
	_, _ = UpdateStar(originStar)

	// 记录转发动态 id
	_, err = redisConn.Do("hset", fmt.Sprintf("star%d:forward_stars", star.PreviousID), star.ID, time.Now().Unix())
	if err != nil {
		return nil, err
	}

	items, _ := redis.Values(redisConn.Do("hkeys", fmt.Sprintf("star%d:forward_stars", star.PreviousID)))
	for _, item := range items {
		id := common.Int(string(item.([]byte)))
		star, _ := SearchStarByID(id)
		star.ForwardNum = originStar.ForwardNum
		_, _ = UpdateStar(star)
	}

	return star, nil
}

func DeleteStar(star *entities.Star) (bool, error) {
	// 如果该动态是转发的
	if star.IsForward == 1 {
		return DeleteForwardStar(star)
	}
	// 如果该动态不是转发的
	return DeleteOriginalStar(star)
}

func DeleteOriginalStar(star *entities.Star) (bool, error) {
	// 删除原动态
	res := mysqlConn.Delete(star)
	if res.Error != nil {
		return false, err
	}
	err := common.DeleteDirectory(star.Images)
	if err != nil {
		log.Printf("delete directory failed: %v\n", err)
		return false, nil
	}

	/**
	items, err := redis.Values(redisConn.Do("hkeys", fmt.Sprintf("star%d:forward_stars", star.ID)))
	if err != nil {
		return false, err
	}
	// 删除转发动态
	for _, item := range items {
		id, _ := strconv.Atoi(string(item.([]byte)))
		star, _ := SearchStarByID(id)
		res = mysqlConn.Delete(star)
		if res.Error != nil {
			return false, err
		}
		err = common.DeleteDirectory(star.Images)
		if err != nil {
			log.Printf("delete directory failed: %v\n", err)
			return false, nil
		}
	}

	// 删除 redis 中的记录
	_, err = redisConn.Do("del", fmt.Sprintf("star%d:forward_stars", star.ID))
	if err != nil {
		return false, err
	}
	 */

	return true, nil
}

func DeleteForwardStar(star *entities.Star) (bool, error) {
	originalStar, _ := SearchStarByID(star.PreviousID)

	// 删除该动态
	mysqlConn.Delete(star)
	err := common.DeleteDirectory(star.Images)
	if err != nil {
		return false, err
	}

	// 删除 redis 中的对应条目
	_, err = redisConn.Do("hdel", fmt.Sprintf("star%d:forward_stars", originalStar.ID), star.ID)
	if err != nil {
		return false, err
	}

	// 原动态转发数减一
	originalStar.ForwardNum--
	_, err = UpdateStar(originalStar)
	if err != nil {
		return false, err
	}

	items, _ := redis.Values(redisConn.Do("hkeys", fmt.Sprintf("star%d:forward_stars", originalStar.ID)))
	for _, item := range items {
		id := common.Int(string(item.([]byte)))
		star, _ := SearchStarByID(id)
		star.ForwardNum = originalStar.ForwardNum
		_, _ = UpdateStar(star)
	}

	return true, err
}

func LikeStar(userID int, star *entities.Star) error {
	// 给一条转发的动态点赞本质是给原动态点赞
	if star.IsForward == 1 {
		originalStar, _ := SearchStarByID(star.PreviousID)
		return LikeStar(userID, originalStar)
	}
	// 点赞数加一
	star.LikeNum++
	_, _ = UpdateStar(star)

	// 更新转发的动态
	items, err := redis.Values(redisConn.Do("hkeys", fmt.Sprintf("star%d:forward_stars", star.ID)))
	if err != nil {
		return err
	}
	for _, item := range items {
		id := common.Int(string(item.([]byte)))
		star, _ := SearchStarByID(id)
		star.LikeNum++
		_, _ = UpdateStar(star)
	}

	// 记录点赞用户 id
	_, err = redisConn.Do("hset", fmt.Sprintf("star%d:like_users", star.ID), userID, time.Now().Unix())
	if err != nil {
		return err
	}

	return nil
}

func CancelLikeStar(userID int, star *entities.Star) error {
	if star.IsForward == 1 {
		originalStar, _ := SearchStarByID(star.PreviousID)
		return CancelLikeStar(userID, originalStar)
	}
	star.LikeNum--
	_, _ = UpdateStar(star)

	items, err := redis.Values(redisConn.Do("hkeys", fmt.Sprintf("star%d:forward_stars", star.ID)))
	if err != nil {
		return err
	}
	for _, item := range items {
		id := common.Int(string(item.([]byte)))
		star, _ := SearchStarByID(id)
		star.LikeNum--
		_, _ = UpdateStar(star)
	}

	_, err = redisConn.Do("hdel", fmt.Sprintf("star%d:like_users", star.ID), userID)
	if err != nil {
		return err
	}

	return nil
}

func GetStarList(userID int) ([]*entities.Star, error) {
	var stars []*entities.Star

	res := mysqlConn.Where("user_id=?", userID).Find(&stars)
	if res.Error != nil {
		return nil, err
	}

	return stars, nil
}

func GetLikeUserList(id int) ([]*entities.User, error) {
	items, err := redis.Values(redisConn.Do("hkeys", fmt.Sprintf("star%d:like_users", id)))
	if err != nil {
		return nil, err
	}

	var users []*entities.User
	for _, item := range items {
		userID := common.Int(string(item.([]byte)))
		user, _ := SearchUserByID(userID)
		users = append(users, user)
	}

	return users, nil
}

func GetForwardStarList(id int) ([]*entities.Star, error) {
	items, err := redis.Values(redisConn.Do("hkeys", fmt.Sprintf("star%d:forward_stars", id)))
	if err != nil {
		return nil, err
	}

	var stars []*entities.Star
	for _, item := range items {
		starID := common.Int(string(item.([]byte)))
		star, _ := SearchStarByID(starID)
		stars = append(stars, star)
	}

	return stars, nil
}
