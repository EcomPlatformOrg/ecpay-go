package ecpayGo

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/url"
	"sort"
	"strconv"
	"strings"
)

// tradeToFormValues converts ECPayTrade to url.Values
func tradeToFormValues(trade *ECPayTrade) url.Values {

	formData := url.Values{}
	formData.Set("MerchantID", trade.MerchantID)
	formData.Set("MerchantTradeNo", trade.MerchantTradeNo)
	formData.Set("MerchantTradeDate", trade.MerchantTradeDate)
	formData.Set("PaymentType", trade.PaymentType)
	formData.Set("TotalAmount", strconv.Itoa(trade.TotalAmount))
	formData.Set("TradeDesc", trade.TradeDesc)
	formData.Set("ItemName", trade.ItemName)
	formData.Set("ReturnURL", trade.ReturnURL)
	formData.Set("ChoosePayment", trade.ChoosePayment)
	formData.Set("EncryptType", strconv.Itoa(trade.EncryptType))

	if trade.StoreID != "" {
		formData.Set("StoreID", trade.StoreID)
	}
	if trade.ClientBackURL != "" {
		formData.Set("ClientBackURL", trade.ClientBackURL)
	}
	if trade.ItemURL != "" {
		formData.Set("ItemURL", trade.ItemURL)
	}
	if trade.Remark != "" {
		formData.Set("Remark", trade.Remark)
	}
	if trade.ChooseSubPayment != "" {
		formData.Set("ChooseSubPayment", trade.ChooseSubPayment)
	}
	if trade.OrderResultURL != "" {
		formData.Set("OrderResultURL", trade.OrderResultURL)
	}
	if trade.NeedExtraPaidInfo != "" {
		formData.Set("NeedExtraPaidInfo", trade.NeedExtraPaidInfo)
	}
	if trade.IgnorePayment != "" {
		formData.Set("IgnorePayment", trade.IgnorePayment)
	}
	if trade.PlatformID != "" {
		formData.Set("PlatformID", trade.PlatformID)
	}
	if trade.CustomField1 != "" {
		formData.Set("CustomField1", trade.CustomField1)
	}
	if trade.CustomField2 != "" {
		formData.Set("CustomField2", trade.CustomField2)
	}
	if trade.CustomField3 != "" {
		formData.Set("CustomField3", trade.CustomField3)
	}
	if trade.CustomField4 != "" {
		formData.Set("CustomField4", trade.CustomField4)
	}
	if trade.Language != "" {
		formData.Set("Language", trade.Language)
	}

	return formData
}

// generateCheckMacValue generates CheckMacValue
func generateCheckMacValue(values url.Values, hashKey string, hashIV string) string {
	// Step (1) 將傳遞參數依照第一個英文字母，由A到Z的順序來排序
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	var sortedQueryString string
	for _, key := range keys {
		sortedQueryString += key + "=" + values.Get(key) + "&"
	}
	// remove trailing '&'
	sortedQueryString = strings.TrimSuffix(sortedQueryString, "&")

	// Step (2) 參數最前面加上HashKey、最後面加上HashIV
	encodedString := hashKey + "&" + sortedQueryString + "&" + hashIV

	// Step (3) 將整串字串進行URL encode
	encodedString = url.QueryEscape(encodedString)

	// Step (4) 轉為小寫
	encodedString = strings.ToLower(encodedString)

	// Step (5) 以SHA256加密方式來產生雜凑值
	hasher := sha256.New()
	hasher.Write([]byte(encodedString))
	hashedValue := hex.EncodeToString(hasher.Sum(nil))

	// Step (6) 再轉大寫產生CheckMacValue
	return strings.ToUpper(hashedValue)
}

// ValidateCheckMacValue validates the CheckMacValue from ECPay's response
func ValidateCheckMacValue(responseValues url.Values, hashKey string, hashIV string) error {
	// Extract the CheckMacValue from the response
	receivedCheckMacValue := responseValues.Get("CheckMacValue")
	if receivedCheckMacValue == "" {
		return errors.New("CheckMacValue is missing from the response")
	}

	// Remove the CheckMacValue from the values as it's not included in the generation
	responseValues.Del("CheckMacValue")

	// Generate the CheckMacValue using the provided values
	expectedCheckMacValue := generateCheckMacValue(responseValues, hashKey, hashIV)

	if expectedCheckMacValue != receivedCheckMacValue {
		return errors.New("CheckMacValue mismatch")
	}

	return nil
}
