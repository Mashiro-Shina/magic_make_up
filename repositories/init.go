package repositories

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/spf13/viper"
	"log"
)

var mysqlConn *gorm.DB
var err error

var (
	mysqlIP string
	mysqlPort string
	mysqlUser string
	mysqlPassword string
	mysqlDatabase string

	redisIP string
	redisPort string
)

var redisConn redis.Conn

func initConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath("./config")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("read config failed: %v", err)
	}
}

func initMysql() {
	mysqlIP = viper.GetString("mysql.ip")
	mysqlPort = viper.GetString("mysql.port")
	mysqlUser = viper.GetString("mysql.user")
	mysqlPassword = viper.GetString("mysql.password")
	mysqlDatabase = viper.GetString("mysql.database")

	mysqlConn, err = gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		mysqlUser, mysqlPassword, mysqlIP, mysqlPort, mysqlDatabase))
	if err != nil {
		log.Fatalf("connect mysql failed: %v", err)
	}
}

func initRedis() {
	redisIP = viper.GetString("redis.ip")
	redisPort = viper.GetString("redis.port")

	redisConn, err = redis.Dial("tcp", redisIP + ":" + redisPort)
	if err != nil {
		log.Fatalf("connect redis failed: %v", err)
	}
}

func init() {
	initConfig()
	initMysql()
	initRedis()
}
