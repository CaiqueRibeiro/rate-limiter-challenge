package strategies

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/assert"
)

func mockNow() time.Time {
	return time.Date(2024, 10, 24, 3, 0, 0, 0, time.Local)
}

func TestRedisLimiterStrategy(t *testing.T) {
	db, clientMock := redismock.NewClientMock()
	ipMaxReqs := 5
	timeWindow := int64(1000)
	token := "dummy_token"

	strategy := NewRedisLimiter(db, mockNow)

	t.Run("Should allow when key is informed for first time", func(t *testing.T) {
		expectedTTL := time.Duration(timeWindow) * time.Millisecond
		key := fmt.Sprintf("limit:%s", token)

		clientMock.ExpectGet(key).RedisNil()
		clientMock.ExpectIncr(key).SetVal(1)
		clientMock.ExpectTTL(key).SetVal(time.Duration(-1))
		clientMock.ExpectExpire(key, time.Duration(timeWindow)).RedisNil()

		request := &Request{
			Key:      token,
			Limit:    int64(ipMaxReqs),
			Duration: time.Duration(timeWindow) * time.Millisecond,
		}

		result, err := strategy.CheckLimit(context.Background(), request)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, Allow, result.Result)
		assert.Equal(t, int64(ipMaxReqs), result.Limit)
		assert.Equal(t, int64(1), result.Total)
		assert.Equal(t, int64(ipMaxReqs)-1, result.Remaining)
		assert.WithinDuration(t, mockNow().Add(expectedTTL), result.ExpiresAt, time.Second)

		clientMock.ClearExpect()
	})

	t.Run("Should allow key exists and limit is not reached yet", func(t *testing.T) {
		expectedTTL := time.Duration(timeWindow) * time.Millisecond
		key := fmt.Sprintf("limit:%s", token)

		clientMock.ExpectGet(key).SetVal("1")
		clientMock.ExpectTTL(key).SetVal(expectedTTL)
		clientMock.ExpectIncr(key).SetVal(2)

		request := &Request{
			Key:      token,
			Limit:    int64(ipMaxReqs),
			Duration: time.Duration(timeWindow) * time.Millisecond,
		}

		result, err := strategy.CheckLimit(context.Background(), request)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, Allow, result.Result)
		assert.Equal(t, int64(ipMaxReqs), result.Limit)
		assert.Equal(t, int64(2), result.Total)
		assert.Equal(t, int64(ipMaxReqs)-2, result.Remaining)
		assert.WithinDuration(t, mockNow().Add(expectedTTL), result.ExpiresAt, time.Second)

		clientMock.ClearExpect()
	})

	t.Run("Should allow key exists and limit is reached", func(t *testing.T) {
		expectedTTL := time.Duration(timeWindow) * time.Millisecond
		key := fmt.Sprintf("limit:%s", token)

		clientMock.ExpectGet(key).SetVal(fmt.Sprint(ipMaxReqs))
		clientMock.ExpectTTL(key).SetVal(expectedTTL)
		clientMock.ExpectIncr(key).SetVal(2)

		request := &Request{
			Key:      token,
			Limit:    int64(ipMaxReqs),
			Duration: time.Duration(timeWindow) * time.Millisecond,
		}

		result, err := strategy.CheckLimit(context.Background(), request)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, Deny, result.Result)
		assert.Equal(t, int64(ipMaxReqs), result.Limit)
		assert.Equal(t, int64(ipMaxReqs), result.Total)
		assert.Equal(t, int64(0), result.Remaining)
		assert.WithinDuration(t, mockNow().Add(expectedTTL), result.ExpiresAt, time.Second)

		clientMock.ClearExpect()
	})
}
