package infra

import (
	"log"

	"github.com/ethereum/go-ethereum/ethclient"
)

var EthClient *ethclient.Client
var EthWsClient *ethclient.Client

func InitEthClient() {
	client, err := ethclient.Dial(TheAppConfig.RPCUrl)
	if err != nil {
		log.Fatalf("Eth client connection is failed: %v", err)
	}
	EthClient = client
	log.Println("Eth client connection success")
}
