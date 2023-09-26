package client

import (
	"net/http"
)

type ECPayClient struct {
	BaseURL    string       `json:"BaseURL"`
	MerchantID string       `json:"MerchantID"`
	HashKey    string       `json:"HashKey"`
	HashIV     string       `json:"HashIV"`
	Client     *http.Client `json:"Client"`
}

// NewECPayClient 初始化一個新的 ECPayClient
func NewECPayClient(baseURL, merchantID, hashKey, hashIV string) *ECPayClient {
	return &ECPayClient{
		BaseURL:    baseURL,
		MerchantID: merchantID,
		HashKey:    hashKey,
		HashIV:     hashIV,
		Client:     &http.Client{},
	}
}
