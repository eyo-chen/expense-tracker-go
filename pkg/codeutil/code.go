package codeutil

import (
	"encoding/base64"
	"errors"
	"reflect"
	"strings"

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

// DecodeCursor decodes cursor from encoded string to map
// fieldSource is used to check if the field exists in ecoded string
func DecodeCursor(encodedString string, fieldSource interface{}) (map[string]string, error) {
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

	result := map[string]string{}
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

		result[key] = value
	}

	return result, nil
}

// EncodeCursor encodes cursor from map to encoded string
// fieldSource is used to get the field value from the source
func EncodeCursor(cursor map[string]string, fieldSource interface{}) (string, error) {
	pairs := make([]string, 0, len(cursor))
	for key, value := range cursor {
		if fieldSource == nil {
			pairs = append(pairs, key+":"+value)
			continue
		}

		// note that we have to use the fieldSource to get the value
		// and set it to encoded string when encoding
		v, ok := getFieldValue(fieldSource, key)
		if !ok {
			return "", ErrFieldNotFound
		}
		pairs = append(pairs, key+":"+v.(string))
	}

	encodedString := base64.StdEncoding.EncodeToString([]byte(strings.Join(pairs, ",")))
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
