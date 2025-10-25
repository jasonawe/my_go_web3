package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
)

func main() {
	// 1. 加载环境变量
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// 2. 读取 RPC_URL
	rpcURL := os.Getenv("RPC_URL")
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		log.Fatalf("连接以太坊节点失败: %v", err)
	}
	defer client.Close()

	// 3. 目标钱包地址
	address := common.HexToAddress(os.Getenv("WALLET_ADDRESS"))

	// 4. 获取 ETH 余额（单位：Wei）
	balanceWei, err := client.BalanceAt(context.Background(), address, nil)
	if err != nil {
		log.Fatalf("查询余额失败: %v", err)
	}

	// 5. 转换为 ETH 单位（1 ETH = 1e18 Wei）
	balanceEth := new(big.Float).Quo(new(big.Float).SetInt(balanceWei), big.NewFloat(1e18))

	fmt.Printf("钱包地址: %s\n", address.Hex())
	fmt.Printf("余额: %s ETH\n", balanceEth.Text('f', 8))
}
