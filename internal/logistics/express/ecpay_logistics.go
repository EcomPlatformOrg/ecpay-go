package express

import (
	"encoding/json"
	"fmt"
	"github.com/EcomPlatformOrg/ECpay-go/internal/client"
	"github.com/EcomPlatformOrg/ECpay-go/internal/helpers"
	"io"
	"log/slog"
	"net/http"
)

// ECPayLogistics is a struct containing information for an ECPay logistics
type ECPayLogistics struct {

	// MerchantID 特店編號
	MerchantID string `json:"MerchantID,omitempty" form:"MerchantID"`

	// MerchantTradeNo 特店交易編號
	MerchantTradeNo string `json:"MerchantTradeNo,omitempty" form:"MerchantTradeNo"`

	// MerchantTradeDate 廠商交易時間
	MerchantTradeDate string `json:"MerchantTradeDate,omitempty" form:"MerchantTradeDate"`

	// LogisticsType 物流類型
	LogisticsType string `json:"LogisticsType,omitempty" form:"LogisticsType"`

	// LogisticsSubType 物流子類型
	LogisticsSubType string `json:"LogisticsSubType,omitempty" form:"LogisticsSubType"`

	// GoodsAmount 商品金額
	GoodsAmount int `json:"GoodsAmount,omitempty" form:"GoodsAmount"`

	// CollectionAmount 代收金額
	CollectionAmount int `json:"CollectionAmount,omitempty" form:"CollectionAmount"`

	// IsCollection 是否代收貨款
	IsCollection string `json:"IsCollection,omitempty" form:"IsCollection"`

	// GoodsName 商品名稱
	GoodsName string `json:"GoodsName,omitempty" form:"GoodsName"`

	// SenderName 寄件人姓名
	SenderName string `json:"SenderName,omitempty" form:"SenderName"`

	// SenderPhone 寄件人電話
	SenderPhone string `json:"SenderPhone,omitempty" form:"SenderPhone"`

	// SenderCellPhone 寄件人手機
	SenderCellPhone string `json:"SenderCellPhone,omitempty" form:"SenderCellPhone"`

	// ReceiverName 收件人姓名
	ReceiverName string `json:"ReceiverName,omitempty" form:"ReceiverName"`

	// ReceiverPhone 收件人電話
	ReceiverPhone string `json:"ReceiverPhone,omitempty" form:"ReceiverPhone"`

	// ReceiverCellPhone 收件人手機
	ReceiverCellPhone string `json:"ReceiverCellPhone,omitempty" form:"ReceiverCellPhone"`

	// ReceiverEmail 收件人email
	ReceiverEmail string `json:"ReceiverEmail,omitempty" form:"ReceiverEmail"`

	// ReceiverStoreID 收件人門市代號
	ReceiverStoreID string `json:"ReceiverStoreID,omitempty" form:"ReceiverStoreID"`

	// ReturnStoreID 退貨門市代號
	ReturnStoreID string `json:"ReturnStoreID,omitempty" form:"ReturnStoreID"`

	// TradeDesc 交易描述
	TradeDesc string `json:"TradeDesc,omitempty" form:"TradeDesc"`

	// ServerReplyURL Server端回覆網址
	ServerReplyURL string `json:"ServerReplyURL,omitempty" form:"ServerReplyURL"`

	// ClientReplyURL Client端回覆網址
	ClientReplyURL string `json:"ClientReplyURL,omitempty" form:"ClientReplyURL"`

	// Remark 備註
	Remark string `json:"Remark,omitempty" form:"Remark"`

	// PlatformID 特約合作平台商代號
	PlatformID string `json:"PlatformID,omitempty" form:"PlatformID"`

	// CheckMacValue 檢查碼
	CheckMacValue string `json:"CheckMacValue,omitempty" form:"CheckMacValue"`

	// RtnCode 目前物流狀態
	RtnCode string `json:"RtnCode,omitempty" form:"RtnCode"`

	// RtnMsg 物流狀態說明
	RtnMsg string `json:"RtnMsg,omitempty" form:"RtnMsg"`

	// AllPayLogisticsID 綠界科技的物流交易編號
	AllPayLogisticsID string `json:"AllPayLogisticsID,omitempty" form:"AllPayLogisticsID"`

	// UpdateStatusDate 物流狀態更新時間
	UpdateStatusDate string `json:"UpdateStatusDate,omitempty" form:"UpdateStatusDate"`

	// ReceiverAddress 收件人地址
	ReceiverAddress string `json:"ReceiverAddress,omitempty" form:"ReceiverAddress"`

	// CVSPaymentNo 寄貨編號
	CVSPaymentNo string `json:"CVSPaymentNo,omitempty" form:"CVSPaymentNo"`

	// CVSValidationNo 驗證碼
	CVSValidationNo string `json:"CVSValidationNo,omitempty" form:"CVSValidationNo"`

	// BookingNote 托運單號
	BookingNote string `json:"BookingNote,omitempty" form:"BookingNote"`
}

// Map is a function that maps the ECPayLogistics struct to the ECPayClient struct
func (e *ECPayLogistics) Map(c client.ECPayClient) (string, error) {

	return "", nil
}

// CreateExpress 綠界物流門市訂單建立
func (e *ECPayLogistics) CreateExpress(c client.ECPayClient) error {

	formData := helpers.LogisticsToFormValues(e)

	checkMacValue := helpers.GenerateCheckMacValue(formData, c.HashKey, c.HashIV)
	formData.Set("CheckMacValue", checkMacValue)

	resp, err := http.PostForm(c.BaseURL, formData)
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
