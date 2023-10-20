# ECpay-go

## 綠界SDK - Golang

[綠界特店管理後台測試環境](https://vendor-stage.ecpay.com.tw/)

### 主要功能：

1. **查詢訂單**
2. **模擬付款**：用於測試特店的主機(server)是否能接收綠界後端發送的付款結果通知。  
   **模擬付款功能路徑**：
    - 進入綠界特店管理後台
    - 選擇「一般訂單查詢」
    - 點選「全方位金流訂單」，並輸入查詢條件
    - 根據查詢結果的訂單點擊「模擬付款」按鈕

### 測試用信用卡資料：

- **一般信用卡**
    - 卡號: `4311-9522-2222-2222`
    - 安全碼: `222`

- **圓夢彈性分期信用卡** (用於測試圓夢分期服務，詳細圓夢分期服務請參考)
    - 卡號: `4938-1777-7777-7777`
    - 安全碼: `222`
# ECPay 交易模組

本文件提供了 ECPay 交易模組的詳細說明。

## 1. 結構定義: ECPayTrade

`ECPayTrade` 結構提供了建立 ECPay 交易所需的所有必要欄位。

- **MerchantID**: 特店編號
- **MerchantTradeNo**: 特店訂單編號
- **MerchantTradeDate**: 特店交易時間
- **PaymentType**: 交易類型
  ... [其他欄位]

## 2. 方法說明

### 2.1 CreateAioPayment

這是一個方法，用於建立 ECPay 的 AIO 交易。

**使用方法**:

```go
trade := ECPayTrade{...} // 初始化交易資料
client := ECPayClient{...} // 初始化客戶端資料
err := trade.CreateAioPayment(client)
if err != nil {
    // 處理錯誤
}
```

### 2.2 tradeToFormValues

這個方法將 ECPayTrade 物件轉換為 HTTP POST 請求適用的 url.Values 格式。

### 2.3 generateCheckMacValue

此方法專門用於生成 ECPay 的 CheckMacValue，以確保交易資料的安全性。它首先依據英文字母順序對參數進行排序，然後附加特定的加密金鑰和初始向量，接著對其進行 URL 編碼，再進行 SHA256 雜湊，最後將其轉換為大寫。


