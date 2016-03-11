package ha

import (
	"testing"
	"time"

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
		assert.Equal(t, int(tps[i].Sub(tps[i-1]).Seconds()*1000/100), 1)
	}
}
