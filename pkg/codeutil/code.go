package codeutil

import (
	"encoding/base64"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/OYE0303/expense-tracker-go/internal/domain"
	"github.com/OYE0303/expense-tracker-go/pkg/logger"
)

var (
	// ErrEmptyEncodedString is an error for empty encoded string
	ErrEmptyEncodedString = errors.New("empty encoded string")

	// ErrInvalidCursor is an error for invalid cursor
	ErrInvalidCursor = errors.New("invalid cursor")

	// ErrInvalidFormatCursor is an error for invalid format cursor
	ErrInvalidFormatCursor = errors.New("invalid format cursor")

	// ErrFieldNotFound is an error for field not found
	ErrFieldNotFound = errors.New("field not found")
)

// DecodeNextKeys decodes encoded next key string to decoded next keys.
// fieldSource is used to check if the field exists in the fieldSource
func DecodeNextKeys(encodedString string, fieldSource interface{}) (domain.DecodedNextKeys, error) {
	if encodedString == "" {
		return nil, ErrEmptyEncodedString
	}

	decodedBytes, err := base64.StdEncoding.DecodeString(encodedString)
	if err != nil {
		logger.Error("Decode Cursor failed", "err", err)
		return nil, ErrInvalidCursor
	}

	decodedString := string(decodedBytes)
	pairs := strings.Split(decodedString, ",")
	if len(pairs) == 0 {
		return nil, ErrInvalidFormatCursor
	}

	result := domain.DecodedNextKeys{}
	for _, pair := range pairs {
		keyValue := strings.Split(pair, ":")
		if len(keyValue) != 2 {
			return nil, ErrInvalidFormatCursor
		}

		key := strings.TrimSpace(keyValue[0])
		value := strings.TrimSpace(keyValue[1])

		// check if the field exists in the fieldSource
		if fieldSource != nil {
			if _, ok := getFieldValue(fieldSource, key); !ok {
				return nil, ErrFieldNotFound
			}
		}

		result = append(result, domain.DecodedNextKeyInfo{
			Field: key,
			Value: value,
		})
	}

	return result, nil
}

// EncodeNextKeys encodes decoded next keys to encoded string.
// fieldSource is used to get the latest field value from the source
func EncodeNextKeys(decodedNextKeys domain.DecodedNextKeys, fieldSource interface{}) (string, error) {
	var pairs string

	for _, key := range decodedNextKeys {
		if fieldSource == nil {
			pairs += key.Field + ":" + key.Value + ","
			continue
		}

		// note that we have to use the fieldSource to get the value
		// and set it to encoded string when encoding
		v, ok := getFieldValue(fieldSource, key.Field)
		if !ok {
			return "", ErrFieldNotFound
		}

		pairs += key.Field + ":" + cvtToString(v) + ","
	}

	// remove the last comma
	pairs = pairs[:len(pairs)-1]

	encodedString := base64.StdEncoding.EncodeToString([]byte(pairs))
	return encodedString, nil
}

func getFieldValue(val interface{}, fieldName string) (interface{}, bool) {
	v := reflect.ValueOf(val)
	field := v.FieldByName(fieldName)
	if !field.IsValid() || !field.CanInterface() {
		return nil, false
	}

	return field.Interface(), true
}

// cvtToString converts any value to string
// e.g. 1 -> "1", 1.1 -> "1.1", true -> "true"
func cvtToString(v interface{}) string {
	switch val := v.(type) {
	case string:
		return val
	case int, int8, int16, int32, int64:
		return fmt.Sprintf("%d", val)
	case uint, uint8, uint16, uint32:
		return fmt.Sprintf("%d", val)
	case float32, float64:
		return fmt.Sprintf("%f", val)
	case bool:
		return fmt.Sprintf("%t", val)
	case time.Time:
		return val.Format(time.RFC3339)
	default:
		return ""
	}
}
