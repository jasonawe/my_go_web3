package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
)

func main() {
	// 加载 .env
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ .env 文件未找到，将使用环境变量")
	}

	rpc := os.Getenv("RPC_URL")
	if rpc == "" {
		log.Fatal("❌ RPC_URL 未配置")
	}

	retryCount := 3
	if r := os.Getenv("RETRY"); r != "" {
		if n, err := strconv.Atoi(r); err == nil {
			retryCount = n
		}
	}

	timeoutSec := 10
	if t := os.Getenv("TIMEOUT"); t != "" {
		if n, err := strconv.Atoi(t); err == nil {
			timeoutSec = n
		}
	}

	var client *ethclient.Client
	var err error

	for i := 0; i < retryCount; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutSec)*time.Second)
		defer cancel()

		client, err = ethclient.DialContext(ctx, rpc)
		if err != nil {
			log.Printf("⚠️ 连接 RPC 失败，第 %d 次重试: %v\n", i+1, err)
			time.Sleep(2 * time.Second)
			continue
		}
		break
	}

	if client == nil {
		log.Fatal("❌ RPC 连接失败，请检查网络或 RPC URL")
	}

	// 查询最新区块
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutSec)*time.Second)
	defer cancel()
	block, err := client.BlockNumber(ctx)
	if err != nil {
		log.Fatal("❌ 获取区块失败:", err)
	}

	fmt.Println("✅ RPC 连接成功！当前区块号:", block)
}
