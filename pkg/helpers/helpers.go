package helpers

import (
	"crypto/sha256"
	"ecpay-go/pkg/client"
	"encoding/hex"
	"fmt"
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
			if field.Type() == reflect.TypeOf(time.Time{}) {
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
			slog.Error(fmt.Sprintf("Unsupported type: %v", field.Kind()))
		}
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

func SendFormData(c client.ECPayClient, formData url.Values) ([]byte, error) {

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
