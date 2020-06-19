package repositories

import (
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"log"
	"magicMakeup/common"
	"magicMakeup/entities"
	"strconv"
	"time"
)

func InsertUser(user *entities.User) (bool, error) {
	if err := mysqlConn.Create(user).Error; err != nil {
		return false, err
	}

	if mysqlConn.NewRecord(user) {
		return false, errors.New("insert new user failed")
	}

	return true, nil
}

func SearchUserByPhone(phone string) (*entities.User, error) {
	var user entities.User
	res := mysqlConn.Where("phone = ?", phone).First(&user)

	if res.Error != nil {
		return nil, res.Error
	}

	return &user, nil
}

func SearchUserByID(id int) (*entities.User, error) {
	var user entities.User
	res := mysqlConn.First(&user, id)

	if res.Error != nil {
		return nil, res.Error
	}

	return &user, nil
}

func UpdateUserProfile(user *entities.User) (bool, error) {
	res := mysqlConn.Save(user)

	if res.Error != nil {
		return false, err
	}

	return true, nil
}

func FollowUser(from string, to string) (bool, error) {
	_, err := redisConn.Do("hset", fmt.Sprintf("user%s:following", from), to, time.Now().Unix())
	if err != nil {
		return false, err
	}
	fromID, _ := strconv.Atoi(from)
	ok, err := AlterFollowing(fromID, 1)
	if !ok || err != nil {
		return false, err
	}

	_, err = redisConn.Do("hset", fmt.Sprintf("user%s:followers", to), from, time.Now().Unix())
	if err != nil {
		return false, err
	}
	toID, _ := strconv.Atoi(to)
	ok, err = AlterFollower(toID, 1)
	if !ok || err != nil {
		return false, err
	}

	return true, nil
}

func AlterFollowing(id int, alter int) (bool, error) {
	user, err := SearchUserByID(id)
	if err != nil {
		return false, err
	}

	res := mysqlConn.Model(user).Update("following_num", user.FollowingNum + alter)
	if res.Error != nil {
		return false, res.Error
	}

	return true, nil
}

func AlterFollower(id int, alter int) (bool, error) {
	user, err := SearchUserByID(id)
	if err != nil {
		return false, err
	}

	res := mysqlConn.Model(user).Update("followers_num", user.FollowersNum + alter)
	if res.Error != nil {
		return false, res.Error
	}

	return true, nil
}

func UnFollowUser(from string, to string) (bool, error) {
	_, err := redisConn.Do("hdel", fmt.Sprintf("user%s:following", from), to)
	if err != nil {
		return false, err
	}
	fromID, _ := strconv.Atoi(from)
	ok, err := AlterFollowing(fromID, -1)
	if !ok || err != nil {
		return false, err
	}

	_, err = redisConn.Do("hdel", fmt.Sprintf("user%s:followers", to), from)
	if err != nil {
		return false, err
	}

	toID, _ := strconv.Atoi(to)
	ok, err = AlterFollower(toID, -1)
	if !ok || err != nil {
		return false, err
	}

	return true, nil
}

func GetFollowingList(id string) ([]string, error) {
	items, err := redis.Values(redisConn.Do("hkeys", fmt.Sprintf("user%s:following", id)))
	if err != nil {
		log.Printf("get following list failed: %v", err)
		return nil, err
	}

	var res []string
	for _, item := range items {
		res = append(res, string(item.([]byte)))
	}

	return res, nil
}

func GetFollowersList(id string) ([]string, error) {
	items, err := redis.Values(redisConn.Do("hkeys", fmt.Sprintf("user%s:followers", id)))
	if err != nil {
		log.Printf("get followers list failed: %v", err)
		return nil, err
	}

	var res []string
	for _, item := range items {
		res = append(res, string(item.([]byte)))
	}

	return res, nil
}

func GetCommonFollowingList(id1, id2 string) ([]string, error) {
	list1, err := GetFollowingList(id1)
	if err != nil {
		return nil, err
	}
	list2, err := GetFollowingList(id2)
	if err != nil {
		return nil, err
	}

	res := common.GetCommon(list1, list2)
	return res, nil
}

func GetMutualFollowingList(id string) ([]string, error) {
	followingList, err := GetFollowingList(id)
	if err != nil {
		return nil, err
	}
	followersList, err := GetFollowersList(id)
	if err != nil {
		return nil, err
	}

	res := common.GetCommon(followingList, followersList)
	return res, nil
}

func UpdateAvatar(id int, avatar string) error {
	user, _ := SearchUserByID(id)
	res := mysqlConn.Model(user).Update("avatar", avatar)

	if res.Error != nil {
		return res.Error
	}

	return nil
}
