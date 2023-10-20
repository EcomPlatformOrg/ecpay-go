package ecpayGo

type ECPayClient struct {
	BaseURL string `json:"BaseURL"`
	HashKey string `json:"HashKey"`
	HashIV  string `json:"HashIV"`
}
