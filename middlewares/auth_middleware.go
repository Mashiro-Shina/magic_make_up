package middlewares

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"magicMakeup/common"
	"magicMakeup/dto"
	"magicMakeup/repositories"
	"magicMakeup/response"
	"net/http"
	"strings"
)

func AuthMiddleWare() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString := ctx.GetHeader("Authorization")

		// 没有 token 或者 token 不合法
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
			response.Response(ctx, http.StatusUnauthorized, 401, nil, "权限不足")
			ctx.Abort()
			return
		}

		tokenString = tokenString[7:]
		// 解析 token
		token, claims, err := common.ParseToken(tokenString)
		if err != nil || !token.Valid {
			response.Response(ctx, http.StatusBadRequest, 401, nil, "token 不合法，权限不足")
			ctx.Abort()
			return
		}

		// 验证用户是否存在，因为即使用户已注销，只要 token 未过期就仍然可以使用
		user, err := repositories.SearchUserByID(claims.UserID)
		if err != nil {
			response.Response(ctx, http.StatusBadRequest, 401, nil, "用户已注销，权限不足")
			return
		}

		ctx.Set(fmt.Sprintf("user%d", user.ID), dto.ToUserDTO(user))
		ctx.Next()
	}
}
