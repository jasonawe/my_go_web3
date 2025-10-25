package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"

    "github.com/ethereum/go-ethereum"
    "github.com/ethereum/go-ethereum/accounts/abi"
    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/ethclient"
    "github.com/ethereum/go-ethereum/core/types"
	"github.com/joho/godotenv"
)

// ERC20 ABI
const erc20ABI = `[{"constant":true,"inputs":[{"name":"_owner","type":"address"}],"name":"balanceOf","outputs":[{"name":"balance","type":"uint256"}],"type":"function"},{"anonymous":false,"inputs":[{"indexed":true,"name":"from","type":"address"},{"indexed":true,"name":"to","type":"address"},{"indexed":false,"name":"value","type":"uint256"}],"name":"Transfer","type":"event"}]`

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ .env 文件未找到，将使用环境变量")
	}

	rpc := os.Getenv("RPC_URL")
	contractAddr := os.Getenv("ERC20_CONTRACT")
	walletAddr := os.Getenv("WALLET_ADDRESS")

	client, err := ethclient.Dial(rpc)
	if err != nil {
		log.Fatal("❌ 连接 RPC 失败:", err)
	}

	// 查询余额
	erc20Addr := common.HexToAddress(contractAddr)
	wallet := common.HexToAddress(walletAddr)

	parsedABI, err := abi.JSON(strings.NewReader(erc20ABI))
	if err != nil {
		log.Fatal("❌ 解析 ABI 失败:", err)
	}

	// balanceOf 调用
	data, err := parsedABI.Pack("balanceOf", wallet)
	if err != nil {
		log.Fatal(err)
	}

	msg := ethereum.CallMsg{
		To:   &erc20Addr,
		Data: data,
	}

	res, err := client.CallContract(context.Background(), msg, nil)
	if err != nil {
		log.Fatal("❌ 查询余额失败:", err)
	}

	balance := new(big.Int).SetBytes(res)
	fmt.Printf("✅ 钱包 %s ERC20 余额: %s\n", walletAddr, balance.String())
}
