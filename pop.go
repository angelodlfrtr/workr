package workr

import (
	"context"
	"fmt"
	"time"

	redislib "github.com/go-redis/redis/v8"
)

var _popScript = redislib.NewScript(`
  local name = KEYS[1]
  local timestamp = KEYS[2]
  local limit = KEYS[3]
  local results = redis.call('zrangebyscore', name, '-inf', timestamp, 'LIMIT', 0, limit)
  if table.getn(results) > 0 then
    redis.call('zrem', name, unpack(results))
  end
  return results
`)

// Pop job from queue
func (wrkr *Workr) Pop() (*QueueItem, error) {
	resultsInt, err := _popScript.Run(
		context.Background(),
		wrkr.RedisClient,
		[]string{wrkr.Config.RedisQueueName, fmt.Sprintf("%d", time.Now().UnixNano()), "1"},
		3,
	).Result()
	if err != nil {
		return nil, err
	}

	if resultsInt == nil {
		return nil, nil
	}

	if results, ok := resultsInt.([]interface{}); ok {
		if len(results) > 0 {
			result := results[0]
			resultStr := result.(string)

			queueItem := &QueueItem{}
			if err := queueItem.Load([]byte(resultStr)); err != nil {
				return nil, err
			}

			return queueItem, nil
		}
	}

	return nil, nil
}
