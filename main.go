package main

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"log"
	"magicMakeup/handlers"
	"magicMakeup/middlewares"
)

var (
	ip string
	port string
)

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath("./config")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("read config failed: %v", err)
	}

	ip = viper.GetString("server.ip")
	port = viper.GetString("server.port")
}

func main() {
	engine := gin.Default()

	engine.POST("/register", handlers.HandleRegister)
	engine.POST("/login", handlers.HandleLogin)
	engine.GET("/account", middlewares.AuthMiddleWare(), handlers.HandleAccount)
	engine.POST("/account/update", middlewares.AuthMiddleWare(), handlers.HandleUpdate)
	engine.GET("/follow/:from/:to", middlewares.AuthMiddleWare(),  handlers.HandleFollowUser)
	engine.GET("/unfollow/:from/:to", middlewares.AuthMiddleWare(), handlers.HandleUnfollowUser)
	engine.GET("/common_following/:id1/:id2",middlewares.AuthMiddleWare(),  handlers.HandleCommonFollowingList)

	userGroup := engine.Group("/user")
	{
		userGroup.GET("/:id/fans", handlers.HandleFollowersList)
		userGroup.GET("/:id/following", handlers.HandleFollowingList)
		userGroup.GET("/:id/mutual_following", middlewares.AuthMiddleWare(), handlers.HandleMutualFollowingList)
		userGroup.POST("/:id/update_avatar",middlewares.AuthMiddleWare(),  handlers.HandleUpdateAvatar)
	}

	_ = engine.Run(ip + ":" + port)
}
