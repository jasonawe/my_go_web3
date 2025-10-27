package task

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"strings"
	"web3_wallect_tracker/info"
	"web3_wallect_tracker/infra"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func WsTask() {
	contracts := strings.Split(infra.TheAppConfig.ERC20Contracts, ",")

	// 启动 ERC20 监听
	for _, c := range contracts {
		go ListenTransfer(common.HexToAddress(c))
	}
}
func ListenTransfer(contractAddr common.Address) {
	client, err := ethclient.Dial(infra.TheAppConfig.WsUrl)
	if err != nil {
		log.Fatalf("❌ RPC 连接失败: %v", err)
	}
	parsedABI, _ := abi.JSON(strings.NewReader(info.WS_ERC20ABI))
	decimals := getDecimals(client, contractAddr, parsedABI)
	symbol := getSymbol(client, contractAddr, parsedABI)

	transferTopic := common.HexToHash(info.TRANSFER_TOPIC)
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
			infra.Rdb.LPush(infra.Ctx, "transfer_events", msg)
		}
	}
}

func getDecimals(client *ethclient.Client, contract common.Address, parsedABI abi.ABI) uint8 {
	data, _ := parsedABI.Pack("decimals")
	res, _ := client.CallContract(infra.Ctx, ethereum.CallMsg{To: &contract, Data: data}, nil)
	var decimals uint8
	_ = parsedABI.UnpackIntoInterface(&decimals, "decimals", res)
	return decimals
}
func getSymbol(client *ethclient.Client, contract common.Address, parsedABI abi.ABI) string {
	data, _ := parsedABI.Pack("symbol")
	res, _ := client.CallContract(infra.Ctx, ethereum.CallMsg{To: &contract, Data: data}, nil)
	var symbol string
	_ = parsedABI.UnpackIntoInterface(&symbol, "symbol", res)
	return symbol
}
