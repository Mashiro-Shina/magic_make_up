package common

import (
	"bufio"
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"os"
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
