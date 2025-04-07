package main

import (
	"encoding/json"
	"fmt"

	"github.com/LBank-exchange/lbank-connector-go/sve"
)

const (
	apiKey    = "e3c55479-e314-4909-92ae-2fae5f02c4cb"
	secretKey = "MIICdgIBADANBgkqhkiG9w0BAQEFAASCAmAwggJcAgEAAoGBAI69cl0BOyIpN5mZmuy33d0mquJl0xXtQMBxlSHduzAuyy3m5PMYiCD8HKaGvqAs2gGWODqFli9EduLH9zuNBHKX0A+KHivCzneasNQE/0oiM5yv1Ad63Ob6xEOzbPqr10RoYPxOlxRWDXWUFWqcOH7qdPyHyW1ArZGmiM1vYEdhAgMBAAECgYAX7FiGjfZDO3U+ISh+FDLzJc/uMfK28hSwLFk6W9dLtAwJnXEx7SKjpJ2Iq3y3i8zeBzdVV55cPbVPPQSKzo+4BNObVsEuCNgRMNdzdGu6Mk+xUHFP412Y85Y0B4x5Y0s4VH5AG1wef9Y4F0qEwyxaKQ6zvs4ZPHSnOOukxNooAQJBAMvuhxWDC7NBCtkOZoFeqkJhx/BYUoj6vnPWzJbjVs1Ebc4FVj066b1DCajyCUXGIk80LqDNRbCGk1z1YWGVyjECQQCzL0ezhZJPLVqmNf+aZM4pDTmdrIy/mYcH+QjInz8wsSVK9kwo8OmedfXD7Bh9BDrrOKzFkSEqBnb6o7KFhNQxAkEAoz/UFW1tPVbxBycW+bM9WpyKAKXDlHIdaf/mkVd2EiYYPJdbHPL/UAnNPthageeFaaAdP45znkdsyjqIdSUC0QJAI5Th8h42HY7uD09twGUAI1rC9DKNiIaeL9EeE2i8DZk/xJEAMqkUWyklcpBxlHHAmXEZrenR4hyCh+b1zlnAIQJAem26eLmGqRFDOLhlWCBfzHtcXNCJ1z95XWcpCXhkt1PX0tn6gn91tZHCJ8Tg8SO3ye5zNv8sA1JoonoVRKwJYA=="
)

var client = sve.NewClient(apiKey, secretKey)

func TestBatchCreateOrder() {
	client.Debug = true
	client.SetHost(sve.LbankApiHost)

	// 定义多个订单
	orders := []map[string]string{
		{
			"symbol":    "lbk_usdt",
			"type":      "buy",
			"price":     "0.01",
			"amount":    "10",
			"custom_id": "test1",
		},
		{
			"symbol":    "lbk_usdt",
			"type":      "sell",
			"price":     "0.02",
			"amount":    "5",
			"custom_id": "test2",
		},
	}

	// 将订单数组转换为 JSON 字符串
	ordersJSON, err := json.Marshal(orders)
	if err != nil {
		fmt.Printf("Failed to marshal orders: %v\n", err)
		return
	}

	// 构造批量下单参数
	data := map[string]string{
		"orders": string(ordersJSON),
	}

	// 调用批量下单
	client.NewOrderService().BatchCreateOrder(data)
}

func main() {
	TestBatchCreateOrder()
}
