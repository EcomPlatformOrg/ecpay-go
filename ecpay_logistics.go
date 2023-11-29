package ECpay_go

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

// ECPayLogistics is a struct containing information for an ECPay logistics
type ECPayLogistics struct {

	// MerchantID 特店編號
	MerchantID string `json:"MerchantID"`

	// MerchantTradeNo 特店交易編號
	MerchantTradeNo string `json:"MerchantTradeNo,omitempty"`

	// MerchantTradeDate 廠商交易時間
	MerchantTradeDate string `json:"MerchantTradeDate,omitempty"`

	// LogisticsType 物流類型
	LogisticsType string `json:"LogisticsType,omitempty"`

	// LogisticsSubType 物流子類型
	LogisticsSubType string `json:"LogisticsSubType"`

	// GoodsAmount 商品金額
	GoodsAmount int `json:"GoodsAmount,omitempty"`

	// CollectionAmount 代收金額
	CollectionAmount int `json:"CollectionAmount,omitempty"`

	// IsCollection 是否代收貨款
	IsCollection string `json:"IsCollection,omitempty"`

	// GoodsName 商品名稱
	GoodsName string `json:"GoodsName,omitempty"`

	// SenderName 寄件人姓名
	SenderName string `json:"SenderName,omitempty"`

	// SenderPhone 寄件人電話
	SenderPhone string `json:"SenderPhone,omitempty"`

	// SenderCellPhone 寄件人手機
	SenderCellPhone string `json:"SenderCellPhone,omitempty"`

	// ReceiverName 收件人姓名
	ReceiverName string `json:"ReceiverName,omitempty"`

	// ReceiverPhone 收件人電話
	ReceiverPhone string `json:"ReceiverPhone,omitempty"`

	// ReceiverCellPhone 收件人手機
	ReceiverCellPhone string `json:"ReceiverCellPhone,omitempty"`

	// ReceiverEmail 收件人email
	ReceiverEmail string `json:"ReceiverEmail,omitempty"`

	// ReceiverStoreID 收件人門市代號
	ReceiverStoreID string `json:"ReceiverStoreID,omitempty"`

	// ReturnStoreID 退貨門市代號
	ReturnStoreID string `json:"ReturnStoreID,omitempty"`

	// TradeDesc 交易描述
	TradeDesc string `json:"TradeDesc,omitempty"`

	// ServerReplyURL Server端回覆網址
	ServerReplyURL string `json:"ServerReplyURL,omitempty"`

	// ClientReplyURL Client端回覆網址
	ClientReplyURL string `json:"ClientReplyURL,omitempty"`

	// Remark 備註
	Remark string `json:"Remark,omitempty"`

	// PlatformID 特約合作平台商代號
	PlatformID string `json:"PlatformID,omitempty"`

	// CheckMacValue 檢查碼
	CheckMacValue string `json:"CheckMacValue"`

	// RtnCode 目前物流狀態
	RtnCode string `json:"RtnCode,omitempty"`

	// RtnMsg 物流狀態說明
	RtnMsg string `json:"RtnMsg,omitempty"`

	// AllPayLogisticsID 綠界科技的物流交易編號
	AllPayLogisticsID string `json:"AllPayLogisticsID,omitempty"`

	// UpdateStatusDate 物流狀態更新時間
	UpdateStatusDate string `json:"UpdateStatusDate,omitempty"`

	// ReceiverAddress 收件人地址
	ReceiverAddress string `json:"ReceiverAddress,omitempty"`

	// CVSPaymentNo 寄貨編號
	CVSPaymentNo string `json:"CVSPaymentNo,omitempty"`

	// CVSValidationNo 驗證碼
	CVSValidationNo string `json:"CVSValidationNo,omitempty"`

	// BookingNote 托運單號
	BookingNote string `json:"BookingNote,omitempty"`
}

// CreateExpress 綠界物流門市訂單建立
func (e *ECPayLogistics) CreateExpress(client ECPayClient) error {

	formData := logisticsToFormValues(e)

	checkMacValue := generateCheckMacValue(formData, client.HashKey, client.HashIV)
	formData.Set("CheckMacValue", checkMacValue)

	resp, err := http.PostForm(client.BaseURL, formData)
	if err != nil {
		slog.Error(fmt.Sprintf("Error sending POST request: %v", err))
		return err
	}
	defer func(Body io.ReadCloser) {
		if err = Body.Close(); err != nil {
			slog.Error(fmt.Sprintf("Error closing response body: %v", err))
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(body, &e); err != nil {
		return err
	}

	return nil

}
