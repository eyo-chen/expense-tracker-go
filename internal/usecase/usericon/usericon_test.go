package usericon

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/eyo-chen/expense-tracker-go/mocks"
	"github.com/eyo-chen/expense-tracker-go/pkg/testutil"
	"github.com/stretchr/testify/suite"
)

var (
	mockCTX = context.Background()
)

type UserIconSuite struct {
	suite.Suite
	uc            *UC
	mockS3Service *mocks.S3Service
}

func TestUserIconSuite(t *testing.T) {
	suite.Run(t, new(UserIconSuite))
}

func (s *UserIconSuite) SetupTest() {
	s.mockS3Service = mocks.NewS3Service(s.T())
	s.uc = New(s.mockS3Service)
}

func (s *UserIconSuite) TearDownTest() {
	s.mockS3Service.AssertExpectations(s.T())
}

func (s *UserIconSuite) TestGetPutObjectURL() {
	for scenario, fn := range map[string]func(s *UserIconSuite, desc string){
		"when no error, return success":            getPutObjectURL_NoError_ReturnSuccessfully,
		"when get object url failed, return error": getPutObjectURL_GetObjectUrlFailed_ReturnError,
	} {
		s.Run(testutil.GetFunName(fn), func() {
			s.SetupTest()
			fn(s, scenario)
			s.TearDownTest()
		})
	}
}

func getPutObjectURL_NoError_ReturnSuccessfully(s *UserIconSuite, desc string) {
	// prepare mock data
	mockUserID := int64(1)
	mockFileName := "test.png"
	mockObjectKey := fmt.Sprintf("user_icons/%d/%s", mockUserID, mockFileName)
	mockTTL := 60 * time.Second
	mockURL := "http://s3.com/test.png"

	// prepare mock service
	s.mockS3Service.On("PutObjectUrl", mockCTX, mockObjectKey, int64(mockTTL.Seconds())).Return(mockURL, nil)

	// test function
	url, err := s.uc.GetPutObjectURL(mockCTX, mockFileName, mockUserID)
	s.Require().NoError(err, desc)
	s.Require().Equal(mockURL, url, desc)
}

func getPutObjectURL_GetObjectUrlFailed_ReturnError(s *UserIconSuite, desc string) {
	// prepare mock data
	mockUserID := int64(1)
	mockFileName := "test.png"
	mockObjectKey := fmt.Sprintf("user_icons/%d/%s", mockUserID, mockFileName)
	mockTTL := 60 * time.Second
	mockErr := errors.New("get object url failed")

	// prepare mock service
	s.mockS3Service.On("PutObjectUrl", mockCTX, mockObjectKey, int64(mockTTL.Seconds())).Return("", mockErr)

	// test function
	url, err := s.uc.GetPutObjectURL(mockCTX, mockFileName, mockUserID)
	s.Require().ErrorIs(err, mockErr, desc)
	s.Require().Empty(url, desc)
}
