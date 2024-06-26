package logistics

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/EcomPlatformOrg/ecpay-go/pkg/helpers"
	"github.com/EcomPlatformOrg/ecpay-go/pkg/model"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// CreateTestData is a method that creates test data for the ECPayLogistics struct and sends it to the ECPayClient server for processing and decryption. The method returns the decrypted
func (e *ECPayLogistics) CreateTestData() (*ECPayLogistics, error) {

	if err := e.EncryptLogistics(); err != nil {
		return nil, err
	}

	e.RqHeader = model.RqHeader{
		Timestamp: strconv.FormatInt(time.Now().Unix(), 10),
	}

	jsonData, err := json.Marshal(e)
	if err != nil {
		slog.Error(fmt.Sprintf("Error marshalling ECPayLogistics struct: %v", err))
		return nil, err
	}

	resp, err := http.Post(e.Client.BaseURL, "application/json", strings.NewReader(string(jsonData)))
	if err != nil {
		slog.Error(fmt.Sprintf("Error sending POST request: %v", err))
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		if err = Body.Close(); err != nil {
			slog.Error(err.Error())
		}
	}(resp.Body)

	responseData := ECPayLogistics{}
	if err = json.NewDecoder(resp.Body).Decode(&responseData); err != nil {
		slog.Error(fmt.Sprintf("Error decoding response body: %v", err))
		return nil, err
	}

	slog.Info(fmt.Sprintf("responseData.TransCode %d", responseData.TransCode))
	slog.Info(fmt.Sprintf("responseData.TransMsg %d", responseData.TransMsg))
	decryptedData := &ECPayLogistics{}
	decryptedDataString, err := helpers.DecryptData(responseData.Data, e.Client.HashKey, e.Client.HashIV)
	if err = json.NewDecoder(bytes.NewReader([]byte(decryptedDataString))).Decode(&decryptedData); err != nil {
		slog.Error(fmt.Sprintf("Error decoding decrypted data: %v", err))
		return nil, err
	}

	return decryptedData, nil
}
