package sve

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/gorilla/websocket"

	"github.com/LBank-exchange/lbank-connector-go/pkg" // 使用完整路径
)

type WsService struct {
	conn *websocket.Conn
	Wc   *WsClient
}

func (w *WsService) KeepAlive(timeout time.Duration) {
	ticker := time.NewTicker(timeout)
	go func() {
		defer ticker.Stop()
		for {
			pingMap := map[string]string{
				"ping":   pkg.RandomUUID(),
				"action": "ping",
			}
			msg, _ := json.Marshal(pingMap)
			w.SendMsg(msg)
			<-ticker.C
		}
	}()
}

func (w *WsService) CreateWsConn() (*websocket.Conn, error) {
	u := url.URL{Scheme: "wss", Host: "www.lbkex.net", Path: "/ws/V2/"}
	if w.Wc.Debug {
		w.Wc.debug("CreateWsConn: %s", u.String())
	}
	dialer := websocket.DefaultDialer
	conn, _, err := dialer.Dial(u.String(), nil)
	if err != nil {
		w.Wc.Logger.Fatal("dial:", err)
	}
	w.conn = conn
	return conn, err
}

func (w *WsService) SendMsg(content []byte) {
	err := w.conn.WriteMessage(websocket.TextMessage, content)
	if err != nil {
		w.Wc.Logger.Fatal("SendMsg:", err)
		return
	}
	if w.Wc.Debug {
		w.Wc.debug("Sent: %s", content)

	}
	log.Printf("Sent: %s", content)
}

func (w *WsService) ReceiveMsg() {
	for {
		_, message, err := w.conn.ReadMessage()
		if err != nil {
			w.Wc.Logger.Fatal("read:", err)
			return
		}
		w.Wc.debug("Received: %s", message)
		fmt.Printf("Received: %s", message)
	}
}
