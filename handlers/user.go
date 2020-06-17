package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"log"
	"magicMakeup/common"
	"magicMakeup/dto"
	"magicMakeup/entities"
	"magicMakeup/repositories"
	"magicMakeup/response"
	"net/http"
	"strconv"
)

func HandleRegister(ctx *gin.Context) {
	name := ctx.PostForm("name")
	password := ctx.PostForm("password")
	phone := ctx.PostForm("phone")

	// 将用户密码加密
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		response.Response(ctx, http.StatusInternalServerError, 500, nil, "内部错误")
		log.Printf("encoding password failed: %s", err)
		return
	}

	user := entities.NewUser(phone, string(hashedPwd), name, "", defaultAvatar)
	_, err = repositories.InsertUser(user)
	if err != nil {
		// 违背 unique
		if e, ok := err.(*mysql.MySQLError); ok && e.Number == 1062 {
			response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "该手机号已被注册")
			log.Print("")
			return
		}
		response.Response(ctx, http.StatusInternalServerError, 500, nil, "内部错误")
		log.Printf("insert user failed: %v", err)
		return
	}

	response.Response(ctx, http.StatusCreated, 201, nil, "注册成功")
}

func HandleLogin(ctx *gin.Context) {
	phone := ctx.PostForm("phone")
	password := ctx.PostForm("password")

	user, err := repositories.SearchUserByPhone(phone)
	if err != nil {
		if err.Error() == "record not found" {
			response.Response(ctx, http.StatusBadRequest, 400, nil, "该手机号尚未注册")
		} else {
			response.Response(ctx, http.StatusInternalServerError, 500, nil, "内部错误")
		}
		log.Printf("login failed: %s", err)
		return
	}

	// 验证密码是否正确
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		response.Response(ctx, http.StatusBadRequest, 400, nil, "密码错误")
		log.Print("login failed: password is not correct")
		return
	}

	// 生成 token
	token, err := common.ReleaseToken(user)
	if err != nil {
		response.Response(ctx, http.StatusInternalServerError, 500, nil, "内部错误")
		log.Printf("release token failed: %v", err)
		return
	}

	response.Response(ctx, http.StatusOK, 200, gin.H{
		"token": token,
		"userID": user.ID,
	}, "登录成功")
}

func HandleAccount(ctx *gin.Context) {
	userID := common.GetUserIdFromToken(ctx)
	user, _ := repositories.SearchUserByID(userID)
	response.Response(ctx, http.StatusOK, 200, gin.H{
		"data": dto.ToUserDTO(user),
	}, "user profile")
}

func HandleUpdate(ctx *gin.Context) {
	userID := common.GetUserIdFromToken(ctx)
	phone := ctx.PostForm("phone")
	password := ctx.PostForm("password")
	encodedPwd, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	signature := ctx.PostForm("signature")
	name := ctx.PostForm("name")

	user, _ := repositories.SearchUserByID(userID)
	user.Phone = phone
	user.Password = string(encodedPwd)
	user.Signature = signature
	user.Name = name
	ok, err := repositories.UpdateUserProfile(user)

	if !ok || err != nil {
		response.Response(ctx, http.StatusInternalServerError, 500, nil, "内部错误")
		log.Printf("update user profile failed: %v", err)
		return
	}

	response.Response(ctx, http.StatusOK, 200, nil, "更新信息成功")
}

func HandleFollowUser(ctx *gin.Context) {
	from := ctx.Param("from")
	to := ctx.Param("to")

	ok, err := repositories.FollowUser(from, to)
	if !ok || err != nil {
		log.Printf("%v\n", err)
		response.Response(ctx, http.StatusInternalServerError, 500, nil, "内部错误")
	}

	response.Response(ctx, http.StatusOK, 200, nil, "关注成功")
}

func HandleUnfollowUser(ctx *gin.Context) {
	from := ctx.Param("from")
	to := ctx.Param("to")

	ok, err := repositories.UnFollowUser(from, to)
	if !ok || err != nil {
		log.Printf("%v\n", err)
		response.Response(ctx, http.StatusInternalServerError, 500, nil, "内部错误")
	}

	response.Response(ctx, http.StatusOK, 200, nil, "取消关注成功")
}

func HandleFollowingList(ctx *gin.Context) {
	userID := ctx.Param("id")
	items, err := repositories.GetFollowingList(userID)

	if err != nil {
		response.Response(ctx, http.StatusInternalServerError, 500, nil, "内部错误")
		return
	}

	userList := transformToUserList(items)
	response.Response(ctx, http.StatusOK, 200, gin.H{
		"userList": userList,
	}, "获取关注列表成功")
}

func HandleFollowersList(ctx *gin.Context) {
	userID := ctx.Param("id")
	items, err := repositories.GetFollowersList(userID)

	if err != nil {
		response.Response(ctx, http.StatusInternalServerError, 500, nil, "内部错误")
		return
	}

	userList := transformToUserList(items)
	response.Response(ctx, http.StatusOK, 200, gin.H{
		"userList": userList,
	}, "获取粉丝列表成功")
}

func HandleCommonFollowingList(ctx *gin.Context) {
	id1 := ctx.Param("id1")
	id2 := ctx.Param("id2")
	list, err := repositories.GetCommonFollowingList(id1, id2)

	if err != nil {
		log.Printf("get common following list failed: %v", err)
		response.Response(ctx, http.StatusInternalServerError, 500, nil, "获取共同关注列表失败")
		return
	}

	userList := transformToUserList(list)
	response.Response(ctx, http.StatusOK, 200, gin.H{
		"userList": userList,
	}, "获取共同关注列表成功")
}

func HandleMutualFollowingList(ctx *gin.Context) {
	userID := ctx.Param("id")
	list, err := repositories.GetMutualFollowingList(userID)

	if err != nil {
		log.Printf("get mutual following list failed: %v", err)
		response.Response(ctx, http.StatusInternalServerError, 500, nil, "获取互相关注列表失败")
		return
	}

	userList := transformToUserList(list)
	response.Response(ctx, http.StatusOK, 200, gin.H{
		"userList": userList,
	}, "获取互相关注列表成功")
}

func transformToUserList(items []string) []*dto.UserDTO {
	var res []*dto.UserDTO
	for _, item := range items {
		id, _ := strconv.Atoi(item)
		user, _ := repositories.SearchUserByID(id)
		res = append(res, dto.ToUserDTO(user))
	}
	return res
}

func HandleUpdateAvatar(ctx *gin.Context) {
	avatar, err := ctx.FormFile("avatar")
	if err != nil {
		log.Printf("接收图片失败: %v", err)
		response.Response(ctx, http.StatusInternalServerError, 500, nil, "接收图片失败")
		return
	}

	userID := ctx.Param("id")
	err = ctx.SaveUploadedFile(avatar, avatarsDir + userID + ".jpg")
	if err != nil {
		log.Printf("保存头像失败: %v", err)
		response.Response(ctx, http.StatusInternalServerError, 500, nil, "保存头像失败")
		return
	}

	id, _ := strconv.Atoi(userID)
	err = repositories.UpdateAvatar(id, avatarsDir + userID + ".jpg")
	if err != nil {
		log.Printf("更新头像失败: %v", err)
		response.Response(ctx, http.StatusInternalServerError, 500, nil, "更换头像失败")
		return
	}

	response.Response(ctx, http.StatusOK, 200, nil, "更换头像成功")
}
