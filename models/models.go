package models

import "time"

// map[string]interface{}
type Parameters struct {
	MerchantID        string
	MerchantTradeNo   string
	MerchantTradeDate string // yyyy/MM/dd HH:mm:ss
	PaymentType       string //aio
	TotalAmount       int
	TradeDesc         string
	ItemName          string
	ReturnURL         string
	ChoosePayment     string
	CheckMacValue     string
	EncryptType       int // 1
	StoreID           string
	ClientBackURL     string
	ItemURL           string
	Remark            string
	ChooseSubPayment  string
	OrderResultURL    string
	NeedExtraPaidInfo string //N
	IgnorePayment     string
	PlatformID        string
	CustomField1      string
	CustomField2      string
	CustomField3      string
	CustomField4      string
	Language          string //ENG KOR JPN CHI
}

type Order struct {
	ID                int               `json:"id"`                  // ID Unique identifier for the order
	UserID            string            `json:"user_id"`             // UserID of the user who placed the order
	ShippingAddressID *int              `json:"shipping_address_id"` // ShippingAddressID of the shipping address
	BillingAddressID  *int              `json:"billing_address_id"`  // BillingAddressID of the billing address
	DateCreated       *time.Time        `json:"date_created"`        // DateCreated Timestamp of the order's creation date
	Status            string            `json:"status"`              // Status of the order (e.g., "In Progress")
	ShippingAddress   map[string]string `json:"shipping_address,omitempty"`
	BillingAddress    map[string]string `json:"billing_address,omitempty"`
	Products          []Category        `json:"products,omitempty"`
	//需要加的欄位
	TotalAmount int    `json:"totalAmount,omitempty"`
	TradeDesc   string `json:"tradeDesc,omitempty"`
}

type Category struct {
	ID               int         `json:"id"`                           // ID Unique identifier for the category
	Name             string      `json:"name"`                         // Name of the category
	Description      string      `json:"description"`                  // Description of the category
	ParentCategoryID *int        `json:"parent_category_id,omitempty"` // ParentCategoryID of the parent category (allows for nested categories)
	SubCategories    []*Category `json:"sub_categories,omitempty"`
}
