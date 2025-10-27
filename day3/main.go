package main

import (
	"fmt"
	"web3_wallect_tracker/api"
	"web3_wallect_tracker/infra"
	"web3_wallect_tracker/task"
)

func main() {
	infra.LoadConfig()
	infra.InitRedisClient()
	infra.InitEthClient()
	r := api.InitEngine()
	serverAddr := fmt.Sprintf("%s:%d", infra.TheAppConfig.SevConfig.Ip, infra.TheAppConfig.SevConfig.Port)

	task.WsTask()
	r.Run(serverAddr)
}
