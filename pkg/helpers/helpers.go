package helpers

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/goccy/go-reflect"
	"log/slog"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

// TradeToFormValues converts ECPayTrade to url.Values
//func TradeToFormValues(trade *trade.ECPayTrade) url.Values {
//
//	formData := url.Values{}
//	formData.Set("MerchantID", trade.MerchantID)
//	formData.Set("MerchantTradeNo", trade.MerchantTradeNo)
//	formData.Set("MerchantTradeDate", trade.MerchantTradeDate)
//	formData.Set("PaymentType", trade.PaymentType)
//	formData.Set("TotalAmount", strconv.Itoa(trade.TotalAmount))
//	formData.Set("TradeDesc", trade.TradeDesc)
//	formData.Set("ItemName", trade.ItemName)
//	formData.Set("ReturnURL", trade.ReturnURL)
//	formData.Set("ChoosePayment", trade.ChoosePayment)
//	formData.Set("EncryptType", strconv.Itoa(trade.EncryptType))
//
//	if trade.StoreID != "" {
//		formData.Set("StoreID", trade.StoreID)
//	}
//	if trade.ClientBackURL != "" {
//		formData.Set("ClientBackURL", trade.ClientBackURL)
//	}
//	if trade.ItemURL != "" {
//		formData.Set("ItemURL", trade.ItemURL)
//	}
//	if trade.Remark != "" {
//		formData.Set("Remark", trade.Remark)
//	}
//	if trade.ChooseSubPayment != "" {
//		formData.Set("ChooseSubPayment", trade.ChooseSubPayment)
//	}
//	if trade.OrderResultURL != "" {
//		formData.Set("OrderResultURL", trade.OrderResultURL)
//	}
//	if trade.NeedExtraPaidInfo != "" {
//		formData.Set("NeedExtraPaidInfo", trade.NeedExtraPaidInfo)
//	}
//	if trade.IgnorePayment != "" {
//		formData.Set("IgnorePayment", trade.IgnorePayment)
//	}
//	if trade.PlatformID != "" {
//		formData.Set("PlatformID", trade.PlatformID)
//	}
//	if trade.CustomField1 != "" {
//		formData.Set("CustomField1", trade.CustomField1)
//	}
//	if trade.CustomField2 != "" {
//		formData.Set("CustomField2", trade.CustomField2)
//	}
//	if trade.CustomField3 != "" {
//		formData.Set("CustomField3", trade.CustomField3)
//	}
//	if trade.CustomField4 != "" {
//		formData.Set("CustomField4", trade.CustomField4)
//	}
//	if trade.Language != "" {
//		formData.Set("Language", trade.Language)
//	}
//
//	return formData
//}
//
//// LogisticsToFormValues 將 ECPayLogistics 結構體轉換為 url.Values。
//func LogisticsToFormValues(logistics *express.ECPayLogistics) url.Values {
//	formData := url.Values{}
//
//	// 將 ECPayLogistics 結構體中的所有字段添加到 formData 中
//	formData.Set("MerchantID", logistics.MerchantID)
//	formData.Set("MerchantTradeNo", logistics.MerchantTradeNo)
//	formData.Set("MerchantTradeDate", logistics.MerchantTradeDate)
//	formData.Set("LogisticsType", logistics.LogisticsType)
//	formData.Set("LogisticsSubType", logistics.LogisticsSubType)
//	formData.Set("GoodsAmount", strconv.Itoa(logistics.GoodsAmount))
//	formData.Set("CollectionAmount", strconv.Itoa(logistics.CollectionAmount))
//	formData.Set("IsCollection", logistics.IsCollection)
//	formData.Set("GoodsName", logistics.GoodsName)
//	formData.Set("SenderName", logistics.SenderName)
//	formData.Set("SenderPhone", logistics.SenderPhone)
//	formData.Set("SenderCellPhone", logistics.SenderCellPhone)
//	formData.Set("ReceiverName", logistics.ReceiverName)
//	formData.Set("ReceiverPhone", logistics.ReceiverPhone)
//	formData.Set("ReceiverCellPhone", logistics.ReceiverCellPhone)
//	formData.Set("ReceiverEmail", logistics.ReceiverEmail)
//	formData.Set("ReceiverStoreID", logistics.ReceiverStoreID)
//	formData.Set("ReturnStoreID", logistics.ReturnStoreID)
//	formData.Set("TradeDesc", logistics.TradeDesc)
//	formData.Set("ServerReplyURL", logistics.ServerReplyURL)
//	formData.Set("ClientReplyURL", logistics.ClientReplyURL)
//	formData.Set("Remark", logistics.Remark)
//	formData.Set("PlatformID", logistics.PlatformID)
//	formData.Set("CheckMacValue", logistics.CheckMacValue)
//	formData.Set("RtnCode", logistics.RtnCode)
//	formData.Set("RtnMsg", logistics.RtnMsg)
//	formData.Set("AllPayLogisticsID", logistics.AllPayLogisticsID)
//	formData.Set("UpdateStatusDate", logistics.UpdateStatusDate)
//	formData.Set("ReceiverAddress", logistics.ReceiverAddress)
//	formData.Set("CVSPaymentNo", logistics.CVSPaymentNo)
//	formData.Set("CVSValidationNo", logistics.CVSValidationNo)
//	formData.Set("BookingNote", logistics.BookingNote)
//
//	return formData
//}

func ReflectFormValues(data any) url.Values {

	values := url.Values{}
	v := reflect.ValueOf(data)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		tag := v.Type().Field(i).Tag.Get("form")

		if tag == "" || !field.IsValid() {
			continue
		}

		var value string
		switch field.Kind() {
		case reflect.String:
			strVal := field.String()
			if strVal != "" {
				value = strVal
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if field.Int() != 0 {
				value = strconv.FormatInt(field.Int(), 10)
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if field.Uint() != 0 {
				value = strconv.FormatUint(field.Uint(), 10)
			}
		case reflect.Float32, reflect.Float64:
			if field.Float() != 0.0 {
				value = strconv.FormatFloat(field.Float(), 'f', -1, 64)
			}
		case reflect.Bool:
			value = strconv.FormatBool(field.Bool())
		case reflect.Struct:
			if field.Type() == reflect.TypeOf(time.Time{}) {
				t := field.Interface().(time.Time)
				if !t.IsZero() {
					value = t.Format(time.RFC3339)
				}
			}
		case reflect.Ptr:
			if !field.IsNil() {
				value = fmt.Sprintf("%v", field.Elem().Interface())
			}
		default:
			slog.Error(fmt.Sprintf("Unsupported type: %v", field.Kind()))
		}
		values.Set(tag, value)
	}

	return values
}

// GenerateCheckMacValue generates CheckMacValue
func GenerateCheckMacValue(values url.Values, hashKey string, hashIV string) string {

	// Step (1) 將傳遞參數依照第一個英文字母，由A到Z的順序來排序
	slog.Info(fmt.Sprintf("Step (1) values: %v", values))
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	slog.Info(fmt.Sprintf("keys: %v", keys))

	var sortedQueryString string
	for _, key := range keys {
		sortedQueryString += key + "=" + values.Get(key) + "&"
	}
	// remove trailing '&'
	sortedQueryString = strings.TrimSuffix(sortedQueryString, "&")
	slog.Info(fmt.Sprintf("sortedQueryString: %v", sortedQueryString))

	// Step (2) 參數最前面加上HashKey、最後面加上HashIV
	encodedString := "HashKey=" + hashKey + "&" + sortedQueryString + "&HashIV=" + hashIV
	slog.Info(fmt.Sprintf("Step (2) encodedString: %v", encodedString))

	// Step (3) 將整串字串進行URL encode
	encodedString = url.QueryEscape(encodedString)
	slog.Info(fmt.Sprintf("Step (3) encodedString: %v", encodedString))

	// Step (4) 轉為小寫
	encodedString = strings.ToLower(encodedString)
	slog.Info(fmt.Sprintf("Step (4) encodedString: %v", encodedString))

	// Step (5) 以SHA256加密方式來產生雜凑值
	hasher := sha256.New()
	hasher.Write([]byte(encodedString))
	hashedValue := hex.EncodeToString(hasher.Sum(nil))
	slog.Info(fmt.Sprintf("Step (5) hashedValue: %v", hashedValue))

	// Step (6) 再轉大寫產生CheckMacValue
	slog.Info(fmt.Sprintf("Step (6) CheckMacValue: %v", strings.ToUpper(hashedValue)))
	return strings.ToUpper(hashedValue)
}
