package trade

import (
	"github.com/EcomPlatformOrg/ecpay-go/pkg/client"
	"github.com/EcomPlatformOrg/ecpay-go/pkg/helpers"
	"github.com/EcomPlatformOrg/ecpay-go/pkg/model"
)

// ECPayTrade is a struct containing information for an ECPay trade
type ECPayTrade struct {

	// BaseModel 通用參數
	model.BaseModel `json:",inline"`

	// Merchant 特店資訊
	model.Merchant `json:",inline"`

	// PaymentType 交易類型, 固定為 'aio'
	PaymentType string `json:"PaymentType,omitempty" form:"PaymentType"`

	// TotalAmount 交易金額 (新台幣, 整數, 無小數點)
	TotalAmount int `json:"TotalAmount,omitempty" form:"TotalAmount"`

	// ItemName 商品名稱 (多筆商品以 # 分隔, 中英數 400 字內)
	ItemName string `json:"ItemName,omitempty" form:"ItemName"`

	// ReturnURL 付款完成通知回傳網址
	ReturnURL string `json:"ReturnURL,omitempty" form:"ReturnURL"`

	// ChoosePayment 選擇預設付款方式
	ChoosePayment string `json:"ChoosePayment,omitempty" form:"ChoosePayment"`

	// EncryptType CheckMacValue加密類型, 使用SHA256加密
	EncryptType int `json:"EncryptType,omitempty" form:"EncryptType"`

	// StoreID 特店旗下店舖代號
	StoreID string `json:"StoreID,omitempty" form:"StoreID"`

	// ClientBackURL Client端返回特店的按鈕連結
	ClientBackURL string `json:"ClientBackURL,omitempty" form:"ClientBackURL"`

	// ItemURL 商品銷售網址
	ItemURL string `json:"ItemURL,omitempty" form:"ItemURL"`

	// ChooseSubPayment 付款子項目
	ChooseSubPayment string `json:"ChooseSubPayment,omitempty" form:"ChooseSubPayment"`

	// OrderResultURL Client端回傳付款結果網址
	OrderResultURL string `json:"OrderResultURL,omitempty" form:"OrderResultURL"`

	// NeedExtraPaidInfo 是否需要額外的付款資訊 (Y: 需要, N: 不需要)
	NeedExtraPaidInfo string `json:"NeedExtraPaidInfo,omitempty" form:"NeedExtraPaidInfo"`

	// IgnorePayment 隱藏付款方式 (當ChoosePayment為ALL時使用)
	IgnorePayment string `json:"IgnorePayment,omitempty" form:"IgnorePayment"`

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
func (e *ECPayTrade) CreateAioPayment(c client.ECPayClient) (string, error) {

	formData := helpers.ReflectFormValues(e)

	checkMacValue := helpers.GenerateCheckMacValue(formData, c.HashKey, c.HashIV)

	formData.Set("CheckMacValue", checkMacValue)

	body, err := helpers.SendFormData(c, formData)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
