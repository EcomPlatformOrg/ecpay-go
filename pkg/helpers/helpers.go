package helpers

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/EcomPlatformOrg/ecpay-go/pkg/client"
	"github.com/goccy/go-reflect"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

func ReflectFormValues(data any) url.Values {
	values := url.Values{}
	reflectStruct(reflect.ValueOf(data), &values)
	return values
}

func reflectStruct(v reflect.Value, values *url.Values) {
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := v.Type().Field(i)
		tag := fieldType.Tag.Get("form")

		if tag == "" {
			if field.Kind() == reflect.Struct && !isTimeStruct(field) {
				reflectStruct(field, values) // 遞迴處理嵌套結構體
			}
			continue
		}

		if !field.IsValid() || !field.CanInterface() {
			continue
		}

		switch field.Kind() {
		case reflect.String:
			strVal := field.String()
			if strVal != "" {
				values.Set(tag, strVal)
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if field.Int() != 0 {
				values.Set(tag, strconv.FormatInt(field.Int(), 10))
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if field.Uint() != 0 {
				values.Set(tag, strconv.FormatUint(field.Uint(), 10))
			}
		case reflect.Float32, reflect.Float64:
			if field.Float() != 0.0 {
				values.Set(tag, strconv.FormatFloat(field.Float(), 'f', -1, 64))
			}
		case reflect.Bool:
			values.Set(tag, strconv.FormatBool(field.Bool()))
		case reflect.Struct:
			if isTimeStruct(field) {
				t := field.Interface().(time.Time)
				if !t.IsZero() {
					values.Set(tag, t.Format(time.RFC3339))
				}
			}
		case reflect.Ptr:
			if !field.IsNil() {
				values.Set(tag, fmt.Sprintf("%v", field.Elem().Interface()))
			}
		default:
			// 處理不支持的類型
			slog.Error("Unsupported type: %v\n", field.Kind())
		}
	}
}

func isTimeStruct(v reflect.Value) bool {
	if v.Type().Name() == "Time" && v.Type().PkgPath() == "time" {
		return true
	}
	return false
}

// EncryptData 使用 ECPay 的加密方式對數據進行加密
func EncryptData(data string, hashKey string, hashIV string) (string, error) {
	// URL 編碼
	urlEncodedData := url.QueryEscape(data)

	block, err := aes.NewCipher([]byte(hashKey))
	if err != nil {
		return "", err
	}

	// PKCS7 Padding
	blockSize := block.BlockSize()
	padding := blockSize - len(urlEncodedData)%blockSize
	text := strings.Repeat(string(rune(padding)), padding)
	urlEncodedData += text

	// 初始化向量 IV
	iv := []byte(hashIV)

	// 加密
	ciphertext := make([]byte, len(urlEncodedData))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, []byte(urlEncodedData))

	// Base64 編碼
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func DecryptData(encryptData, hashKey, hashIV string) (string, error) {
	// 將 Base64 字符串解碼為原始加密數據
	encryptedData, err := base64.StdEncoding.DecodeString(encryptData)
	if err != nil {
		return "", err
	}

	// 根據給定的密鑰創建 cipher.Block
	block, err := aes.NewCipher([]byte(hashKey))
	if err != nil {
		return "", err
	}

	// 檢查加密數據的長度是否合法
	if len(encryptedData) < aes.BlockSize {
		return "", errors.New("ciphertext too short")
	}

	// 使用 CBC 模式解密數據
	mode := cipher.NewCBCDecrypter(block, []byte(hashIV))
	decryptedData := make([]byte, len(encryptedData))
	mode.CryptBlocks(decryptedData, encryptedData)

	// 移除 PKCS7 填充
	padding := int(decryptedData[len(decryptedData)-1])
	if padding < 1 || padding > aes.BlockSize || padding > len(decryptedData) {
		return "", errors.New("invalid padding")
	}
	decryptedData = decryptedData[:len(decryptedData)-padding]

	// URL 解碼
	result, err := url.QueryUnescape(string(decryptedData))
	if err != nil {
		return "", err
	}

	return result, nil
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

func SendFormData(c *client.ECPayClient, formData url.Values) ([]byte, error) {

	resp, err := http.PostForm(c.BaseURL, formData)
	if err != nil {
		slog.Error(fmt.Sprintf("Error sending POST request: %v", err))
		return nil, fmt.Errorf("error sending POST request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		if err = Body.Close(); err != nil {
			slog.Error(fmt.Sprintf("Error closing response body: %v", err))
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	return body, nil
}
