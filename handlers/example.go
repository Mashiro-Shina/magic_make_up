package handlers

import (
	"github.com/gin-gonic/gin"
	"magicMakeup/common"
	"magicMakeup/dto"
	"magicMakeup/repositories"
	"magicMakeup/response"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func HandleGetExamples(ctx *gin.Context) {
	stars, err := repositories.GetAllStars()
	if err != nil {
		response.Response(ctx, http.StatusInternalServerError, 500, nil, "内部错误")
		return
	}

	var starDTOs []*dto.StarDTO

	err = filepath.Walk(repositories.ExampleImgsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		parts := strings.Split(path, "/")
		if len(parts) != 3 {
			return nil
		}
		description := parts[1]
		image, _ := common.ReadAndEncodingImage("./" + path)
		starDTOs = append(starDTOs, &dto.StarDTO{
			Content: description,
			Images: image,
		})

		return nil
	})

	if err != nil {
		response.Response(ctx, http.StatusInternalServerError, 500, nil, "内部错误")
		return
	}

	for _, star := range stars {
		starDTOs = append(starDTOs, dto.ToStarDTO(star))
	}

	response.Response(ctx, http.StatusOK, 200, gin.H{
		"examples": starDTOs,
	}, "获取示例成功")
}
