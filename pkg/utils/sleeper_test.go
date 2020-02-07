package utils_test

import (
	"context"
	"io"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/ysmood/kit/pkg/utils"
)

func TestSleep(t *T) {
	utils.Sleep(0.01)
}

func TestPause(t *T) {
	go utils.Pause()
}

func TestBackoffSleeperWakeNow(t *T) {
	utils.E(utils.BackoffSleeper(0, 0, nil)(context.Background()))
}

func TestRetry(t *T) {
	count := 0
	s1 := utils.BackoffSleeper(1, 5, nil)
	s2 := utils.BackoffSleeper(2, 5, nil)
	s := utils.MergeSleepers(s1, s2)

	err := utils.Retry(context.Background(), s, func() (bool, error) {
		if count > 5 {
			return true, io.EOF
		}
		count++
		return false, nil
	})

	assert.EqualError(t, err, io.EOF.Error())
}

func TestRetryCancel(t *T) {
	ctx, cancel := context.WithCancel(context.Background())
	go cancel()
	s := utils.BackoffSleeper(time.Second, time.Second, nil)

	err := utils.Retry(ctx, s, func() (bool, error) {
		return false, nil
	})

	assert.EqualError(t, err, context.Canceled.Error())
}

func TestCountSleeperErr(t *T) {
	ctx := context.Background()
	s := utils.CountSleeper(5)
	for i := 0; i < 5; i++ {
		_ = s(ctx)
	}
	assert.Errorf(t, s(ctx), utils.ErrMaxSleepCount.Error())
}

func TestCountSleeperCancel(t *T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	s := utils.CountSleeper(5)
	assert.Errorf(t, s(ctx), context.Canceled.Error())
}
