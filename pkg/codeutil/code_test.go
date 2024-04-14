package codeutil_test

import (
	"encoding/base64"
	"fmt"
	"testing"

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
		"when the encoded string is valid, return cursor map": decodeCursor_ValidEncodedString_ReturnCursorMap,
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
	result, err := codeutil.DecodeCursor(encodedString)
	s.Require().Nil(result, desc)
	s.Require().Equal(codeutil.ErrEmptyEncodedString, err, desc)
}

func decodeCursor_DecodedEmptyString_ReturnErr(s *CodeUtilSuite, desc string) {
	// prepare encoded string
	cursorKey := ""
	encodedString := base64.StdEncoding.EncodeToString([]byte(cursorKey))

	fmt.Println("encodedString: ", encodedString)

	// action
	result, err := codeutil.DecodeCursor(encodedString)
	s.Require().Nil(result, desc)
	s.Require().Equal(codeutil.ErrEmptyEncodedString, err, desc)
}

func decodeCursor_InvalidFormatCursor_ReturnErr(s *CodeUtilSuite, desc string) {
	// prepare encoded string
	cursorKey := "id:123,main_category_id"
	encodedString := base64.StdEncoding.EncodeToString([]byte(cursorKey))

	// action
	result, err := codeutil.DecodeCursor(encodedString)
	s.Require().Nil(result, desc)
	s.Require().Equal(codeutil.ErrInvalidFormatCursor, err, desc)
}

func decodeCursor_ValidEncodedString_ReturnCursorMap(s *CodeUtilSuite, desc string) {
	// prepare encoded string
	cursorKey := "id:123,main_category_id:456"
	encodedString := base64.StdEncoding.EncodeToString([]byte(cursorKey))

	// prepare expected result
	cursorMap := map[string]string{
		"id":               "123",
		"main_category_id": "456",
	}

	// action
	result, err := codeutil.DecodeCursor(encodedString)
	s.Require().NoError(err)
	s.Require().Equal(cursorMap, result)
}

func (s *CodeUtilSuite) TestEncodeCursor() {
	cursor := map[string]string{
		"id": "123",
	}

	encodedString := codeutil.EncodeCursor(cursor)

	// check
	encodedBytes, err := base64.StdEncoding.DecodeString(encodedString)
	s.Require().NoError(err)
	s.Require().Equal("id:123", string(encodedBytes))
}
