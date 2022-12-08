package ws

import (
	"encoding/json"
	"jojo-live/util"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type WsMessage struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

func (m WsMessage) ToJson() []byte {
	d, _ := json.Marshal(m)
	return d
}

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

	util.WSConnMap[ws.RemoteAddr().String()] = ws

	if err := ws.WriteMessage(websocket.TextMessage, WsMessage{
		Type: "status",
		Data: util.GetStatus(),
	}.ToJson()); err != nil {
		log.Println(err)
		// return
	}

	defer func() {
		ws.Close()
		delete(util.WSConnMap, ws.RemoteAddr().String())
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
	for _, conn := range util.WSConnMap {
		if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			log.Println(err)
			// return
		}
	}
}
