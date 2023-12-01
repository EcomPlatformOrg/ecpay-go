package model

type BaseModel struct {
	// TradeDesc 交易描述
	TradeDesc string `json:"TradeDesc,omitempty" form:"TradeDesc"`

	// Remark 備註
	Remark string `json:"Remark,omitempty" form:"Remark"`

	// PlatformID 特約合作平台商代號
	PlatformID string `json:"PlatformID,omitempty" form:"PlatformID"`

	// CheckMacValue 檢查碼
	CheckMacValue string `json:"CheckMacValue,omitempty" form:"CheckMacValue"`
}

type Merchant struct {
	// MerchantID 特店編號
	MerchantID string `json:"MerchantID,omitempty" form:"MerchantID"`

	// MerchantTradeNo 特店交易編號
	MerchantTradeNo string `json:"MerchantTradeNo,omitempty" form:"MerchantTradeNo"`

	// MerchantTradeDate 廠商交易時間
	MerchantTradeDate string `json:"MerchantTradeDate,omitempty" form:"MerchantTradeDate"`
}

type ConvenienceStore struct {
	// CVSPaymentNo 寄貨編號
	CVSPaymentNo string `json:"CVSPaymentNo,omitempty" form:"CVSPaymentNo"`

	// CVSStoreID 寄貨門市代號
	CVSStoreID string `json:"CVSStoreID,omitempty" form:"CVSStoreID"`

	// CVSStoreName 寄貨門市名稱
	CVSStoreName string `json:"CVSStoreName,omitempty" form:"CVSStoreName"`

	// CVSAddress 寄貨門市地址
	CVSAddress string `json:"CVSAddress,omitempty" form:"CVSAddress"`

	// CVSOutSide 寄貨門市是否為外縣市
	CVSOutSide string `json:"CVSOutSide,omitempty" form:"CVSOutSide"`

	// CVSValidationNo 驗證碼
	CVSValidationNo string `json:"CVSValidationNo,omitempty" form:"CVSValidationNo"`
}

type Goods struct {
	// GoodsName 商品名稱
	GoodsName string `json:"GoodsName,omitempty" form:"GoodsName"`

	// GoodsAmount 商品金額
	GoodsAmount int `json:"GoodsAmount,omitempty" form:"GoodsAmount"`
}

type Sender struct {
	// SenderName 寄件人姓名
	SenderName string `json:"SenderName,omitempty" form:"SenderName"`

	// SenderPhone 寄件人電話
	SenderPhone string `json:"SenderPhone,omitempty" form:"SenderPhone"`

	// SenderCellPhone 寄件人手機
	SenderCellPhone string `json:"SenderCellPhone,omitempty" form:"SenderCellPhone"`
}

type Receiver struct {
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

	// ReceiverAddress 收件人地址
	ReceiverAddress string `json:"ReceiverAddress,omitempty" form:"ReceiverAddress"`
}
