package codeutil_test

import (
	"encoding/base64"
	"strings"
	"testing"

	"github.com/OYE0303/expense-tracker-go/internal/domain"
	"github.com/OYE0303/expense-tracker-go/pkg/codeutil"
	"github.com/OYE0303/expense-tracker-go/pkg/testutil"
	"github.com/stretchr/testify/suite"
)

type CodeUtilSuite struct {
	suite.Suite
}

func TestEncodeSuite(t *testing.T) {
	suite.Run(t, new(CodeUtilSuite))
}

func (s *CodeUtilSuite) TestDecodeCursor() {
	for scenario, fn := range map[string]func(*CodeUtilSuite, string){
		"when the encoded string is empty, return an error":   decodeCursor_EncodedEmptyString_ReturnErr,
		"when the decoded string is empty, return an error":   decodeCursor_DecodedEmptyString_ReturnErr,
		"when the cursor format is invalid, return an error":  decodeCursor_InvalidFormatCursor_ReturnErr,
		"when the source field is not found, return an error": decodeCursor_SourceFieldNotFound_ReturnErr,
		"when the encoded string is valid, return cursor map": decodeCursor_ValidEncodedString_ReturnCursorMap,
		"when the source field is correct, return cursor map": decodeCursor_WithCorrectSourceField_ReturnCursorMap,
	} {
		s.Run(testutil.GetFunName(fn), func() {
			fn(s, scenario)
		})
	}
}

func decodeCursor_EncodedEmptyString_ReturnErr(s *CodeUtilSuite, desc string) {
	// prepare encoded string
	encodedString := ""

	// action
	result, err := codeutil.DecodeCursor(encodedString, nil)
	s.Require().Nil(result, desc)
	s.Require().Equal(codeutil.ErrEmptyEncodedString, err, desc)
}

func decodeCursor_DecodedEmptyString_ReturnErr(s *CodeUtilSuite, desc string) {
	// prepare encoded string
	cursorKey := ""
	encodedString := base64.StdEncoding.EncodeToString([]byte(cursorKey))

	// action
	result, err := codeutil.DecodeCursor(encodedString, nil)
	s.Require().Nil(result, desc)
	s.Require().Equal(codeutil.ErrEmptyEncodedString, err, desc)
}

func decodeCursor_InvalidFormatCursor_ReturnErr(s *CodeUtilSuite, desc string) {
	// prepare encoded string
	cursorKey := "ID:123,MainCategID"
	encodedString := base64.StdEncoding.EncodeToString([]byte(cursorKey))

	// action
	result, err := codeutil.DecodeCursor(encodedString, nil)
	s.Require().Nil(result, desc)
	s.Require().Equal(codeutil.ErrInvalidFormatCursor, err, desc)
}

func decodeCursor_SourceFieldNotFound_ReturnErr(s *CodeUtilSuite, desc string) {
	// prepare encoded string
	cursorKey := "ID:123,MainCategID:456"
	encodedString := base64.StdEncoding.EncodeToString([]byte(cursorKey))

	// prepare field source
	fieldSource := struct {
		ID string
	}{}

	// action
	result, err := codeutil.DecodeCursor(encodedString, fieldSource)
	s.Require().Nil(result, desc)
	s.Require().Equal(codeutil.ErrFieldNotFound, err, desc)
}

func decodeCursor_ValidEncodedString_ReturnCursorMap(s *CodeUtilSuite, desc string) {
	// prepare encoded string
	cursorKey := "ID:123,MainCategID:456"
	encodedString := base64.StdEncoding.EncodeToString([]byte(cursorKey))

	// prepare expected result
	cursorMap := domain.DecodedNextKey{
		"ID":          "123",
		"MainCategID": "456",
	}

	// action
	result, err := codeutil.DecodeCursor(encodedString, nil)
	s.Require().NoError(err, desc)
	s.Require().Equal(cursorMap, result, desc)
}

func decodeCursor_WithCorrectSourceField_ReturnCursorMap(s *CodeUtilSuite, desc string) {
	// prepare encoded string
	cursorKey := "ID:123,MainCategID:456"
	encodedString := base64.StdEncoding.EncodeToString([]byte(cursorKey))

	// prepare field source
	fieldSource := struct {
		ID          string
		MainCategID string
	}{}

	// prepare expected result
	cursorMap := domain.DecodedNextKey{
		"ID":          "123",
		"MainCategID": "456",
	}

	// action
	result, err := codeutil.DecodeCursor(encodedString, fieldSource)
	s.Require().NoError(err, desc)
	s.Require().Equal(cursorMap, result, desc)
}

func (s *CodeUtilSuite) TestEncodeCursor() {
	for scenario, fn := range map[string]func(*CodeUtilSuite, string){
		"when the field is not found, return an error":            encodeCursor_FieldNotFound_ReturnErr,
		"when the cursor map is valid, return encoded string":     encodeCursor_ValidCursorMap_ReturnEncodedString,
		"when the field source is correct, return encoded string": encodeCursor_WithCorrectFieldSource_ReturnEncodedString,
	} {
		s.Run(testutil.GetFunName(fn), func() {
			fn(s, scenario)
		})
	}

}

func encodeCursor_FieldNotFound_ReturnErr(s *CodeUtilSuite, desc string) {
	// prepare cursor map
	cursorMap := domain.DecodedNextKey{
		"ID":          "123",
		"MainCategID": "456",
	}

	// prepare field source
	fieldSource := struct {
		ID string
	}{}

	// action
	result, err := codeutil.EncodeCursor(cursorMap, fieldSource)
	s.Require().Empty(result, desc)
	s.Require().Equal(codeutil.ErrFieldNotFound, err, desc)
}

// In the EncodeCursor function, the output is Base64 encoded string, and it's random
// For example, the same input "ID:123,MainCategID:456" can be encoded to "SUQ6MTIzLE1haW5DYWdlZElEOjQ1Ng==" or "SUQ6MTIzLE1haW5DYWdlZElEOjewdwe==
// So, we can't check the exact value of the encoded string
// Also, the output of the decoded string is random too
// For example, the encoded string can be decoded to "ID:123,MainCategID:456" or "MainCategID:456,ID:123"
// The only way we can check is to check the number of pairs and the value of the pairs respectively (using for loop)
func encodeCursor_ValidCursorMap_ReturnEncodedString(s *CodeUtilSuite, desc string) {
	// prepare cursor map
	cursorMap := domain.DecodedNextKey{
		"ID":          "123",
		"MainCategID": "456",
	}

	// prepare expected result
	expectedCursorMap := domain.DecodedNextKey{
		"ID":          "123",
		"MainCategID": "456",
	}

	// action
	result, err := codeutil.EncodeCursor(cursorMap, nil)
	s.Require().NoError(err, desc)

	// check decoded string
	decodedBytes, err := base64.StdEncoding.DecodeString(result)
	s.Require().NoError(err, desc)

	// using for loop to check the value of the pairs
	decodedString := string(decodedBytes)
	pairs := strings.Split(decodedString, ",")
	s.Require().Equal(2, len(pairs), desc)
	for _, pair := range pairs {
		keyValue := strings.Split(pair, ":")
		s.Require().Equal(2, len(keyValue), desc)

		key := strings.TrimSpace(keyValue[0])
		value := strings.TrimSpace(keyValue[1])
		v, ok := expectedCursorMap[key]
		s.Require().True(ok, desc)
		s.Require().Equal(v, value, desc)
	}
}

func encodeCursor_WithCorrectFieldSource_ReturnEncodedString(s *CodeUtilSuite, desc string) {
	// prepare cursor map
	cursorMap := domain.DecodedNextKey{
		"ID":          "123",
		"MainCategID": "456",
	}

	// prepare field source
	fieldSource := struct {
		ID          string
		MainCategID string
	}{
		ID:          "123",
		MainCategID: "456new",
	}

	// prepare expected result
	expectedCursorMap := domain.DecodedNextKey{
		"ID":          "123",
		"MainCategID": "456new",
	}

	// action
	result, err := codeutil.EncodeCursor(cursorMap, fieldSource)
	s.Require().NoError(err, desc)

	// check encoded string
	decodedBytes, err := base64.StdEncoding.DecodeString(result)
	s.Require().NoError(err, desc)

	// using for loop to check the value of the pairs
	decodedString := string(decodedBytes)
	pairs := strings.Split(decodedString, ",")
	s.Require().Equal(2, len(pairs), desc)
	for _, pair := range pairs {
		keyValue := strings.Split(pair, ":")
		s.Require().Equal(2, len(keyValue), desc)

		key := strings.TrimSpace(keyValue[0])
		value := strings.TrimSpace(keyValue[1])
		v, ok := expectedCursorMap[key]
		s.Require().True(ok, desc)
		s.Require().Equal(v, value, desc)
	}
}
