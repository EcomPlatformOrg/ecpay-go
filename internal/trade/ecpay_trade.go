package trade

import (
	"fmt"
	"github.com/EcomPlatformOrg/ECpay-go/internal/client"
	"github.com/EcomPlatformOrg/ECpay-go/internal/helpers"
	"io"
	"log/slog"
	"net/http"
)

// ECPayTrade is a struct containing information for an ECPay trade
type ECPayTrade struct {
	// MerchantID 特店編號
	MerchantID string `json:"MerchantID,omitempty" form:"MerchantID"`

	// MerchantTradeNo 特店訂單編號 (唯一值, 不可重複使用)
	MerchantTradeNo string `json:"MerchantTradeNo,omitempty" form:"MerchantTradeNo"`

	// MerchantTradeDate 特店交易時間, 格式: yyyy/MM/dd HH:mm:ss
	MerchantTradeDate string `json:"MerchantTradeDate" form:"MerchantTradeDate"`

	// PaymentType 交易類型, 固定為 'aio'
	PaymentType string `json:"PaymentType,omitempty" form:"PaymentType"`

	// TotalAmount 交易金額 (新台幣, 整數, 無小數點)
	TotalAmount int `json:"TotalAmount,omitempty" form:"TotalAmount"`

	// TradeDesc 交易描述 (不可有特殊字元)
	TradeDesc string `json:"TradeDesc,omitempty" form:"TradeDesc"`

	// ItemName 商品名稱 (多筆商品以 # 分隔, 中英數 400 字內)
	ItemName string `json:"ItemName,omitempty" form:"ItemName"`

	// ReturnURL 付款完成通知回傳網址
	ReturnURL string `json:"ReturnURL,omitempty" form:"ReturnURL"`

	// ChoosePayment 選擇預設付款方式
	ChoosePayment string `json:"ChoosePayment,omitempty" form:"ChoosePayment"`

	// CheckMacValue 檢查碼
	CheckMacValue string `json:"CheckMacValue,omitempty" form:"CheckMacValue"`

	// EncryptType CheckMacValue加密類型, 使用SHA256加密
	EncryptType int `json:"EncryptType,omitempty" form:"EncryptType"`

	// StoreID 特店旗下店舖代號
	StoreID string `json:"StoreID,omitempty" form:"StoreID"`

	// ClientBackURL Client端返回特店的按鈕連結
	ClientBackURL string `json:"ClientBackURL,omitempty" form:"ClientBackURL"`

	// ItemURL 商品銷售網址
	ItemURL string `json:"ItemURL,omitempty" form:"ItemURL"`

	// Remark 備註欄位
	Remark string `json:"Remark,omitempty" form:"Remark"`

	// ChooseSubPayment 付款子項目
	ChooseSubPayment string `json:"ChooseSubPayment,omitempty" form:"ChooseSubPayment"`

	// OrderResultURL Client端回傳付款結果網址
	OrderResultURL string `json:"OrderResultURL,omitempty" form:"OrderResultURL"`

	// NeedExtraPaidInfo 是否需要額外的付款資訊 (Y: 需要, N: 不需要)
	NeedExtraPaidInfo string `json:"NeedExtraPaidInfo,omitempty" form:"NeedExtraPaidInfo"`

	// IgnorePayment 隱藏付款方式 (當ChoosePayment為ALL時使用)
	IgnorePayment string `json:"IgnorePayment,omitempty" form:"IgnorePayment"`

	// PlatformID 特約合作平台商代號
	PlatformID string `json:"PlatformID,omitempty" form:"PlatformID"`

	// CustomField1 自訂名稱欄位1
	CustomField1 string `json:"CustomField1,omitempty" form:"CustomField1"`

	// CustomField2 自訂名稱欄位2
	CustomField2 string `json:"CustomField2,omitempty" form:"CustomField2"`

	// CustomField3 自訂名稱欄位3
	CustomField3 string `json:"CustomField3,omitempty" form:"CustomField3"`

	// CustomField4 自訂名稱欄位4
	CustomField4 string `json:"CustomField4,omitempty" form:"CustomField4"`

	// Language 語系設定 (ENG: 英語, KOR: 韓語, JPN: 日語, CHI: 簡體中文)
	Language string `json:"Language,omitempty" form:"Language"`
}

// CreateAioPayment sends an HTTP POST request to create a payment transaction with AioPayment method.
// It takes an ECPayClient as a parameter and returns the response body as a string and an error, if any.
// If an error occurs during the request, it will be returned.
func (e *ECPayTrade) CreateAioPayment(client client.ECPayClient) (string, error) {

	formData := helpers.ReflectFormValues(e)

	checkMacValue := helpers.GenerateCheckMacValue(formData, client.HashKey, client.HashIV)

	formData.Set("CheckMacValue", checkMacValue)

	// 發送 HTTP POST 請求
	resp, err := http.PostForm(client.BaseURL, formData)
	if err != nil {
		slog.Error(fmt.Sprintf("Error sending POST request: %v", err))
		return "", err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	defer func(Body io.ReadCloser) {
		if err = Body.Close(); err != nil {
			slog.Error(err.Error())
		}
	}(resp.Body)

	return string(body), nil
}
