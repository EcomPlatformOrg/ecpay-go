package main

import (
	"github.com/linshan1234/ECPay-go/models"
	"github.com/linshan1234/ECPay-go/sdk"
)

func main() {
	//訂單範例
	orderData := models.Order{
		ID:          1676530197384,
		TotalAmount: 999,
		TradeDesc:   "Desc",
	}

	sdk.CreateOrder(orderData)
}
