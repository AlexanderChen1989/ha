package ha

import (
	"fmt"
	"time"

	"golang.org/x/net/context"
)

func recoverWrapper(fn func()) (err error) {
	defer func() {
		switch e := recover().(type) {
		case nil:
		case error:
			err = e
		case interface{}:
			err = fmt.Errorf("%s", e)
		}
	}()

	fn()
	return
}

type config struct {
	maxRestart   uint
	noLimit      bool
	restartDelay time.Duration
	onStop       func(error)
	cancelCtx    context.Context
}

// CancelCtx when ctx cancled, fn will not restart again
func CancelCtx(ctx context.Context) func(conf *config) {
	return func(conf *config) {
		conf.cancelCtx = ctx
	}
}

// OnStop run onStop when fn stopped
func OnStop(onStop func(error)) func(conf *config) {
	return func(conf *config) {
		conf.onStop = onStop
	}
}

// MaxRestart max restart times
func MaxRestart(max uint) func(conf *config) {
	return func(conf *config) {
		conf.noLimit = false
		conf.maxRestart = max
	}
}

// RestartDelay wait for d time before restart fn
func RestartDelay(d time.Duration) func(conf *config) {
	return func(conf *config) {
		conf.restartDelay = d
	}
}

// Watch run then watch fn, if fn returned, run it again
func Watch(fn func(), setups ...func(*config)) {
	conf := &config{
		noLimit:   true,
		onStop:    func(error) {},
		cancelCtx: context.Background(),
	}
	for _, setup := range setups {
		setup(conf)
	}
	for {
		conf.onStop(recoverWrapper(fn))
		select {
		case <-conf.cancelCtx.Done():
			return
		default:
		}
		time.Sleep(conf.restartDelay)
		if conf.noLimit {
			continue
		}
		if conf.maxRestart--; conf.maxRestart <= 0 {
			return
		}
	}
}
