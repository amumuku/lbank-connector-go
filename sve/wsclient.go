package sve

import (
	"go.uber.org/zap"

	"github.com/LBank-exchange/lbank-connector-go/pkg" // 使用完整路径
)

type WsClient struct {
	ApiKey    string
	SecretKey string
	Host      string
	Debug     bool
	Logger    *zap.SugaredLogger
}

//var host = "https://www.lbkex.net"

func NewWsClient(apiKey, secretKey string) *WsClient {
	devLogger := pkg.InitLogger()
	wsClient := &WsClient{
		ApiKey:    apiKey,
		SecretKey: secretKey,
		Host:      host,
		Logger:    devLogger,
	}
	return wsClient
}

func (wc *WsClient) debug(msg string, args ...interface{}) {
	if wc.Debug {
		wc.Logger.Debugf(msg, args...)
	}
}

func (wc *WsClient) SetHost(host string) {
	wc.Host = host
}

func (wc *WsClient) NewWsService() *WsService {
	return &WsService{Wc: wc}
}

func (wc *WsClient) NewWsMarketService() *WsMarketService {
	ws := wc.NewWsService()
	return &WsMarketService{ws}
}
