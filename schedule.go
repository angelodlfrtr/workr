package workr

import (
	"context"
	"time"
)

// Schedule a job at defined time
func (wrkr *Workr) Schedule(at time.Time, jj Job) error {
	jobData, err := jj.Bytes()
	if err != nil {
		return err
	}

	queueItem := &QueueItem{
		JobName: jj.Name(),
		Data:    jobData,
	}

	member, err := queueItem.Bytes()
	if err != nil {
		return err
	}

	score := at.UnixNano()
	zMember := NewZ(float64(score), member)

	return wrkr.RedisClient.ZAdd(
		context.Background(),
		wrkr.Config.RedisQueueName,
		zMember,
	).Err()
}

// Now run job now
func (wrkr *Workr) Now(jj Job) error {
	return wrkr.Schedule(time.Now(), jj)
}
