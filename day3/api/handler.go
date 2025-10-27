package api

import (
	"log"
	"net/http"
	"web3_wallect_tracker/service"

	"github.com/gin-gonic/gin"
)

func getEthBalanceHandler(ctx *gin.Context) {
	address := ctx.Param("address")
	ret, err := service.GetEthBalance(address)
	if err != nil {
		log.Fatalf("get eth balance failed, error:%v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "get eth balance failed"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"eth_balance": ret})
}
func getWalletBalanceHandler(ctx *gin.Context) {
	address := ctx.Param("address")
	ret, err := service.GetWalletBalance(address)
	if err != nil {
		log.Fatalf("get wallet balance failed, error:%v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "get wallet balance failed"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"wallet_balance": ret})
}
