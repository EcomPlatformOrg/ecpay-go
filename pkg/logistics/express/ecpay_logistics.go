package logistics

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/EcomPlatformOrg/ecpay-go/pkg/client"
	"github.com/EcomPlatformOrg/ecpay-go/pkg/helpers"
	"github.com/EcomPlatformOrg/ecpay-go/pkg/model"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"
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

	// ServiceType 服務型態 固定帶4
	ServiceType string `json:"ServiceType,omitempty" form:"ServiceType"`

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

	// RqHeader
	RqHeader model.RqHeader `json:",inline"`

	// Data
	Data string `json:"Data"`

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

// CreateTestData is a method that creates test data for the ECPayLogistics struct and sends it to the ECPayClient server for processing and decryption. The method returns the decrypted
func (e *ECPayLogistics) CreateTestData(c client.ECPayClient) (*ECPayLogistics, error) {

	jsonBytes, err := json.Marshal(e)
	if err != nil {
		slog.Error(fmt.Sprintf("Error marshalling ECPayLogistics struct: %v", err))
		return nil, err
	}

	jsonString := string(jsonBytes)
	encryptedData, err := helpers.EncryptData(jsonString, c.HashKey, c.HashIV)
	if err != nil {
		slog.Error(fmt.Sprintf("Error encrypting data: %v", err))
		return nil, err
	}

	e.Data = encryptedData
	e.RqHeader = model.RqHeader{
		Timestamp: strconv.FormatInt(time.Now().Unix(), 10),
	}

	jsonData, err := json.Marshal(e)
	if err != nil {
		slog.Error(fmt.Sprintf("Error marshalling ECPayLogistics struct: %v", err))
		return nil, err
	}

	resp, err := http.Post(c.BaseURL, "application/json", strings.NewReader(string(jsonData)))
	if err != nil {
		slog.Error(fmt.Sprintf("Error sending POST request: %v", err))
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		if err = Body.Close(); err != nil {
			slog.Error(err.Error())
		}
	}(resp.Body)

	responseData := ECPayLogistics{}
	if err = json.NewDecoder(resp.Body).Decode(&responseData); err != nil {
		slog.Error(fmt.Sprintf("Error decoding response body: %v", err))
		return nil, err
	}

	decryptedData := &ECPayLogistics{}
	decryptedDataString, err := helpers.DecryptData(responseData.Data, c.HashKey, c.HashIV)
	if err = json.NewDecoder(bytes.NewReader([]byte(decryptedDataString))).Decode(&decryptedData); err != nil {
		slog.Error(fmt.Sprintf("Error decoding decrypted data: %v", err))
		return nil, err
	}

	return decryptedData, nil
}
