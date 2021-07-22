// Package workr ...
package workr

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

var _tickInterval = 100 * time.Millisecond

// Workr define main worker instance
type Workr struct {
	Config      *Config
	RedisClient *redis.Client

	jobs      map[string]Job
	stopCH    chan bool
	stoppedCH chan bool
}

// New workr with config cfg
func New(cfg *Config) (*Workr, error) {
	cfg.SetDefaults()
	wrkr := &Workr{
		Config: cfg,
		jobs:   map[string]Job{},
	}

	// Test redis conn
	if err := wrkr.RedisConnect(); err != nil {
		wrkr.Config.Logger.Error("Error while connecting to redis", zap.Error(err))
		return nil, err
	}

	return wrkr, nil
}

// Work warch for jobs, blocking
func (wrkr *Workr) Work() {
	wrkr.stopCH = make(chan bool, 1)
	wrkr.stoppedCH = make(chan bool, 1)

	wrkr.Config.Logger.Debug("Staring worker loop")

	go wrkr.loop()
}

// Shutdown workr
// @TODO: handle running jobs end
func (wrkr *Workr) Shutdown(ctx context.Context) error {
	wrkr.Config.Logger.Debug("Shutdown requested")

	select {
	case wrkr.stopCH <- true:
	default:
	}

	select {
	case <-ctx.Done():
		wrkr.Config.Logger.Debug("Shutdown done due to ctx", zap.Error(ctx.Err()))
		return ctx.Err()
	case wrkr.stoppedCH <- true:
		wrkr.Config.Logger.Debug("Shutdown done")
		return nil
	}
}

// loop on redis queue
func (wrkr *Workr) loop() {
	for {
		select {
		case <-wrkr.stopCH:
			wrkr.stoppedCH <- true
			return
		default:
		}

		queueItem, err := wrkr.Pop()
		if err != nil {
			wrkr.Config.Logger.Error("Error while reading redis queue", zap.Error(err))
			time.Sleep(_tickInterval * 2)
			continue
		}

		if queueItem == nil {
			time.Sleep(_tickInterval)
			continue
		}

		wrkr.Config.Logger.Debug("Job received", zap.Any("queue_item", queueItem))

		if jd, ok := wrkr.jobs[queueItem.JobName]; ok {
			wrkr.Config.Logger.Debug("Job recognized", zap.Any("job", jd), zap.Any("queue_item", queueItem))

			jn := jd.New()
			if err := jn.Load(queueItem.Data); err != nil {
				wrkr.Config.Logger.Error("Error while login queue item", zap.Error(err))
				continue
			}

			wrkr.Config.Logger.Info("Job will run", zap.Any("job", jd), zap.Any("queue_item", queueItem))

			go func() {
				if err := jn.Run(wrkr); err != nil {
					wrkr.Config.Logger.Error(
						"Error while running job",
						zap.Error(err),
						zap.Any("job", jn),
						zap.Any("queue_item", queueItem),
					)
				}

				wrkr.Config.Logger.Info("Job done", zap.Any("job", jd), zap.Any("queue_item", queueItem))
			}()
		} else {
			wrkr.Config.Logger.Warn("Job not recognized", zap.Any("job", jd), zap.Any("queue_item", queueItem))
		}
	}
}
