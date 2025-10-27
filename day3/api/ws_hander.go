package api

import (
	"log"
	"net/http"
	"time"
	"web3_wallect_tracker/infra"

	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true }, // 允许跨域
}

func WsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket Upgrade失败:", err)
		return
	}
	defer conn.Close()

	for {
		// 阻塞直到有数据
		result, err := infra.Rdb.BLPop(infra.Ctx, time.Second, "transfer_events").Result() // 1秒超时保持活跃
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
