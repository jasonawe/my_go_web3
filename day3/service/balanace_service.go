package service

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"
	"web3_wallect_tracker/info"
	"web3_wallect_tracker/infra"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func GetEthBalance(address string) (string, error) {
	account := common.HexToAddress(address)
	balance, err := infra.EthClient.BalanceAt(context.Background(), account, nil)
	if err != nil {
		return "", err
	}
	ethValue := new(big.Float).Quo(new(big.Float).SetInt(balance), big.NewFloat(1e18))
	return ethValue.Text('f', 6), nil
}

func GetWalletBalance(address string) (*info.WalletResponse, error) {
	rpc_url := os.Getenv("RPC_URL")
	client, err := ethclient.Dial(rpc_url)
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to ETH:%v", err)
	}
	defer client.Close()
	addr := common.HexToAddress(address)
	ctx := context.Background()
	// 获取eth的余额
	balance, err := client.BalanceAt(ctx, addr, nil)
	if err != nil {
		return nil, fmt.Errorf("get ETH balance error:%v", err)
	}
	// 格式化
	ethVal := new(big.Float).Quo(new(big.Float).SetInt(balance), big.NewFloat(1e18))
	var tokenInfos []info.TokenInfo

	erc20ABI, _ := abi.JSON(strings.NewReader(info.ERC20ABI))

	for symbol, contraceAddr := range info.TokenList {
		contract := common.HexToAddress(contraceAddr)
		data, err := erc20ABI.Pack("balanceOf", addr)
		if err != nil {
			log.Printf("get contract Pack %s by address %s error:%v", contract, address, err)
			continue
		}
		callMsg := ethereum.CallMsg{To: &contract, Data: data}
		ret, err := client.CallContract(ctx, callMsg, nil)
		if err != nil {
			log.Printf("get contract CallContract %s by address %s error:%v", contract, address, err)
			continue
		}
		var balance = new(big.Int)

		err = erc20ABI.UnpackIntoInterface(&balance, "balanceOf", ret)
		if err != nil {
			log.Printf("get contract UnpackIntoInterface %s by address %s error:%v", contract, address, err)
			continue
		}
		tokenVal := new(big.Float).Quo(new(big.Float).SetInt(balance), big.NewFloat(1e6)) // 假设 6 位精度
		tokenInfos = append(tokenInfos, info.TokenInfo{
			Address: contraceAddr,
			Symbol:  symbol,
			Balance: tokenVal.Text('f', 2),
		})
	}
	return &info.WalletResponse{
		Address:    address,
		EthBalance: ethVal.Text('f', 4),
		Tokens:     tokenInfos,
	}, nil
}
