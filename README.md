# LBank Connector Go Connector


## Supported API Endpoints:
- Account: `examples/account/rsa/account.go`
- Base: `examples/account/rsa/base.go`
- Market: `examples/account/rsa/market.go`
- Order: `examples/account/rsa/order.go`
- Spot: `examples/account/rsa/spot.go`
- Wallet: `examples/account/rsa/wallet.go`
- WithDraw: `examples/account/rsa/withdraw.go`


## Installation
```shell
go get github.com/LBank-exchange/lbank-connector-go
```

## Import
```golang
import (
    "github.com/LBank-exchange/lbank-connector-go"
)
```
## Authentication
```go
client := sve.NewClient("yourApiKey", "yourSecretKey")
// Debug Mode
client.Debug = true
```

## REST API
Create an  example
```go
package main

import (
	"github.com/LBank-exchange/lbank-connector-go"
)

// rsa
const (
	apiKey    = ""
	secretKey = ""
)

var client = sve.NewClient(apiKey, secretKey)

func TestUserInfo() {
	client.Debug = true
	client.SetHost(sve.LbankApiHost)
	data := map[string]string{}
	client.NewAccountService().UserInfo(data)
}
func main() {
	TestUserInfo()
}
```
