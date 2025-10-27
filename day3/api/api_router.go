package api

import "github.com/gin-gonic/gin"

func InitEngine() *gin.Engine {
	r := gin.Default()
	// api路由
	api := r.Group("/api")
	{
		api.GET("/balance/eth/:address", getEthBalanceHandler)
		api.GET("/balance/wallet/:address", getWalletBalanceHandler)
	}
	// WebSocket 路由
	r.GET("/ws", func(c *gin.Context) {
		WsHandler(c.Writer, c.Request)
	})

	// 静态文件挂载在 /static
	r.Static("/static", "./static")
	return r
}
