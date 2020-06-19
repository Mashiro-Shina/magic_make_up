package handlers

import (
	"github.com/gin-gonic/gin"
	"log"
	"magicMakeup/common"
	"magicMakeup/dto"
	"magicMakeup/entities"
	"magicMakeup/repositories"
	"magicMakeup/response"
	"net/http"
	"time"
)

func HandlePublishStar(ctx *gin.Context) {
	userID := common.GetUserIdFromToken(ctx)
	content := ctx.PostForm("content")
	images, err := ctx.FormFile("images")
	publishTime := time.Now().Unix()

	// 暂时插入数据
	star := entities.NewStar(userID, content, "", publishTime, "0")
	_, err = repositories.InsertStar(star)
	if err != nil {
		log.Printf("insert star failed: %v\n", err)
		response.Response(ctx, http.StatusInternalServerError, 500, nil, "内部错误")
		return
	}
	if images != nil {
		// 生成图片保存路径
		imgPath, err := common.GenerateStarImagePath(repositories.StarImgsDir, star.ID, 1)
		_ = ctx.SaveUploadedFile(images, imgPath)
		star.Images = imgPath

		// 更新图片路径
		_, err = repositories.UpdateStar(star)
		if err != nil {
			log.Printf("insert star failed: %v\n", err)
			response.Response(ctx, http.StatusInternalServerError, 500, nil, "内部错误")
			return
		}

	}
	response.Response(ctx, http.StatusOK, 200, gin.H{
		"star": dto.ToStarDTO(star),
	}, "发布动态成功")
}

func HandleForwardStar(ctx *gin.Context) {
	userID := common.GetUserIdFromToken(ctx)
	content := ctx.PostForm("content")
	images, err := ctx.FormFile("images")
	publishTime := time.Now().Unix()

	// 暂时插入数据
	star := entities.NewStar(userID, content, "", publishTime, "1")
	_, err = repositories.InsertStar(star)
	if err != nil {
		log.Printf("insert star failed: %v\n", err)
		response.Response(ctx, http.StatusInternalServerError, 500, nil, "内部错误")
		return
	}
	if images != nil {
		// 生成图片保存路径
		imgPath, err := common.GenerateStarImagePath(repositories.StarImgsDir, star.ID, 1)
		_ = ctx.SaveUploadedFile(images, imgPath)
		star.Images = imgPath

		// 更新图片路径
		_, err = repositories.UpdateStar(star)
		if err != nil {
			log.Printf("insert star failed: %v\n", err)
			response.Response(ctx, http.StatusInternalServerError, 500, nil, "内部错误")
			return
		}

	}
	starID := ctx.Param("starID")
	// 获取原动态
	originalStar, err := repositories.SearchStarByID(common.Int(starID))
	if err != nil {
		log.Printf("search star failed: %v\n", err)
		response.Response(ctx, http.StatusInternalServerError, 500, nil, "内部错误")
		return
	}

	if originalStar.IsForward == 1 {
		originalStar, _ = repositories.SearchStarByID(originalStar.PreviousID)
	}

	// 设置原动态 id
	star.PreviousID = originalStar.ID
	star.ForwardNum = originalStar.ForwardNum + 1
	star.LikeNum = originalStar.LikeNum
	star.CommentNum = originalStar.CommentNum

	_, err = repositories.ForwardStar(star)
	if err != nil {
		response.Response(ctx, http.StatusInternalServerError, 500, nil, "内部错误")
		return
	}

	response.Response(ctx, http.StatusOK, 200, gin.H{
		"star": dto.ToStarDTO(star),
	}, "转发动态成功")
}

func HandleUpdateStar(ctx *gin.Context) {
	starID := ctx.Param("starID")
	content := ctx.PostForm("content")
	images, _ := ctx.FormFile("images")
	userID := common.GetUserIdFromToken(ctx)

	star, _ := repositories.SearchStarByID(common.Int(starID))
	if star.UserID != userID {
		response.Response(ctx, http.StatusBadRequest, 400, nil, "没有修改权限")
		return
	}

	star.Content = content
	_ = common.DeleteDirectory(star.Images)
	imgPath, _ := common.GenerateStarImagePath(repositories.StarImgsDir, star.ID, 1)
	_ = ctx.SaveUploadedFile(images, imgPath)
	star.Images = imgPath

	_, err := repositories.UpdateStar(star)
	if err != nil {
		log.Printf("update star failed: %v\n", err)
		response.Response(ctx, http.StatusInternalServerError, 500, nil, "内部错误")
		return
	}

	response.Response(ctx, http.StatusOK, 200, gin.H{
		"star": dto.ToStarDTO(star),
	}, "更新动态成功")
}

func HandleDeleteStar(ctx *gin.Context) {
	starID := ctx.Param("starID")
	star, _ := repositories.SearchStarByID(common.Int(starID))
	if star.UserID != common.GetUserIdFromToken(ctx) {
		response.Response(ctx, http.StatusBadRequest, 400, nil, "权限不足")
		return
	}

	ok, err := repositories.DeleteStar(star)

	if !ok || err != nil {
		log.Printf("delete star failed: %v\n", err)
		response.Response(ctx, http.StatusInternalServerError, 500, nil, "内部错误")
		return
	}

	response.Response(ctx, http.StatusOK, 200, nil, "删除动态成功")
}

func HandleLikeStar(ctx *gin.Context) {
	starID := ctx.Param("starID")
	star, err := repositories.SearchStarByID(common.Int(starID))
	if err != nil {
		response.Response(ctx, http.StatusInternalServerError, 500, nil, "内部错误")
		return
	}

	err = repositories.LikeStar(common.GetUserIdFromToken(ctx), star)
	if err != nil {
		response.Response(ctx, http.StatusInternalServerError, 500, nil, "内部错误")
		return
	}

	response.Response(ctx, http.StatusOK, 200, nil, "点赞成功")
}

func HandleCancelLikeStar(ctx *gin.Context) {
	starID := ctx.Param("starID")
	star, err := repositories.SearchStarByID(common.Int(starID))
	if err != nil {
		response.Response(ctx, http.StatusInternalServerError, 500, nil, "内部错误")
		return
	}

	err = repositories.CancelLikeStar(common.GetUserIdFromToken(ctx), star)
	if err != nil {
		response.Response(ctx, http.StatusInternalServerError, 500, nil, "内部错误")
		return
	}

	response.Response(ctx, http.StatusOK, 200, nil, "取消点赞成功")
}

func HandleGetStar(ctx *gin.Context) {
	starID := ctx.Param("starID")
	star, err := repositories.SearchStarByID(common.Int(starID))
	if err != nil {
		response.Response(ctx, http.StatusInternalServerError, 500, nil, "内部错误")
		return
	}

	if star == nil {
		response.Response(ctx, http.StatusBadRequest, 400, nil, "该动态已删除")
		return
	}

	user, _ := repositories.SearchUserByID(star.UserID)

	response.Response(ctx, http.StatusOK, 200, gin.H{
		"star": dto.ToStarDTO(star),
		"user": dto.ToUserDTO(user),
	}, "获取动态成功")
}

func HandleStarList(ctx *gin.Context) {
	userID := common.Int(ctx.Param("id"))

	stars, err := repositories.GetStarList(userID)
	if err != nil {
		response.Response(ctx, http.StatusInternalServerError, 500, nil, "内部错误")
		return
	}

	res := make([]*dto.StarDTO, len(stars))
	for i, star := range stars {
		res[i] = dto.ToStarDTO(star)
	}

	response.Response(ctx, http.StatusOK, 200, gin.H{
		"starList": res,
	}, "获取动态列表成功")
}

func HandleForwardStarList(ctx *gin.Context) {
	starID := ctx.Param("starID")

	stars, err := repositories.GetForwardStarList(common.Int(starID))
	if err != nil {
		response.Response(ctx, http.StatusInternalServerError, 500, nil, "内部错误")
		return
	}

	var res []*dto.StarDTO
	for _, star := range stars {
		res = append(res, dto.ToStarDTO(star))
	}
	response.Response(ctx, http.StatusOK, 200, gin.H{
		"starList": res,
	}, "获取转发动态列表成功")
}

func HandleLikeUserList(ctx *gin.Context) {
	starID := ctx.Param("starID")

	users, err := repositories.GetLikeUserList(common.Int(starID))
	if err != nil {
		response.Response(ctx, http.StatusInternalServerError, 500, nil, "内部错误")
		return
	}

	var res []*dto.UserDTO
	for _, user := range users {
		res = append(res, dto.ToUserDTO(user))
	}
	response.Response(ctx, http.StatusOK, 200, gin.H{
		"userList": res,
	}, "获取点赞用户列表成功")
}

