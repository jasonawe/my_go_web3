package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

const ERC20_ABI = `[{"constant":true,"inputs":[{"name":"_owner","type":"address"}],
"name":"balanceOf","outputs":[{"name":"balance","type":"uint256"}],"type":"function"},
{"constant":true,"inputs":[],"name":"decimals","outputs":[{"name":"","type":"uint8"}],"type":"function"},
{"constant":true,"inputs":[],"name":"symbol","outputs":[{"name":"","type":"string"}],"type":"function"}]`

func main() {
	// 创建 gin 实例
	r := gin.Default()
	// 加载 .env
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ .env 文件未找到，将使用环境变量")
	}

	rpcURL := os.Getenv("RPC_URL")
	// 初始化以太坊客户端 (Infura Sepolia)
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		log.Fatalf("❌ 无法连接以太坊网络: %v", err)
	}
	log.Println("✅ 已连接到 Sepolia 网络")

	// 定义 /balance/:address 接口
	r.GET("/balance/:address", func(c *gin.Context) {
		address := c.Param("address")
		if !common.IsHexAddress(address) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的以太坊地址"})
			return
		}

		acc := common.HexToAddress(address)
		balance, err := client.BalanceAt(context.Background(), acc, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ethValue := new(big.Float).Quo(new(big.Float).SetInt(balance), big.NewFloat(1e18))
		c.JSON(http.StatusOK, gin.H{
			"address": address,
			"balance": fmt.Sprintf("%f ETH", ethValue),
		})
	})
	r.GET("/token_balance/:contract/:address", func(c *gin.Context) {
		contractAddr := c.Param("contract")
		userAddr := c.Param("address")

		if !common.IsHexAddress(contractAddr) || !common.IsHexAddress(userAddr) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的地址"})
			return
		}

		contract := common.HexToAddress(contractAddr)
		user := common.HexToAddress(userAddr)

		// 解析 ABI
		parsedABI, err := abi.JSON(strings.NewReader(ERC20_ABI))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ABI解析失败"})
			return
		}

		// 创建合约调用对象
		// callOpts := &bind.CallOpts{Context: context.Background()}

		// 调用 symbol()
		symbolData, err := parsedABI.Pack("symbol")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "symbol 调用打包失败"})
			return
		}
		symbolResult, err := client.CallContract(context.Background(), ethereum.CallMsg{
			To:   &contract,
			Data: symbolData,
		}, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "symbol 调用失败"})
			return
		}

		var tokenSymbol string
		_ = parsedABI.UnpackIntoInterface(&tokenSymbol, "symbol", symbolResult)

		// 调用 decimals()
		decData, _ := parsedABI.Pack("decimals")
		decResult, _ := client.CallContract(context.Background(), ethereum.CallMsg{
			To:   &contract,
			Data: decData,
		}, nil)
		var decimals uint8
		_ = parsedABI.UnpackIntoInterface(&decimals, "decimals", decResult)

		// 调用 balanceOf(address)
		data, err := parsedABI.Pack("balanceOf", user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "balanceOf 调用打包失败"})
			return
		}
		result, err := client.CallContract(context.Background(), ethereum.CallMsg{
			To:   &contract,
			Data: data,
		}, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "调用 balanceOf 失败"})
			return
		}

		// 解码结果
		var balance *big.Int
		err = parsedABI.UnpackIntoInterface(&balance, "balanceOf", result)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "结果解码失败"})
			return
		}

		// 转换为浮点数表示
		divisor := new(big.Float).SetFloat64(float64(new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil).Int64()))
		balanceFloat := new(big.Float).Quo(new(big.Float).SetInt(balance), divisor)

		c.JSON(http.StatusOK, gin.H{
			"token":   tokenSymbol,
			"address": userAddr,
			"balance": fmt.Sprintf("%f", balanceFloat),
		})
	})
	// 启动服务
	log.Println("🚀 服务已启动: http://localhost:8080")
	r.Run(":8080")
}
