package handlers

import (
	"github.com/gin-gonic/gin"
	"magicMakeup/common"
	"magicMakeup/repositories"
	"magicMakeup/response"
	"net/http"
)

func HandleGetCommentLikeNotifications(ctx *gin.Context) {
	userID := common.GetUserIdFromToken(ctx)

	users, comments, err := repositories.GetCommentLikeNotifications(userID)
	if err != nil {
		response.Response(ctx, http.StatusInternalServerError, 500, nil, "内部错误")
		return
	}

	response.Response(ctx, http.StatusOK, 200, gin.H{
		"comments": comments,
		"users": users,
	}, "获取评论点赞通知成功")
}

func HandleGetCommentReplyNotifications(ctx *gin.Context) {
	userID := common.GetUserIdFromToken(ctx)

	users, comments, err := repositories.GetCommentReplyNotifications(userID)
	if err != nil {
		response.Response(ctx, http.StatusInternalServerError, 500, nil, "内部错误")
		return
	}

	response.Response(ctx, http.StatusOK, 200, gin.H{
		"comments": comments,
		"users": users,
	}, "获取评论回复通知成功")
}

func HandleGetStarLikeNotifications(ctx *gin.Context) {
	userID := common.GetUserIdFromToken(ctx)

	users, stars, err := repositories.GetStarLikeNotifications(userID)
	if err != nil {
		response.Response(ctx, http.StatusInternalServerError, 500, nil, "内部错误")
		return
	}

	response.Response(ctx, http.StatusOK, 200, gin.H{
		"stars": stars,
		"users": users,
	}, "获取动态点赞通知成功")
}

func HandleGetStarForwardNotifications(ctx *gin.Context) {
	userID := common.GetUserIdFromToken(ctx)

	users, stars, err := repositories.GetStarForwardNotifications(userID)
	if err != nil {
		response.Response(ctx, http.StatusInternalServerError, 500, nil, "内部错误")
		return
	}

	response.Response(ctx, http.StatusOK, 200, gin.H{
		"stars": stars,
		"users": users,
	}, "获取动态转发通知成功")
}
