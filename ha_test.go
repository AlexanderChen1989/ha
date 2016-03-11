package ha

import (
	"testing"
	"time"

	"golang.org/x/net/context"

	"github.com/stretchr/testify/assert"
)

func TestRestartTimes(t *testing.T) {
	times := 0
	Watch(func() {
		times += 1
		panic("")
	}, RestartTimes(100))
	assert.Equal(t, 100, times)
}

func TestRestartDelay(t *testing.T) {
	const num = 100
	var tps []time.Time
	Watch(func() {
		tps = append(tps, time.Now())
		panic("")
	}, RestartDelay(num*time.Millisecond), RestartTimes(10))

	for i := range tps {
		if i == 0 {
			continue
		}
		assert.Equal(t, int(tps[i].Sub(tps[i-1]).Seconds()*1000/num), 1)
	}
}

func TestOnStop(t *testing.T) {
	errLogged := false
	Watch(func() {
		panic("")
	}, RestartTimes(1), OnStop(func(err error) {
		errLogged = true
	}))

	assert.True(t, errLogged)
}

func TestCancelCtx(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	start := time.Now()
	Watch(func() {
		panic("")
	}, CancelCtx(ctx))
	assert.Equal(t, int(time.Now().Sub(start).Seconds()), 2)
}
