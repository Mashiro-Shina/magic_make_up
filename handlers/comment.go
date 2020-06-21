package handlers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"magicMakeup/common"
	"magicMakeup/dto"
	"magicMakeup/entities"
	"magicMakeup/rabbitmq"
	"magicMakeup/repositories"
	"magicMakeup/response"
	"net/http"
	"time"
)

func HandlePublishComment(ctx *gin.Context) {
	commentID := <- commentIDChannel
	userID := common.GetUserIdFromToken(ctx)
	starID := common.Int(ctx.Param("starID"))
	replyID := common.Int(ctx.Param("replyID"))
	content := ctx.PostForm("content")
	publishTime := time.Now().Unix()

	comment := entities.NewComment(commentID, userID, starID, replyID, content, publishTime)

	jsonComment, err := json.Marshal(comment)
	if err != nil {
		log.Printf("encoding comment failed: %v\n", err)
		response.Response(ctx, http.StatusInternalServerError, 500, nil, "内部错误")
		return
	}

	rabbitmq.RabbitmqConn.PublishComment(jsonComment)

	//repositories.InsertComment(comment)
	response.Response(ctx, http.StatusOK, 200, gin.H{
		"comment": dto.ToCommentDTO(comment),
	}, "发布评论成功")
}

func HandleDeleteComment(ctx *gin.Context) {
	commentID := ctx.Param("commentID")

	err := repositories.DeleteComment(common.Int(commentID))
	if err != nil {
		response.Response(ctx, http.StatusInternalServerError, 500, nil, "内部错误")
		return
	}

	response.Response(ctx, http.StatusOK, 200, nil, "删除评论成功")
}

func HandleGetComment(ctx *gin.Context) {
	commentID := ctx.Param("commentID")
	comment, err := repositories.SearchCommentByID(common.Int(commentID))
	if err != nil {
		if err.Error() == "record not found" {
			response.Response(ctx, http.StatusBadRequest, 400, nil, "该评论已删除")
			return
		}
		response.Response(ctx, http.StatusInternalServerError, 500, nil, "内部错误")
		return
	}

	userID := comment.UserID
	user, err := repositories.SearchUserByID(userID)
	if err != nil {
		response.Response(ctx, http.StatusInternalServerError, 200, nil, "内部错误")
		return
	}

	response.Response(ctx, http.StatusOK, 200, gin.H{
		"comment": dto.ToCommentDTO(comment),
		"user": dto.ToUserDTO(user),
	}, "获取评论成功")
}

func HandleGetReplyCommentList(ctx *gin.Context) {
	commentID := ctx.Param("commentID")

	comments, err := repositories.GetReplyCommentList(common.Int(commentID))
	if err != nil {
		response.Response(ctx, http.StatusInternalServerError, 500, nil, "内部错误")
		return
	}

	users := make([]*dto.UserDTO, len(comments))
	for i, comment := range comments {
		user, _ := repositories.SearchUserByID(comment.UserID)
		users[i] = dto.ToUserDTO(user)
	}

	response.Response(ctx, http.StatusOK, 200, gin.H{
		"comments": comments,
		"users": users,
	}, "获取回复列表成功")
}

func HandleLikeComment(ctx *gin.Context) {
	userID := common.GetUserIdFromToken(ctx)
	commentID := ctx.Param("commentID")

	err := repositories.LikeComment(userID, common.Int(commentID))
	if err != nil {
		response.Response(ctx, http.StatusInternalServerError, 500, nil, "内部错误")
		return
	}

	response.Response(ctx, http.StatusOK, 200, nil, "点赞成功")
}

func HandleCancelLikeComment(ctx *gin.Context) {
	userID := common.GetUserIdFromToken(ctx)
	commentID := ctx.Param("commentID")

	err := repositories.CancelLikeComment(userID, common.Int(commentID))
	if err != nil {
		response.Response(ctx, http.StatusInternalServerError, 500, nil, "内部错误")
		return
	}

	response.Response(ctx, http.StatusOK, 200, nil, "取消点赞成功")
}

func HandleCommentList(ctx *gin.Context) {
	starID := ctx.Param("starID")

	comments, err := repositories.GetCommentList(common.Int(starID))
	if err != nil {
		response.Response(ctx, http.StatusInternalServerError, 500, nil, "内部错误")
		return
	}

	users := make([]*dto.UserDTO, len(comments))
	commentDTOs := make([]*dto.CommentDTO, len(comments))
	replyComments := make([][]*dto.CommentDTO, len(comments))
	replyUsers := make([][]*dto.UserDTO, len(comments))
	for i, comment := range comments {
		user, _ := repositories.SearchUserByID(comment.UserID)
		users[i] = dto.ToUserDTO(user)
		commentDTOs[i] = dto.ToCommentDTO(comment)

		replys, _ := repositories.GetReplyCommentList(comment.ID)
		for _, reply := range replys {
			replyComments[i] = append(replyComments[i], dto.ToCommentDTO(reply))
			replyUser, _ := repositories.SearchUserByID(reply.UserID)
			replyUsers[i] = append(replyUsers[i], dto.ToUserDTO(replyUser))
		}
	}

	response.Response(ctx, http.StatusOK, 200, gin.H{
		"comments": commentDTOs,
		"users": users,
		"replyComments": replyComments,
		"replyUsers": replyUsers,
	}, "获取评论列表成功")
}

func HandleLikeCommentUserList(ctx *gin.Context) {
	commentID := ctx.Param("commentID")

	users, err := repositories.GetLikeCommentList(common.Int(commentID))
	if err != nil {
		response.Response(ctx, http.StatusInternalServerError, 500, nil, "内部错误")
		return
	}

	res := make([]*dto.UserDTO, len(users))
	for i, user := range users {
		res[i] = dto.ToUserDTO(user)
	}

	response.Response(ctx, http.StatusOK, 200, gin.H{
		"users": res,
	}, "获取点赞用户列表成功")
}

