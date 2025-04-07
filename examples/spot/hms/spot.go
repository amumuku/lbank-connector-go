package main

import "github.com/LBank-exchange/lbank-connector-go/sve"

const (
	apiKey    = "7c991120-e3ec-4b9a-8fcc-aa9d533d1969"
	secretKey = "F129FDE895B6AE7850CE4E72D2A796CF"
)

var client = sve.NewClient(apiKey, secretKey)

func TestCreateOrder() {
	client.Debug = true
	client.SetHost(sve.LbankApiHost)
	data := map[string]string{
		"symbol":    "lbk_usdt",
		"type":      "buy",
		"price":     "0.01",
		"amount":    "10",
		"custom_id": "test",
	}
	client.NewOrderService().CreateOrder(data)
}

func main() {
	TestCreateOrder()
}
