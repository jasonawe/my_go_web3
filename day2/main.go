package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

var (
	ctx         = context.Background()
	redisClient *redis.Client
	upgrader    = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true }, // 允许跨域
	}
)

const ERC20_ABI = `[
    {"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"from","type":"address"},{"indexed":true,"internalType":"address","name":"to","type":"address"},{"indexed":false,"internalType":"uint256","name":"value","type":"uint256"}],"name":"Transfer","type":"event"},
    {"constant":true,"inputs":[],"name":"decimals","outputs":[{"name":"","type":"uint8"}],"type":"function"},
    {"constant":true,"inputs":[],"name":"symbol","outputs":[{"name":"","type":"string"}],"type":"function"}
]`

func main() {
	_ = godotenv.Load()
	rpcURL := os.Getenv("RPC_URL")
	contracts := strings.Split(os.Getenv("ERC20_CONTRACTS"), ",")

	redisClient = redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	})
	defer redisClient.Close()

	r := gin.Default()

	// WebSocket 路由
	r.GET("/ws", func(c *gin.Context) {
		handleWS(c.Writer, c.Request)
	})

	// 静态文件挂载在 /static
	r.Static("/static", "./static")

	// 启动 ERC20 监听
	for _, c := range contracts {
		go listenTransfer(rpcURL, common.HexToAddress(c))
	}

	log.Println("🚀 Web3 Dashboard running on http://localhost:8080/static/index.html")
	r.Run(":8080")
}

func listenTransfer(rpcURL string, contractAddr common.Address) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		log.Fatalf("❌ RPC 连接失败: %v", err)
	}
	defer client.Close()

	parsedABI, _ := abi.JSON(strings.NewReader(ERC20_ABI))
	decimals := getDecimals(client, contractAddr, parsedABI)
	symbol := getSymbol(client, contractAddr, parsedABI)

	transferTopic := common.HexToHash("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef")
	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddr},
		Topics:    [][]common.Hash{{transferTopic}},
	}

	logsCh := make(chan types.Log)
	sub, err := client.SubscribeFilterLogs(context.Background(), query, logsCh)
	if err != nil {
		log.Fatalf("❌ 订阅失败: %v", err)
	}

	for {
		select {
		case err := <-sub.Err():
			log.Printf("❌ 订阅错误: %v", err)
		case vLog := <-logsCh:
			var transferEvent struct {
				From  common.Address
				To    common.Address
				Value *big.Int
			}
			_ = parsedABI.UnpackIntoInterface(&transferEvent, "Transfer", vLog.Data)
			transferEvent.From = common.HexToAddress(vLog.Topics[1].Hex())
			transferEvent.To = common.HexToAddress(vLog.Topics[2].Hex())

			divisor := new(big.Float).SetFloat64(float64(new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil).Int64()))
			realValue := new(big.Float).Quo(new(big.Float).SetInt(transferEvent.Value), divisor)

			msg := fmt.Sprintf("%s | From: %s To: %s Value: %f %s",
				contractAddr.Hex(),
				transferEvent.From.Hex(),
				transferEvent.To.Hex(),
				realValue,
				symbol,
			)
			log.Println(msg)
			redisClient.LPush(ctx, "transfer_events", msg)
		}
	}
}

func getDecimals(client *ethclient.Client, contract common.Address, parsedABI abi.ABI) uint8 {
	data, _ := parsedABI.Pack("decimals")
	res, _ := client.CallContract(ctx, ethereum.CallMsg{To: &contract, Data: data}, nil)
	var decimals uint8
	_ = parsedABI.UnpackIntoInterface(&decimals, "decimals", res)
	return decimals
}

func getSymbol(client *ethclient.Client, contract common.Address, parsedABI abi.ABI) string {
	data, _ := parsedABI.Pack("symbol")
	res, _ := client.CallContract(ctx, ethereum.CallMsg{To: &contract, Data: data}, nil)
	var symbol string
	_ = parsedABI.UnpackIntoInterface(&symbol, "symbol", res)
	return symbol
}

func handleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket Upgrade失败:", err)
		return
	}
	defer conn.Close()

	for {
		// 阻塞直到有数据
		result, err := redisClient.BLPop(ctx, time.Second, "transfer_events").Result() // 1秒超时保持活跃
		if err == redis.Nil {
			continue
		} else if err != nil {
			log.Println("Redis error:", err)
			continue
		}

		err = conn.WriteMessage(websocket.TextMessage, []byte(result[1]))
		if err != nil {
			log.Println("WebSocket写入失败:", err)
			break
		}
	}
}
