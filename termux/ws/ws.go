package ws

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var (
	WSConnMap = make(map[string]*websocket.Conn)
)

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// webSocket请求ping 返回pong
func Ws(c *gin.Context) {
	//升级get请求为webSocket协议
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	WSConnMap[ws.RemoteAddr().String()] = ws

	defer func() {
		ws.Close()
		delete(WSConnMap, ws.RemoteAddr().String())
	}()

	for {
		//读取ws中的数据
		mt, message, err := ws.ReadMessage()
		if err != nil {
			break
		}
		if string(message) == "ping" {
			message = []byte("pong")
		}
		//写入ws数据
		err = ws.WriteMessage(mt, message)
		if err != nil {
			break
		}
	}
}

func WsBroadcast(msg []byte) {
	for _, conn := range WSConnMap {
		if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			log.Println(err)
			// return
		}
	}
}
