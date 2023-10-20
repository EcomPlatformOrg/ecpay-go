package main

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

// ECPayTrade ECPay交易資料
type ECPayTrade struct {
	// MerchantID 特店編號
	MerchantID string `json:"MerchantID"`

	// MerchantTradeNo 特店訂單編號 (唯一值, 不可重複使用)
	MerchantTradeNo string `json:"MerchantTradeNo"`

	// MerchantTradeDate 特店交易時間, 格式: yyyy/MM/dd HH:mm:ss
	MerchantTradeDate string `json:"MerchantTradeDate"`

	// PaymentType 交易類型, 固定為 'aio'
	PaymentType string `json:"PaymentType"`

	// TotalAmount 交易金額 (新台幣, 整數, 無小數點)
	TotalAmount int `json:"TotalAmount"`

	// TradeDesc 交易描述 (不可有特殊字元)
	TradeDesc string `json:"TradeDesc"`

	// ItemName 商品名稱 (多筆商品以 # 分隔, 中英數 400 字內)
	ItemName string `json:"ItemName"`

	// ReturnURL 付款完成通知回傳網址
	ReturnURL string `json:"ReturnURL"`

	// ChoosePayment 選擇預設付款方式
	ChoosePayment string `json:"ChoosePayment"`

	// CheckMacValue 檢查碼
	CheckMacValue string `json:"CheckMacValue"`

	// EncryptType CheckMacValue加密類型, 使用SHA256加密
	EncryptType int `json:"EncryptType"`

	// StoreID 特店旗下店舖代號
	StoreID string `json:"StoreID,omitempty"`

	// ClientBackURL Client端返回特店的按鈕連結
	ClientBackURL string `json:"ClientBackURL,omitempty"`

	// ItemURL 商品銷售網址
	ItemURL string `json:"ItemURL,omitempty"`

	// Remark 備註欄位
	Remark string `json:"Remark,omitempty"`

	// ChooseSubPayment 付款子項目
	ChooseSubPayment string `json:"ChooseSubPayment,omitempty"`

	// OrderResultURL Client端回傳付款結果網址
	OrderResultURL string `json:"OrderResultURL,omitempty"`

	// NeedExtraPaidInfo 是否需要額外的付款資訊 (Y: 需要, N: 不需要)
	NeedExtraPaidInfo string `json:"NeedExtraPaidInfo,omitempty"`

	// IgnorePayment 隱藏付款方式 (當ChoosePayment為ALL時使用)
	IgnorePayment string `json:"IgnorePayment,omitempty"`

	// PlatformID 特約合作平台商代號
	PlatformID string `json:"PlatformID,omitempty"`

	// CustomField1 自訂名稱欄位1
	CustomField1 string `json:"CustomField1,omitempty"`

	// CustomField2 自訂名稱欄位2
	CustomField2 string `json:"CustomField2,omitempty"`

	// CustomField3 自訂名稱欄位3
	CustomField3 string `json:"CustomField3,omitempty"`

	// CustomField4 自訂名稱欄位4
	CustomField4 string `json:"CustomField4,omitempty"`

	// Language 語系設定 (ENG: 英語, KOR: 韓語, JPN: 日語, CHI: 簡體中文)
	Language string `json:"Language,omitempty"`
}

// CreateAioPayment 建立ECPay交易 (AIO)
func (e *ECPayTrade) CreateAioPayment(client ECPayClient) error {

	formData := tradeToFormValues(e)

	checkMacValue := generateCheckMacValue(formData, client.HashKey, client.HashIV)

	formData.Set("CheckMacValue", checkMacValue)

	// 發送 HTTP POST 請求
	resp, err := http.PostForm(client.BaseURL, formData)
	if err != nil {
		slog.Error(fmt.Sprintf("Error sending POST request: %v", err))
		return err
	}

	defer func(Body io.ReadCloser) {
		if err = Body.Close(); err != nil {
			slog.Error(err.Error())
		}
	}(resp.Body)

	return nil
}
