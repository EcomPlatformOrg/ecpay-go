package logistics

import (
	"ecpay-go/pkg/client"
	"ecpay-go/pkg/helpers"
	"ecpay-go/pkg/model"
	"encoding/json"
)

// ECPayLogistics is a struct containing information for an ECPay logistics
type ECPayLogistics struct {

	// LogisticsType 物流類型
	LogisticsType string `json:"LogisticsType,omitempty" form:"LogisticsType"`

	// LogisticsSubType 物流子類型
	LogisticsSubType string `json:"LogisticsSubType,omitempty" form:"LogisticsSubType"`

	// CollectionAmount 代收金額
	CollectionAmount int `json:"CollectionAmount,omitempty" form:"CollectionAmount"`

	// IsCollection 是否代收貨款
	IsCollection string `json:"IsCollection,omitempty" form:"IsCollection"`

	// ReturnStoreID 退貨門市代號
	ReturnStoreID string `json:"ReturnStoreID,omitempty" form:"ReturnStoreID"`

	// AllPayLogisticsID 綠界科技的物流交易編號
	AllPayLogisticsID string `json:"AllPayLogisticsID,omitempty" form:"AllPayLogisticsID"`

	// UpdateStatusDate 物流狀態更新時間
	UpdateStatusDate string `json:"UpdateStatusDate,omitempty" form:"UpdateStatusDate"`

	// ServerReplyURL Server端回覆網址
	ServerReplyURL string `json:"ServerReplyURL,omitempty" form:"ServerReplyURL"`

	// ClientReplyURL Client端回覆網址
	ClientReplyURL string `json:"ClientReplyURL,omitempty" form:"ClientReplyURL"`

	// BookingNote 托運單號
	BookingNote string `json:"BookingNote,omitempty" form:"BookingNote"`

	// ExtraData 額外資訊
	ExtraData string `json:"ExtraData,omitempty" form:"ExtraData"`

	// RtnCode 目前物流狀態
	RtnCode string `json:"RtnCode,omitempty" form:"RtnCode"`

	// RtnMsg 物流狀態說明
	RtnMsg string `json:"RtnMsg,omitempty" form:"RtnMsg"`

	// BaseModel 通用參數
	BaseModel model.BaseModel `json:",inline"`

	// GoodsAmount 商品金額
	Merchant model.Merchant `json:",inline"`

	// Goods 商品資訊
	Goods model.Goods `json:",inline"`

	// Sender 寄件人資訊
	Sender model.Sender `json:",inline"`

	// Receiver 收件人資訊
	Receiver model.Receiver `json:",inline"`

	// ConvenienceStore 超商取貨相關資訊
	ConvenienceStore model.ConvenienceStore `json:",inline"`
}

// Map is a function that maps the ECPayLogistics struct to the ECPayClient struct
func (e *ECPayLogistics) Map(c client.ECPayClient) (string, error) {

	formData := helpers.ReflectFormValues(e)

	body, err := helpers.SendFormData(c, formData)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// CreateExpress 綠界物流門市訂單建立
func (e *ECPayLogistics) CreateExpress(c client.ECPayClient) error {

	formData := helpers.ReflectFormValues(e)

	checkMacValue := helpers.GenerateCheckMacValue(formData, c.HashKey, c.HashIV)
	formData.Set("CheckMacValue", checkMacValue)

	body, err := helpers.SendFormData(c, formData)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(body, &e); err != nil {
		return err
	}

	return nil

}
