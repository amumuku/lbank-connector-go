package sve

import (
	"go.uber.org/zap"

	"github.com/LBank-exchange/lbank-connector-go/pkg" // 使用完整路径
)

type Client struct {
	ApiKey    string
	SecretKey string
	Host      string
	Debug     bool
	Logger    *zap.SugaredLogger
}

var host = "https://www.lbkex.net"

func NewClient(apiKey, secretKey string) *Client {
	devLogger := pkg.InitLogger()
	c := &Client{
		ApiKey:    apiKey,
		SecretKey: secretKey,
		Host:      host,
		Logger:    devLogger,
	}
	return c
}

func (c *Client) debug(msg string, args ...interface{}) {
	if c.Debug {
		c.Logger.Debugf(msg, args...)
	}
}

func (c *Client) SetHost(host string) {
	c.Host = host
}

func (c *Client) NewSpotService() *SpotService {
	hs := c.NewHttpService()
	return &SpotService{c: c, hs: hs}
}
func (c *Client) NewMarketService() *MarketService {
	hs := c.NewHttpService()
	return &MarketService{c: c, hs: hs}
}
func (c *Client) NewWalletService() *WalletService {
	hs := c.NewHttpService()
	return &WalletService{c: c, hs: hs}
}
func (c *Client) NewBaseService() *BaseService {
	hs := c.NewHttpService()
	return &BaseService{c: c, hs: hs}
}
func (c *Client) NewAccountService() *AccountService {
	hs := c.NewHttpService()
	return &AccountService{c: c, hs: hs}
}

func (c *Client) NewOrderService() *OrderService {
	hs := c.NewHttpService()
	return &OrderService{c: c, hs: hs}
}

func (c *Client) NewWithDrawService() *WithDrawService {
	hs := c.NewHttpService()
	return &WithDrawService{c: c, hs: hs}
}

func (c *Client) NewHttpService() *HttpService {
	return &HttpService{c: c}
}
