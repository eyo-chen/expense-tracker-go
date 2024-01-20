package testutil

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type utils struct {
	suite.Suite
}

func TestUtils(t *testing.T) {
	suite.Run(t, new(utils))
}

func functionName()  {}
func function_name() {}
func FunctionName()  {}
func Function_Name() {}
func FUNCTIONNAME()  {}
func functionname()  {}
func f()             {}

func (s *utils) TestGetFunName() {
	tests := []struct {
		Desc  string
		Input func()
		Exp   string
	}{
		{
			Desc:  "camelCase",
			Input: functionName,
			Exp:   "functionName",
		},
		{
			Desc:  "snake_case",
			Input: function_name,
			Exp:   "function_name",
		},
		{
			Desc:  "CamelCase",
			Input: FunctionName,
			Exp:   "FunctionName",
		},
		{
			Desc:  "Snake_Case",
			Input: Function_Name,
			Exp:   "Function_Name",
		},
		{
			Desc:  "UPPERCASE",
			Input: FUNCTIONNAME,
			Exp:   "FUNCTIONNAME",
		},
		{
			Desc:  "lowercase",
			Input: functionname,
			Exp:   "functionname",
		},
		{
			Desc:  "single letter",
			Input: f,
			Exp:   "f",
		},
		{
			Desc:  "nil",
			Input: nil,
			Exp:   "",
		},
	}

	for _, test := range tests {
		s.Run(test.Desc, func() {
			s.Require().Equal(test.Exp, GetFunName(test.Input))
		})
	}
}
