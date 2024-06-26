package logistics

import (
	"bytes"
	"encoding/json"
	"errors"
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

	// TempLogisticsID 物流交易編號
	TempLogisticsID string `json:"TempLogisticsID,omitempty" form:"TempLogisticsID"`

	// LogisticsID 綠界物流訂單編號
	LogisticsID string `json:"LogisticsID,omitempty"`

	// LogisticsStatus 物流狀態
	LogisticsStatus string `json:"LogisticsStatus,omitempty"`

	// LogisticsStatusName 貨態代碼訊息
	LogisticsStatusName string `json:"LogisticsStatusName,omitempty"`

	// LogisticsType 物流類型
	LogisticsType string `json:"LogisticsType,omitempty" form:"LogisticsType"`

	// LogisticsSubType 物流子類型
	LogisticsSubType string `json:"LogisticsSubType,omitempty" form:"LogisticsSubType"`

	// LogisticsSubType 物流子類型
	LogisticsSelection string `json:"LogisticsURL,omitempty"`

	// CollectionAmount 代收金額
	CollectionAmount int `json:"CollectionAmount,omitempty" form:"CollectionAmount"`

	// IsCollection 是否代收貨款
	IsCollection string `json:"IsCollection,omitempty" form:"IsCollection"`

	// Temperature 溫層
	Temperature string `json:"Temperature,omitempty" form:"Temperature"`

	// Specification 規格
	Specification string `json:"Specification,omitempty" form:"Specification"`

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
	RtnCode int `json:"RtnCode,omitempty" form:"RtnCode"`

	// RtnMsg 物流狀態說明
	RtnMsg string `json:"RtnMsg,omitempty" form:"RtnMsg"`

	// ScheduledPickupTime 預定取件時段
	ScheduledPickupTime string `json:"ScheduledPickupTime,omitempty" json:"ScheduledPickupTime"`

	// EnableSelectDeliveryTime 是否允許選擇送達時間
	EnableSelectDeliveryTime string `json:"EnableSelectDeliveryTime,omitempty" form:"EnableSelectDeliveryTime"`

	// RqHeader
	RqHeader model.RqHeader `json:"RqHeader"`

	TransCode int `json:"TransCode"`

	TransMsg string `json:"TransMsg"`

	// Data
	Data string `json:"Data"`

	// ResultData
	ResultData string `json:"ResultData"`

	// BaseModel 通用參數
	model.BaseModel `json:",inline"`

	// GoodsAmount 商品金額
	model.Merchant `json:",inline"`

	// Goods 商品資訊
	model.Goods `json:",inline"`

	// Sender 寄件人資訊
	model.Sender `json:",inline"`

	// Receiver 收件人資訊
	model.Receiver `json:",inline"`

	// ConvenienceStore 超商取貨相關資訊
	model.ConvenienceStore `json:",inline"`
}

// Map is a function that maps the ECPayLogistics struct to the ECPayClient struct
func (e *ECPayLogistics) Map() (string, error) {

	formData := helpers.ReflectFormValues(e)

	body, err := helpers.SendFormData(e.Client, formData)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// CreateExpress 綠界物流門市訂單建立
func (e *ECPayLogistics) CreateExpress() error {

	formData := helpers.ReflectFormValues(e)

	checkMacValue := helpers.GenerateCheckMacValue(formData, e.Client.HashKey, e.Client.HashIV)
	formData.Set("CheckMacValue", checkMacValue)

	body, err := helpers.SendFormData(e.Client, formData)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(body, &e); err != nil {
		return err
	}

	return nil
}

// EncryptLogistics is a method that encrypts the ECPayLogistics struct using the helpers.EncryptData function and sets the encrypted data to the "Data" field of the struct.
// It takes no arguments and returns an error if there was an error marshalling the struct or encrypting the data, otherwise it returns nil.
func (e *ECPayLogistics) EncryptLogistics() error {

	jsonBytes, err := json.Marshal(e)
	if err != nil {
		slog.Error(fmt.Sprintf("Error marshalling ECPayLogistics struct: %v", err))
		return err
	}

	jsonString := string(jsonBytes)
	encryptedData, err := helpers.EncryptData(jsonString, e.Client.HashKey, e.Client.HashIV)
	if err != nil {
		slog.Error(fmt.Sprintf("Error encrypting data: %v", err))
		return err
	}

	e.Data = encryptedData
	return nil
}

func (e *ECPayLogistics) DecryptLogistics(body []byte) error {

	if err := json.NewDecoder(bytes.NewReader(body)).Decode(&e); err != nil {
		slog.Error(fmt.Sprintf("Error decoding response body: %v", err))
		return err
	}

	slog.Info(fmt.Sprintf("Body : %s", string(body)))
	slog.Info(fmt.Sprintf("TransCode : %d", e.TransCode))
	slog.Info(fmt.Sprintf("TransMsg : %s", e.TransMsg))
	slog.Info(fmt.Sprintf("Data : %s", e.Data))
	decryptedDataString, err := helpers.DecryptData(e.Data, e.Client.HashKey, e.Client.HashIV)
	if err = json.NewDecoder(bytes.NewReader([]byte(decryptedDataString))).Decode(&e); err != nil {
		slog.Error(fmt.Sprintf("Error decoding decrypted data: %v", err))
		return err
	}

	return nil
}

func (e *ECPayLogistics) RedirectToLogisticsSelection() (string, error) {

	if err := e.EncryptLogistics(); err != nil {
		return "", err
	}

	e.RqHeader = model.RqHeader{
		Timestamp: strconv.FormatInt(time.Now().Unix(), 10),
	}

	jsonData, err := json.Marshal(e)
	if err != nil {
		slog.Error(fmt.Sprintf("Error marshalling ECPayLogistics struct: %v", err))
		return "", err
	}

	resp, err := http.Post(e.Client.BaseURL, "application/json", strings.NewReader(string(jsonData)))
	if err != nil {
		slog.Error(fmt.Sprintf("Error sending POST request: %v", err))
		return "", err
	}

	defer func(Body io.ReadCloser) {
		if err = Body.Close(); err != nil {
			slog.Error(err.Error())
		}
	}(resp.Body)

	responseData := &ECPayLogistics{
		BaseModel: model.BaseModel{
			Client: &client.ECPayClient{
				HashKey: e.Client.HashKey,
				HashIV:  e.Client.HashIV,
			},
		},
	}

	body, _ := io.ReadAll(resp.Body)
	if err = responseData.DecryptLogistics(body); err != nil {
		return string(body), err
	}

	return responseData.RtnMsg, nil
}

func (e *ECPayLogistics) UpdateTempTrade() error {

	if err := e.EncryptLogistics(); err != nil {
		return err
	}

	e.RqHeader = model.RqHeader{
		Timestamp: strconv.FormatInt(time.Now().Unix(), 10),
	}

	jsonData, err := json.Marshal(e)
	if err != nil {
		slog.Error(fmt.Sprintf("Error marshalling ECPayLogistics struct: %v", err))
		return err
	}

	resp, err := http.Post(e.Client.BaseURL, "application/json", strings.NewReader(string(jsonData)))
	if err != nil {
		slog.Error(fmt.Sprintf("Error sending POST request: %v", err))
		return err
	}

	defer func(Body io.ReadCloser) {
		if err = Body.Close(); err != nil {
			slog.Error(err.Error())
		}
	}(resp.Body)

	responseData := &ECPayLogistics{
		BaseModel: model.BaseModel{
			Client: &client.ECPayClient{
				HashKey: e.Client.HashKey,
				HashIV:  e.Client.HashIV,
			},
		},
	}

	body, _ := io.ReadAll(resp.Body)
	if err = responseData.DecryptLogistics(body); err != nil {
		return err
	}

	if responseData.RtnCode != 1 {
		slog.Info(fmt.Sprintf("responseData.RtnCode : %d", responseData.RtnCode))
		slog.Info(fmt.Sprintf("responseData.RtnMsg : %s", responseData.RtnMsg))
		return errors.New(fmt.Sprintf("建立暫存物流訂單失敗 失敗原因 : %s", responseData.RtnMsg))
	}

	return nil
}

func (e *ECPayLogistics) CreateByTempTrade() (string, error) {

	if err := e.EncryptLogistics(); err != nil {
		return "", err
	}

	e.RqHeader = model.RqHeader{
		Timestamp: strconv.FormatInt(time.Now().Unix(), 10),
	}

	jsonData, err := json.Marshal(e)
	if err != nil {
		slog.Error(fmt.Sprintf("Error marshalling ECPayLogistics struct: %v", err))
		return "", err
	}

	resp, err := http.Post(e.Client.BaseURL, "application/json", strings.NewReader(string(jsonData)))
	if err != nil {
		slog.Error(fmt.Sprintf("Error sending POST request: %v", err))
		return "", err
	}

	defer func(Body io.ReadCloser) {
		if err = Body.Close(); err != nil {
			slog.Error(err.Error())
		}
	}(resp.Body)

	responseData := &ECPayLogistics{
		BaseModel: model.BaseModel{
			Client: &client.ECPayClient{
				HashKey: e.Client.HashKey,
				HashIV:  e.Client.HashIV,
			},
		},
	}

	body, _ := io.ReadAll(resp.Body)
	if err = responseData.DecryptLogistics(body); err != nil {
		return "", err
	}

	if responseData.RtnCode != 1 {
		return "", errors.New(fmt.Sprintf("建立正式物流訂單失敗 失敗原因 : %s", responseData.RtnMsg))
	}

	return responseData.LogisticsID, nil
}

//func (e *ECPayLogistics)CreateTestData() (string, error) {
//
//}
