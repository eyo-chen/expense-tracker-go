package codeutil

import (
	"encoding/base64"
	"errors"
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
)

// DecodeCursor decodes cursor from encoded string to map
func DecodeCursor(encodedString string) (map[string]string, error) {
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
		result[key] = value
	}

	return result, nil
}

// EncodeCursor encodes cursor from map to encoded string
func EncodeCursor(cursor map[string]string) string {
	pairs := make([]string, 0, len(cursor))
	for key, value := range cursor {
		pairs = append(pairs, key+":"+value)
	}

	encodedString := base64.StdEncoding.EncodeToString([]byte(strings.Join(pairs, ",")))
	return encodedString
}
