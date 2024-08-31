package redisservice

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/eyo-chen/expense-tracker-go/internal/domain"
	"github.com/eyo-chen/expense-tracker-go/internal/model/interfaces"
	"github.com/eyo-chen/expense-tracker-go/pkg/dockerutil"
	"github.com/eyo-chen/expense-tracker-go/pkg/logger"
	"github.com/eyo-chen/expense-tracker-go/pkg/testutil"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
)

var (
	mockCTX = context.Background()
)

type redisServiceSuite struct {
	suite.Suite
	dk           *dockerutil.Container
	redis        *redis.Client
	redisService interfaces.RedisService
}

func TestRedisServiceSuitee(t *testing.T) {
	suite.Run(t, new(redisServiceSuite))
}

func (s *redisServiceSuite) SetupSuite() {
	dk := dockerutil.RunDocker(dockerutil.ImageRedis)
	s.dk = dk

	client := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("localhost:%s", dk.Port),
	})
	s.redis = client

	s.redisService = New(client)

	logger.Register()
}

func (s *redisServiceSuite) TearDownSuite() {
	s.redis.Close()
	s.dk.PurgeDocker()
}

func (s *redisServiceSuite) SetupTest() {
	s.redisService = New(s.redis)
}

func (s *redisServiceSuite) TearDownTest() {
	s.redis.FlushAll(mockCTX)
}

func (s *redisServiceSuite) TestGetByFunc() {
	for scenario, fn := range map[string]func(s *redisServiceSuite, desc string){
		"when cache miss, return value and cache val": getByFunc_CacheMiss_ReturnValueAndCacheVal,
		"when cache hit, return value":                getByFunc_CacheHit_ReturnValue,
		"when cache failed, return error":             getByFunc_CacheFailed_ReturnError,
	} {
		s.Run(testutil.GetFunName(fn), func() {
			s.SetupTest()
			fn(s, scenario)
			s.TearDownTest()
		})
	}
}

func getByFunc_CacheMiss_ReturnValueAndCacheVal(s *redisServiceSuite, desc string) {
	// prepare mock data
	mockKey := "key"
	mockValue := "value"
	mockTTL := time.Duration(0)
	mockFlag := false

	mockGetFun := func() (string, error) {
		mockFlag = true
		return mockValue, nil
	}

	// action
	value, err := s.redisService.GetByFunc(mockCTX, mockKey, mockTTL, mockGetFun)

	// assertion
	s.Require().NoError(err, desc)
	s.Require().Equal(mockValue, value, desc)
	s.Require().True(mockFlag, desc)

	// check cache
	cacheValue, err := s.redis.Get(mockCTX, mockKey).Result()
	s.Require().NoError(err)
	s.Require().Equal(mockValue, cacheValue)
}

func getByFunc_CacheHit_ReturnValue(s *redisServiceSuite, desc string) {
	// prepare mock data
	mockKey := "key"
	mockValue := "value"
	mockTTL := time.Duration(0)
	mockFlag := false

	mockGetFun := func() (string, error) {
		mockFlag = true
		return "", nil
	}

	// set value to cache
	err := s.redis.Set(mockCTX, mockKey, mockValue, 0).Err()
	s.Require().NoError(err)

	// action
	value, err := s.redisService.GetByFunc(mockCTX, mockKey, mockTTL, mockGetFun)

	// assertion
	s.Require().NoError(err, desc)
	s.Require().Equal(mockValue, value, desc)
	s.Require().False(mockFlag, desc)
}

func getByFunc_CacheFailed_ReturnError(s *redisServiceSuite, desc string) {
	// prepare mock data
	mockKey := "key"
	mockTTL := time.Duration(0)
	mockErr := fmt.Errorf("cache failed")

	mockGetFun := func() (string, error) {
		return "", mockErr
	}

	// test function
	value, err := s.redisService.GetByFunc(mockCTX, mockKey, mockTTL, mockGetFun)
	s.Require().ErrorIs(err, mockErr, desc)
	s.Require().Empty(value, desc)

	// check cache
	cacheValue, err := s.redis.Get(mockCTX, mockKey).Result()
	s.Require().ErrorIs(err, redis.Nil)
	s.Require().Empty(cacheValue)
}

func (s *redisServiceSuite) TestGetDel() {
	for scenario, fn := range map[string]func(s *redisServiceSuite, desc string){
		"when key exists, return value and delete from cache": testGetDel_KeyExists_ReturnValueAndDelete,
		"when key doesn't exist, return ErrCacheMiss":         testGetDel_KeyNotExists_ReturnErrCacheMiss,
	} {
		s.Run(testutil.GetFunName(fn), func() {
			s.SetupTest()
			fn(s, scenario)
			s.TearDownTest()
		})
	}
}

func testGetDel_KeyExists_ReturnValueAndDelete(s *redisServiceSuite, desc string) {
	// Prepare mock data
	mockKey := "test_key"
	mockValue := "test_value"

	// Set value in Redis
	err := s.redis.Set(mockCTX, mockKey, mockValue, 0).Err()
	s.Require().NoError(err, desc)

	// Action
	value, err := s.redisService.GetDel(mockCTX, mockKey)

	// Assertion
	s.Require().NoError(err, desc)
	s.Require().Equal(mockValue, value, desc)

	// Check if key was deleted
	_, err = s.redis.Get(mockCTX, mockKey).Result()
	s.Require().ErrorIs(err, redis.Nil, desc)
}

func testGetDel_KeyNotExists_ReturnErrCacheMiss(s *redisServiceSuite, desc string) {
	mockKey := "non_existent_key"

	// Action
	value, err := s.redisService.GetDel(mockCTX, mockKey)

	// Assertion
	s.Require().ErrorIs(err, domain.ErrCacheMiss, desc)
	s.Require().Empty(value, desc)
}

func (s *redisServiceSuite) TestSet() {
	for scenario, fn := range map[string]func(s *redisServiceSuite, desc string){
		"when set value, should store in redis": testSet_ValueStoredInRedis,
	} {
		s.Run(testutil.GetFunName(fn), func() {
			s.SetupTest()
			fn(s, scenario)
			s.TearDownTest()
		})
	}
}

func testSet_ValueStoredInRedis(s *redisServiceSuite, desc string) {
	// prepare mock data
	mockKey := "test_key"
	mockValue := "test_value"
	ttl := time.Duration(0)

	// action
	err := s.redisService.Set(mockCTX, mockKey, mockValue, ttl)

	// assertion
	s.Require().NoError(err, desc)

	// check cache
	value, err := s.redis.Get(mockCTX, mockKey).Result()
	s.Require().NoError(err, desc)
	s.Require().Equal(mockValue, value, desc)
}
