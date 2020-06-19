package rabbitmq

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
)

var (
	ip string
	port string
	user string
	password string
	vhost string
	commentQueue string

	url string
)

var RabbitmqConn *RabbitMQ

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath("./config")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("read config failed: %v", err)
	}

	ip = viper.GetString("rabbitmq.ip")
	port = viper.GetString("rabbitmq.port")
	user = viper.GetString("rabbitmq.user")
	password = viper.GetString("rabbitmq.password")
	vhost = viper.GetString("rabbitmq.vhost")
	commentQueue = viper.GetString("rabbitmq.comment_queue")

	url = fmt.Sprintf("amqp://%s:%s@%s:%s/%s", user, password, ip, port, vhost)
	RabbitmqConn = NewRabbitMQ()
}
