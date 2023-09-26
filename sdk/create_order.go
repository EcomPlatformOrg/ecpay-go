package sdk

import (
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"

	"github.com/linshan1234/ECPay-go/client"
	"github.com/linshan1234/ECPay-go/models"
)

type Parameter map[string]string

func CreateOrder(orderData models.Order) error {

	initClient := &client.ECPayClient{
		BaseURL:    "https://payment-stage.ecpay.com.tw/Cashier/AioCheckOut/V5",
		MerchantID: "3002607",
		HashKey:    "pwFHCqoQZGmho4w6",
		HashIV:     "EkRm7iFT261dpevs",
	}

	initPara := Parameter{
		"MerchantID":      initClient.MerchantID,
		"MerchantTradeNo": "test" + strconv.Itoa(orderData.ID),
		// "MerchantTradeDate": time.Now().Format("2006/01/02 15:04:05"),
		"MerchantTradeDate": "2023/09/25 00:00:00",
		"PaymentType":       "aio",
		"TotalAmount":       strconv.Itoa(orderData.TotalAmount),
		"TradeDesc":         orderData.TradeDesc,
		"ItemName":          "A#B#C",
		"ReturnURL":         "https://www.ecpay.com.tw/return_url.php",
		"ChoosePayment":     "Credit",
		"EncryptType":       "1",
		//非必填
		"ClientBackURL":     "https://www.ecpay.com.tw/client_back_url.php",
		"Remark":            "交易備註",
		"OrderResultURL":    "https://www.ecpay.com.tw/order_result_url.php",
		"NeedExtraPaidInfo": "Y",
		// "ItemURL":           "",
		// "ChooseSubPayment":  "",
		// "StoreID":           "",
		// "IgnorePayment":     "",
		// "PlatformID":        "",
		// "CustomField1":      "",
		// "CustomField2":      "",
		// "CustomField3":      "",
		// "CustomField4":      "",
		// "Language":          "",
	}

	//轉成formdata
	formDataString := ""
	for key, value := range initPara {
		formDataString += key + "=" + value + "&"
	}

	CheckMacValue := GenerateCheckValue(initClient, initPara)

	//加入Hash
	finalPara := formDataString + "CheckMacValue=" + CheckMacValue
	fmt.Println("final formdata:", finalPara)

	// 發送至綠界
	client := &http.Client{}
	req, err := http.NewRequest("POST", initClient.BaseURL, strings.NewReader(finalPara))
	if err != nil {
		fmt.Print("創建請求錯誤")
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Print("發送請求錯誤", err)
	} else {
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("讀取回應時發生錯誤:", err)
		} else {
			fmt.Println(string(body))
		}
	}
	return nil

}

func GenerateCheckValue(client *client.ECPayClient, initPara map[string]string) (CheckMacValue string) {
	//排序初始化參數
	var keys []string
	for key := range initPara {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	initformData := ""
	for i := 0; i < len(keys); i++ {
		initformData += keys[i] + "=" + initPara[keys[i]] + "&"
	}
	initformData = strings.TrimSuffix(initformData, "&")
	// fmt.Println("initformData:", initformData)
	// initformDataString := initformData.Encode()

	//加上HashKey、HashIV
	finalPara := "HashKey=" + client.HashKey + "&" + initformData + "&" + "HashIV=" + client.HashIV
	fmt.Println("Hash initformData:", finalPara)

	// url encode
	initformData = url.QueryEscape(finalPara)
	fmt.Println("urlencode formdata:", initformData)

	//轉為小寫
	initformData = strings.ToLower(initformData)
	//SHA256加密

	hash := sha256.Sum256([]byte(initformData))
	encodedHash := fmt.Sprintf("%x", hash)

	hashString := strings.ToUpper(encodedHash)
	fmt.Println("hashString:", hashString)
	return hashString
}
