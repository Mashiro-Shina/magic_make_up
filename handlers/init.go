package handlers

import (
	"github.com/spf13/viper"
	"log"
)

var (
	avatarsDir string
	defaultAvatar string
)

func initConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath("./config")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("read config failed: %v", err)
	}
}

func init() {
	avatarsDir = viper.GetString("avatars_dir")
	defaultAvatar = viper.GetString("avatars_dir") + "default.jpg"
}
