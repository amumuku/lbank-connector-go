package main

import (
	"encoding/json"
	"fmt"

	"github.com/LBank-exchange/lbank-connector-go/sve"
)

const (
	apiKey    = "7c991120-e3ec-4b9a-8fcc-aa9d533d1969"
	secretKey = "F129FDE895B6AE7850CE4E72D2A796CF"
)

var client = sve.NewClient(apiKey, secretKey)

// func TestCreateOrder() {
// 	client.Debug = true
// 	client.SetHost(sve.LbankApiHost)
// 	data := map[string]string{
// 		"symbol": "lbk_usdt",
// 		"size":   "1",
// 	}
// 	client.NewOrderService().CreateOrder(data)
// }

// func TestBatchCreateOrder() {
// 	client.Debug = true
// 	client.SetHost(sve.LbankApiHost)

// 	// 定义多个订单
// 	orders := []map[string]string{
// 		{
// 			"symbol":    "lbk_usdt",
// 			"type":      "buy",
// 			"price":     "0.01",
// 			"amount":    "10",
// 			"custom_id": "test5",
// 		},
// 		{
// 			"symbol":    "lbk_usdt",
// 			"type":      "sell",
// 			"price":     "0.02",
// 			"amount":    "5",
// 			"custom_id": "test6",
// 		},
// 	}

// 	// 将订单数组转换为 JSON 字符串
// 	ordersJSON, err := json.Marshal(orders)
// 	if err != nil {
// 		fmt.Printf("Failed to marshal orders: %v\n", err)
// 		return
// 	}

// 	// 构造批量下单参数
// 	data := map[string]string{
// 		"orders": string(ordersJSON),
// 	}

// 	// 调用批量下单
// 	client.NewOrderService().BatchCreateOrder(data)
// }

// func TestOrdersInfoNoDeal() {
// 	client.Debug = true
// 	client.SetHost("https://www.lbkex.net") // 如果 LbankApiHost 未定义

// 	// 构造请求参数
// 	data := map[string]string{
// 		"symbol":       "lbk_usdt",
// 		"current_page": "1",
// 		"page_length":  "10",
// 	}

// 	// 调用 OrdersInfoNoDeal
// 	spotSvc := client.NewSpotService()
// 	spotSvc.OrdersInfoNoDeal(data)
// 	hs := spotSvc.GetHttpService()
// 	// 检查请求是否成功
// 	if hs.Error != nil {
// 		fmt.Printf("查询未成交订单失败: %v\n", hs.Error)
// 		return
// 	}

// 	// 定义响应结构体
// 	type OrderInfo struct {
// 		Symbol              string  `json:"symbol"`
// 		OrderID             string  `json:"orderId"`
// 		ClientOrderID       string  `json:"clientOrderId"`
// 		Price               float64 `json:"price"`
// 		OrigQty             float64 `json:"origQty"`
// 		ExecutedQty         float64 `json:"executedQty"`
// 		CummulativeQuoteQty float64 `json:"cummulativeQuoteQty"`
// 		Status              int     `json:"status"`
// 		Type                string  `json:"type"`
// 		Time                int64   `json:"time"`
// 		UpdateTime          int64   `json:"updateTime"`
// 		OrigQuoteOrderQty   float64 `json:"origQuoteOrderQty"`
// 	}

// 	type Response struct {
// 		Msg    string `json:"msg"`
// 		Result bool   `json:"result"`
// 		Data   struct {
// 			Total       int         `json:"total"`
// 			PageLength  int         `json:"page_length"`
// 			Orders      []OrderInfo `json:"orders"`
// 			CurrentPage int         `json:"current_page"`
// 		} `json:"data"`
// 		ErrorCode int   `json:"error_code"`
// 		Ts        int64 `json:"ts"`
// 	}

// 	// 解析响应
// 	var resp Response
// 	err := json.Unmarshal([]byte(hs.Text), &resp)
// 	if err != nil {
// 		fmt.Printf("解析响应失败: %v\n", err)
// 		return
// 	}

// 	// 检查 API 返回是否成功
// 	if !resp.Result || resp.ErrorCode != 0 {
// 		fmt.Printf("API 返回失败: msg=%s, error_code=%d\n", resp.Msg, resp.ErrorCode)
// 		return
// 	}

// 	// 输出订单信息
// 	for _, order := range resp.Data.Orders {
// 		fmt.Printf("订单: Symbol=%s, OrderID=%s, ClientOrderID=%s, Price=%.2f, OrigQty=%.2f, ExecutedQty=%.2f, Status=%d, Type=%s, Time=%d, UpdateTime=%d\n",
// 			order.Symbol, order.OrderID, order.ClientOrderID, order.Price, order.OrigQty, order.ExecutedQty, order.Status, order.Type, order.Time, order.UpdateTime)
// 	}
// }

// func TestOrdersAndCancel1() {
// 	client.Debug = true
// 	client.SetHost("https://www.lbkex.net")

// 	// 1. 查询未成交订单
// 	data := map[string]string{
// 		"symbol":       "lbk_usdt",
// 		"current_page": "1",
// 		"page_length":  "10",
// 	}

// 	spotSvc := client.NewSpotService()
// 	respBytes, err := spotSvc.OrdersInfoNoDeal(data)
// 	if err != nil {
// 		fmt.Printf("查询未成交订单失败: %v\n", err)
// 		return
// 	}

// 	// 定义响应结构体
// 	type OrderInfo struct {
// 		Symbol              string  `json:"symbol"`
// 		OrderID             string  `json:"orderId"`
// 		ClientOrderID       string  `json:"clientOrderId"`
// 		Price               float64 `json:"price"`
// 		OrigQty             float64 `json:"origQty"`
// 		ExecutedQty         float64 `json:"executedQty"`
// 		CummulativeQuoteQty float64 `json:"cummulativeQuoteQty"`
// 		Status              int     `json:"status"`
// 		Type                string  `json:"type"`
// 		Time                int64   `json:"time"`
// 		UpdateTime          int64   `json:"updateTime"`
// 		OrigQuoteOrderQty   float64 `json:"origQuoteOrderQty"`
// 	}

// 	type Response struct {
// 		Msg    string `json:"msg"`
// 		Result bool   `json:"result"`
// 		Data   struct {
// 			Total       int         `json:"total"`
// 			PageLength  int         `json:"page_length"`
// 			Orders      []OrderInfo `json:"orders"`
// 			CurrentPage int         `json:"current_page"`
// 		} `json:"data"`
// 		ErrorCode int   `json:"error_code"`
// 		Ts        int64 `json:"ts"`
// 	}

// 	// 解析查询响应
// 	var resp Response
// 	if err := json.Unmarshal(respBytes, &resp); err != nil {
// 		fmt.Printf("解析查询响应失败: %v\n", err)
// 		return
// 	}

// 	if !resp.Result || resp.ErrorCode != 0 {
// 		fmt.Printf("查询 API 返回失败: msg=%s, error_code=%d\n", resp.Msg, resp.ErrorCode)
// 		return
// 	}

// 	// 输出未成交订单
// 	fmt.Println("未成交订单：")
// 	for _, order := range resp.Data.Orders {
// 		fmt.Printf("Symbol=%s, OrderID=%s, ClientOrderID=%s, Price=%.2f, OrigQty=%.2f, ExecutedQty=%.2f, Status=%d, Type=%s, Time=%d\n",
// 			order.Symbol, order.OrderID, order.ClientOrderID, order.Price, order.OrigQty, order.ExecutedQty, order.Status, order.Type, order.Time)
// 	}

// 	// 2. 批量撤单
// 	var cancelOrders []map[string]string
// 	for _, order := range resp.Data.Orders {
// 		if order.Status == 0 || order.Status == 1 { // 未成交或部分成交
// 			cancelOrders = append(cancelOrders, map[string]string{
// 				"symbol":            order.Symbol,
// 				"origClientOrderId": order.ClientOrderID,
// 			})
// 		}
// 	}

// 	if len(cancelOrders) > 0 {
// 		cancelResp, err := spotSvc.CancelClientOrders(cancelOrders)
// 		if err != nil {
// 			fmt.Printf("批量撤单失败: %v\n", err)
// 			return
// 		}

// 		// 解析撤单响应（假设返回格式为 {"result":true,"order_id":"xxx"}）
// 		type CancelResponse struct {
// 			Result    bool   `json:"result"`
// 			OrderID   string `json:"order_id"`
// 			ErrorCode int    `json:"error_code"`
// 		}

// 		for i, respText := range cancelResp {
// 			var cancelResult CancelResponse
// 			if err := json.Unmarshal([]byte(respText), &cancelResult); err != nil {
// 				fmt.Printf("解析撤单响应失败 (订单 %s): %v\n", cancelOrders[i]["origClientOrderId"], err)
// 				continue
// 			}
// 			if cancelResult.Result {
// 				fmt.Printf("成功取消订单: ClientOrderID=%s, OrderID=%s\n", cancelOrders[i]["origClientOrderId"], cancelResult.OrderID)
// 			} else {
// 				fmt.Printf("取消订单失败: ClientOrderID=%s, error_code=%d\n", cancelOrders[i]["origClientOrderId"], cancelResult.ErrorCode)
// 			}
// 		}
// 	} else {
// 		fmt.Println("没有需要撤单的订单")
// 	}
// }

func TestOrdersAndCancel() {
	client.Debug = true
	client.SetHost("https://www.lbkex.net")

	// 1. 查询未成交订单
	spotSvc := client.NewSpotService()
	queryData := map[string]string{
		"symbol":       "lbk_usdt",
		"current_page": "1",
		"page_length":  "10",
	}
	respBytes, err := spotSvc.OrdersInfoNoDeal(queryData)
	if err != nil {
		fmt.Printf("查询未成交订单失败: %v\n", err)
		return
	}

	// 定义查询响应结构体
	type OrderInfo struct {
		Symbol              string  `json:"symbol"`
		OrderID             string  `json:"orderId"`
		ClientOrderID       string  `json:"clientOrderId"`
		Price               float64 `json:"price"`
		OrigQty             float64 `json:"origQty"`
		ExecutedQty         float64 `json:"executedQty"`
		CummulativeQuoteQty float64 `json:"cummulativeQuoteQty"`
		Status              int     `json:"status"`
		Type                string  `json:"type"`
		Time                int64   `json:"time"`
		UpdateTime          int64   `json:"updateTime"`
		OrigQuoteOrderQty   float64 `json:"origQuoteOrderQty"`
	}

	type QueryResponse struct {
		Msg    string `json:"msg"`
		Result bool   `json:"result"`
		Data   struct {
			Total       int         `json:"total"`
			PageLength  int         `json:"page_length"`
			Orders      []OrderInfo `json:"orders"`
			CurrentPage int         `json:"current_page"`
		} `json:"data"`
		ErrorCode int   `json:"error_code"`
		Ts        int64 `json:"ts"`
	}

	// 解析查询响应
	var queryResp QueryResponse
	if err := json.Unmarshal(respBytes, &queryResp); err != nil {
		fmt.Printf("解析查询响应失败: %v\n", err)
		return
	}

	if !queryResp.Result || queryResp.ErrorCode != 0 {
		fmt.Printf("查询 API 返回失败: msg=%s, error_code=%d\n", queryResp.Msg, queryResp.ErrorCode)
		return
	}

	// 输出未成交订单
	fmt.Println("未成交订单：")
	for _, order := range queryResp.Data.Orders {
		fmt.Printf("Symbol=%s, OrderID=%s, ClientOrderID=%s, Price=%.2f, OrigQty=%.2f, ExecutedQty=%.2f, Status=%d, Type=%s\n",
			order.Symbol, order.OrderID, order.ClientOrderID, order.Price, order.OrigQty, order.ExecutedQty, order.Status, order.Type)
	}

	// 2. 准备批量撤单参数
	var cancelOrders []map[string]string
	for _, order := range queryResp.Data.Orders {
		if order.Status == 0 || order.Status == 1 { // 未成交或部分成交
			cancelOrders = append(cancelOrders, map[string]string{
				"symbol":            order.Symbol,
				"origClientOrderId": order.ClientOrderID,
			})
		}
	}

	if len(cancelOrders) == 0 {
		fmt.Println("没有需要撤单的订单")
		return
	}

	// 3. 调用 CancelClientOrders
	orderSvc := client.NewOrderService()
	cancelResp, err := orderSvc.CancelClientOrders(cancelOrders)
	if err != nil {
		fmt.Printf("批量撤单失败: %v\n", err)
		return
	}

	// 定义撤单响应结构体（根据实际响应调整）
	type CancelResponse struct {
		Result    bool     `json:"result"`
		Msg       string   `json:"msg"`
		ErrorCode int      `json:"error_code"`
		Ts        int64    `json:"ts"`
		Data      []string `json:"data"` // 假设返回取消的订单 ID 列表
	}

	// 解析撤单响应
	var cancelResult CancelResponse
	if err := json.Unmarshal(cancelResp, &cancelResult); err != nil {
		fmt.Printf("解析撤单响应失败: %v\n", err)
		return
	}

	if !cancelResult.Result || cancelResult.ErrorCode != 0 {
		fmt.Printf("撤单 API 返回失败: msg=%s, error_code=%d\n", cancelResult.Msg, cancelResult.ErrorCode)
		return
	}

	// 输出撤单结果
	fmt.Println("批量撤单成功：")
	for i, orderID := range cancelResult.Data {
		fmt.Printf("已取消订单: ClientOrderID=%s, OrderID=%s\n", cancelOrders[i]["origClientOrderId"], orderID)
	}
}

func main() {
	// TestCreateOrder()
	// TestOrdersInfoNoDeal()
	// TestOrdersAndCancel1()
	TestOrdersAndCancel()

}
