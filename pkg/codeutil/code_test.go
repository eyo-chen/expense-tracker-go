package codeutil_test

import (
	"encoding/base64"
	"testing"

	"github.com/eyo-chen/expense-tracker-go/internal/domain"
	"github.com/eyo-chen/expense-tracker-go/pkg/codeutil"
	"github.com/eyo-chen/expense-tracker-go/pkg/logger"
	"github.com/eyo-chen/expense-tracker-go/pkg/testutil"
	"github.com/stretchr/testify/suite"
)

type CodeUtilSuite struct {
	suite.Suite
}

func TestEncodeSuite(t *testing.T) {
	suite.Run(t, new(CodeUtilSuite))
}

func (s *CodeUtilSuite) SetupSuite() {
	logger.Register()
}

func (s *CodeUtilSuite) TestDecodeNextKeys() {
	for scenario, fn := range map[string]func(*CodeUtilSuite, string){
		"when the encoded string is empty, return an error":   decodeNextKeys_EncodedEmptyString_ReturnErr,
		"when the decoded string is empty, return an error":   decodeNextKeys_DecodedEmptyString_ReturnErr,
		"when the cursor format is invalid, return an error":  decodeNextKeys_InvalidFormatCursor_ReturnErr,
		"when the source field is not found, return an error": decodeNextKeys_SourceFieldNotFound_ReturnErr,
		"when the encoded string is valid, return cursor map": decodeNextKeys_ValidEncodedString_ReturnCursorMap,
		"when the source field is correct, return cursor map": decodeNextKeys_WithCorrectSourceField_ReturnCursorMap,
	} {
		s.Run(testutil.GetFunName(fn), func() {
			fn(s, scenario)
		})
	}
}

func decodeNextKeys_EncodedEmptyString_ReturnErr(s *CodeUtilSuite, desc string) {
	// prepare encoded string
	encodedString := ""

	// action
	result, err := codeutil.DecodeNextKeys(encodedString, nil)
	s.Require().Nil(result, desc)
	s.Require().Equal(codeutil.ErrEmptyEncodedString, err, desc)
}

func decodeNextKeys_DecodedEmptyString_ReturnErr(s *CodeUtilSuite, desc string) {
	// prepare encoded string
	cursorKey := ""
	encodedString := base64.StdEncoding.EncodeToString([]byte(cursorKey))

	// action
	result, err := codeutil.DecodeNextKeys(encodedString, nil)
	s.Require().Nil(result, desc)
	s.Require().Equal(codeutil.ErrEmptyEncodedString, err, desc)
}

func decodeNextKeys_InvalidFormatCursor_ReturnErr(s *CodeUtilSuite, desc string) {
	// prepare encoded string
	cursorKey := "ID:123,MainCategID"
	encodedString := base64.StdEncoding.EncodeToString([]byte(cursorKey))

	// action
	result, err := codeutil.DecodeNextKeys(encodedString, nil)
	s.Require().Nil(result, desc)
	s.Require().Equal(codeutil.ErrInvalidFormatCursor, err, desc)
}

func decodeNextKeys_SourceFieldNotFound_ReturnErr(s *CodeUtilSuite, desc string) {
	// prepare encoded string
	cursorKey := "ID:123,MainCategID:456"
	encodedString := base64.StdEncoding.EncodeToString([]byte(cursorKey))

	// prepare field source
	fieldSource := struct {
		ID string
	}{}

	// action
	result, err := codeutil.DecodeNextKeys(encodedString, fieldSource)
	s.Require().Nil(result, desc)
	s.Require().Equal(codeutil.ErrFieldNotFound, err, desc)
}

func decodeNextKeys_ValidEncodedString_ReturnCursorMap(s *CodeUtilSuite, desc string) {
	// prepare encoded string
	cursorKey := "MainCategID:456,ID:123"
	encodedString := base64.StdEncoding.EncodeToString([]byte(cursorKey))

	// prepare expected result
	cursorMap := domain.DecodedNextKeys{
		{Field: "MainCategID", Value: "456"},
		{Field: "ID", Value: "123"},
	}

	// action
	result, err := codeutil.DecodeNextKeys(encodedString, nil)
	s.Require().NoError(err, desc)
	s.Require().Equal(cursorMap, result, desc)
}

func decodeNextKeys_WithCorrectSourceField_ReturnCursorMap(s *CodeUtilSuite, desc string) {
	// prepare encoded string
	cursorKey := "MainCategID:456,ID:123"
	encodedString := base64.StdEncoding.EncodeToString([]byte(cursorKey))

	// prepare field source
	fieldSource := struct {
		ID          string
		MainCategID string
	}{}

	// prepare expected result
	cursorMap := domain.DecodedNextKeys{
		{Field: "MainCategID", Value: "456"},
		{Field: "ID", Value: "123"},
	}

	// action
	result, err := codeutil.DecodeNextKeys(encodedString, fieldSource)
	s.Require().NoError(err, desc)
	s.Require().Equal(cursorMap, result, desc)
}

func (s *CodeUtilSuite) TestEncodeNextKeys() {
	for scenario, fn := range map[string]func(*CodeUtilSuite, string){
		"when the field is not found, return an error":            encodeNextKeys_FieldNotFound_ReturnErr,
		"when the cursor map is valid, return encoded string":     encodeNextKeys_ValidCursorMap_ReturnEncodedString,
		"when the field source is correct, return encoded string": encodeNextKeys_WithCorrectFieldSource_ReturnEncodedString,
	} {
		s.Run(testutil.GetFunName(fn), func() {
			fn(s, scenario)
		})
	}

}

func encodeNextKeys_FieldNotFound_ReturnErr(s *CodeUtilSuite, desc string) {
	// prepare next keys
	nextKeys := domain.DecodedNextKeys{
		{Field: "ID", Value: "123"},
		{Field: "MainCategID", Value: "456"},
	}

	// prepare field source
	fieldSource := struct {
		ID string
	}{}

	// action
	result, err := codeutil.EncodeNextKeys(nextKeys, fieldSource)
	s.Require().Empty(result, desc)
	s.Require().Equal(codeutil.ErrFieldNotFound, err, desc)
}

func encodeNextKeys_ValidCursorMap_ReturnEncodedString(s *CodeUtilSuite, desc string) {
	// prepare next keys
	nextKeys := domain.DecodedNextKeys{
		{Field: "ID", Value: "123"},
		{Field: "MainCategID", Value: "456"},
	}

	// prepare expected result
	expResult := "ID:123,MainCategID:456"

	// action
	result, err := codeutil.EncodeNextKeys(nextKeys, nil)
	s.Require().NoError(err, desc)

	// check decoded string
	decodedBytes, err := base64.StdEncoding.DecodeString(result)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, string(decodedBytes), desc)
}

func encodeNextKeys_WithCorrectFieldSource_ReturnEncodedString(s *CodeUtilSuite, desc string) {
	// prepare next keys
	nextKeys := domain.DecodedNextKeys{
		{Field: "MainCategID", Value: "456"},
		{Field: "ID", Value: "123"},
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
	expResult := "MainCategID:456new,ID:123"

	// action
	result, err := codeutil.EncodeNextKeys(nextKeys, fieldSource)
	s.Require().NoError(err, desc)

	// check encoded string
	decodedBytes, err := base64.StdEncoding.DecodeString(result)
	s.Require().NoError(err, desc)
	s.Require().Equal(expResult, string(decodedBytes), desc)
}
