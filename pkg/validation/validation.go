package validation

import (
	"errors"
	"github.com/EcomPlatformOrg/ecpay-go/pkg/helpers"
	"net/url"
)

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
	expectedCheckMacValue := helpers.GenerateCheckMacValue(responseValues, hashKey, hashIV)

	if expectedCheckMacValue != receivedCheckMacValue {
		return errors.New("CheckMacValue mismatch")
	}

	return nil
}
