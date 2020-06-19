package common

import (
	"bufio"
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"os"
	"strconv"
	"strings"
)

// 从 token　中读取用户 id
func GetUserIdFromToken(ctx *gin.Context) int {
	tokenString := ctx.GetHeader("Authorization")[7:]
	_, claims, _ := ParseToken(tokenString)
	return claims.UserID
}

//　获取两个列表的公共元素
func GetCommon(list1 []string, list2 []string) []string {
	common := make(map[string]bool, len(list1))
	for _, item := range list1 {
		common[item] = true
	}

	var res []string
	for _, item := range list2 {
		if common[item] {
			res = append(res, item)
		}
	}

	return res
}

// 读取图片并编码为 base64 格式
func ReadAndEncodingImage(path string) (string, error) {
	image, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer image.Close()

	fileInfo, _ := image.Stat()
	fileSize := fileInfo.Size()
	buffer := make([]byte, fileSize)

	fReader := bufio.NewReader(image)
	fReader.Read(buffer)

	base64Str := base64.StdEncoding.EncodeToString(buffer)
	return base64Str, nil
}

func GenerateStarImagePath(basePath string, starID int, imageID int) (string, error) {
	subDir := basePath + strconv.Itoa(starID)
	_, err := os.Stat(subDir)
	if err != nil {
		if os.IsNotExist(err) {
			os.Mkdir(subDir, os.ModePerm)
		} else {
			return "", err
		}
	}

	return subDir + "/" + strconv.Itoa(imageID) + ".jpg", nil
}


func DeleteDirectory(path string) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return nil
	}

	err = os.RemoveAll(path)
	if err != nil {
		return err
	}

	dirs := strings.Split(path, "/")
	err = os.Remove(strings.Join(dirs[:len(dirs)-1], "/"))
	if err != nil {
		return err
	}

	return nil
}

func Int(s string) int {
	v, _ := strconv.Atoi(s)
	return v
}
