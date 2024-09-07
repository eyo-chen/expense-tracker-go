package s3

import (
	"context"
	"errors"
	"testing"

	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/eyo-chen/expense-tracker-go/internal/adapter/interfaces"
	"github.com/eyo-chen/expense-tracker-go/mocks"
	"github.com/eyo-chen/expense-tracker-go/pkg/logger"
	"github.com/eyo-chen/expense-tracker-go/pkg/testutil"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

var (
	mockCTX        = context.Background()
	mockBucketName = "test-bucket"
)

type s3ServiceSuite struct {
	suite.Suite
	s3Service           interfaces.S3Service
	mockS3Client        *mocks.S3Client
	mockS3PresignClient *mocks.S3PresignClient
}

func TestS3ServiceSuite(t *testing.T) {
	suite.Run(t, new(s3ServiceSuite))
}

func (s *s3ServiceSuite) SetupSuite() {
	logger.Register()
}

func (s *s3ServiceSuite) SetupTest() {
	s.mockS3Client = new(mocks.S3Client)
	s.mockS3PresignClient = new(mocks.S3PresignClient)

	s.s3Service = New(mockBucketName, s.mockS3Client, s.mockS3PresignClient)
}

func (s *s3ServiceSuite) TearDownTest() {
	s.mockS3Client.AssertExpectations(s.T())
	s.mockS3PresignClient.AssertExpectations(s.T())
}

func (s *s3ServiceSuite) TestPutObjectUrl() {
	for scenario, fn := range map[string]func(s *s3ServiceSuite, desc string){
		"when no error, put object successfully":      putObject_NoError_ReturnSuccessfully,
		"when get presigned URL failed, return error": putObject_GetObjectUrlFailed_ReturnError,
	} {
		s.Run(testutil.GetFunName(fn), func() {
			s.SetupTest()
			fn(s, scenario)
			s.TearDownTest()
		})
	}
}

func putObject_NoError_ReturnSuccessfully(s *s3ServiceSuite, desc string) {
	// prepare mock data
	mockObjectKey := "test-object-key"
	mockLifetimeSecs := int64(3600)
	mockPutObjectInput := &s3.PutObjectInput{
		Bucket: &mockBucketName,
		Key:    &mockObjectKey,
	}
	mockPresignedURL := "https://test-bucket.s3.amazonaws.com/test-object-key"
	mockPresignedHTTPRequest := &v4.PresignedHTTPRequest{
		URL: mockPresignedURL,
	}
	mockFunc := mock.AnythingOfType("func(*s3.PresignOptions)")

	// prepare mock service
	s.mockS3PresignClient.On("PresignPutObject", mockCTX, mockPutObjectInput, mockFunc).
		Return(mockPresignedHTTPRequest, nil).Once()

	// action
	url, err := s.s3Service.PutObjectUrl(mockCTX, mockObjectKey, mockLifetimeSecs)

	// assertion
	s.Require().NoError(err, desc)
	s.Require().Equal(mockPresignedURL, url, desc)
}

func putObject_GetObjectUrlFailed_ReturnError(s *s3ServiceSuite, desc string) {
	// prepare mock data
	mockObjectKey := "test-object-key"
	mockLifetimeSecs := int64(3600)
	mockPutObjectInput := &s3.PutObjectInput{
		Bucket: &mockBucketName,
		Key:    &mockObjectKey,
	}
	mockFunc := mock.AnythingOfType("func(*s3.PresignOptions)")
	mockErr := errors.New("failed to get presigned URL")

	// prepare mock service
	s.mockS3PresignClient.On("PresignPutObject", mockCTX, mockPutObjectInput, mockFunc).
		Return(nil, mockErr).Once()

	// action
	url, err := s.s3Service.PutObjectUrl(mockCTX, mockObjectKey, mockLifetimeSecs)

	// assertion
	s.Require().Empty(url, desc)
	s.Require().ErrorIs(err, mockErr, desc)
}

func (s *s3ServiceSuite) TestGetObjectUrl() {
	for scenario, fn := range map[string]func(s *s3ServiceSuite, desc string){
		"when no error, get object URL successfully":  getObjectUrl_NoError_ReturnSuccessfully,
		"when get presigned URL failed, return error": getObjectUrl_GetPresignedUrlFailed_ReturnError,
	} {
		s.Run(testutil.GetFunName(fn), func() {
			s.SetupTest()
			fn(s, scenario)
			s.TearDownTest()
		})
	}
}

func getObjectUrl_NoError_ReturnSuccessfully(s *s3ServiceSuite, desc string) {
	// prepare mock data
	mockObjectKey := "test-object-key"
	mockLifetimeSecs := int64(3600)
	mockGetObjectUrlInput := &s3.GetObjectInput{
		Bucket: &mockBucketName,
		Key:    &mockObjectKey,
	}
	mockPresignedURL := "https://test-bucket.s3.amazonaws.com/test-object-key"
	mockPresignedHTTPRequest := &v4.PresignedHTTPRequest{
		URL: mockPresignedURL,
	}
	mockFunc := mock.AnythingOfType("func(*s3.PresignOptions)")

	// prepare mock service
	s.mockS3PresignClient.On("PresignGetObject", mockCTX, mockGetObjectUrlInput, mockFunc).
		Return(mockPresignedHTTPRequest, nil).Once()

	// action
	url, err := s.s3Service.GetObjectUrl(mockCTX, mockObjectKey, mockLifetimeSecs)

	// assertion
	s.Require().NoError(err, desc)
	s.Require().Equal(mockPresignedURL, url, desc)
}

func getObjectUrl_GetPresignedUrlFailed_ReturnError(s *s3ServiceSuite, desc string) {
	// prepare mock data
	mockObjectKey := "test-object-key"
	mockLifetimeSecs := int64(3600)
	mockGetObjectUrlInput := &s3.GetObjectInput{
		Bucket: &mockBucketName,
		Key:    &mockObjectKey,
	}
	mockFunc := mock.AnythingOfType("func(*s3.PresignOptions)")
	mockErr := errors.New("failed to get presigned URL")

	// prepare mock service
	s.mockS3PresignClient.On("PresignGetObject", mockCTX, mockGetObjectUrlInput, mockFunc).
		Return(nil, mockErr).Once()

	// action
	url, err := s.s3Service.GetObjectUrl(mockCTX, mockObjectKey, mockLifetimeSecs)

	// assertion
	s.Require().Empty(url, desc)
	s.Require().ErrorIs(err, mockErr, desc)
}

func (s *s3ServiceSuite) TestDeleteObject() {
	for scenario, fn := range map[string]func(s *s3ServiceSuite, desc string){
		"when no error, delete object successfully": deleteObject_NoError_DeleteSuccessfully,
		"when delete object failed, return error":   deleteObject_DeleteObjectFailed_ReturnError,
	} {
		s.Run(testutil.GetFunName(fn), func() {
			s.SetupTest()
			fn(s, scenario)
			s.TearDownTest()
		})
	}
}

func deleteObject_NoError_DeleteSuccessfully(s *s3ServiceSuite, desc string) {
	// prepare mock data
	mockObjectKey := "test-object-key"
	mockDeleteObjectInput := &s3.DeleteObjectInput{
		Bucket: &mockBucketName,
		Key:    &mockObjectKey,
	}

	// prepare mock service
	s.mockS3Client.On("DeleteObject", mockCTX, mockDeleteObjectInput).
		Return(&s3.DeleteObjectOutput{}, nil).Once()

	// action
	err := s.s3Service.DeleteObject(mockCTX, mockObjectKey)

	// assertion
	s.Require().NoError(err, desc)
}

func deleteObject_DeleteObjectFailed_ReturnError(s *s3ServiceSuite, desc string) {
	// prepare mock data
	mockObjectKey := "test-object-key"
	mockDeleteObjectInput := &s3.DeleteObjectInput{
		Bucket: &mockBucketName,
		Key:    &mockObjectKey,
	}
	mockErr := errors.New("failed to delete object")

	// prepare mock service
	s.mockS3Client.On("DeleteObject", mockCTX, mockDeleteObjectInput).
		Return(nil, mockErr).Once()

	// action
	err := s.s3Service.DeleteObject(mockCTX, mockObjectKey)

	// assertion
	s.Require().ErrorIs(err, mockErr, desc)
}
