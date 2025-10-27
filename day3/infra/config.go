package infra

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	RPCUrl         string
	WsUrl          string
	RedisAddr      string
	ERC20Contracts string

	SevConfig *SeverConfig
}
type SeverConfig struct {
	Ip   string
	Port uint64
}

var TheAppConfig *AppConfig

func LoadConfig() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("无法加载 .env 文件，使用默认配置")
	}

	TheAppConfig = &AppConfig{
		RPCUrl:         os.Getenv("RPC_URL"),
		WsUrl:          os.Getenv("WS_URL"),
		RedisAddr:      os.Getenv("REDIS_ADDR"),
		ERC20Contracts: os.Getenv("ERC20_CONTRACTS"),
	}
	port, err := strconv.ParseUint(os.Getenv("SERVER_PORT"), 10, 32)
	if err != nil {
		log.Println("无法加载 .env 文件，使用默认配置")
	}
	sevConfig := SeverConfig{
		Ip:   os.Getenv("SERVER_IP"),
		Port: port,
	}
	TheAppConfig.SevConfig = &sevConfig
}
