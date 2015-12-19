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
			err = fmt.Errorf("%s", err)
		}
	}()

	fn()
	return
}

type config struct {
	max     uint
	noLimit bool
	wait    time.Duration
	onStop  func(error)
	ctx     context.Context
}

// Context when ctx cancled, fn will not restart again
func Context(ctx context.Context) func(conf *config) {
	return func(conf *config) {
		conf.ctx = ctx
	}
}

// OnStop run onStop when fn stopped
func OnStop(onStop func(error)) func(conf *config) {
	return func(conf *config) {
		conf.onStop = onStop
	}
}

// Max max restart times
func Max(max uint) func(conf *config) {
	return func(conf *config) {
		conf.noLimit = false
		conf.max = max
	}
}

// Wait wait for d time before rerun fn
func Wait(d time.Duration) func(conf *config) {
	return func(conf *config) {
		conf.wait = d
	}
}

// Watch run then watch fn, if fn returned, run it again
func Watch(fn func(), setups ...func(*config)) {
	conf := &config{
		noLimit: true,
		onStop:  func(error) {},
		ctx:     context.Background(),
	}
	for _, setup := range setups {
		setup(conf)
	}
	for {
		conf.onStop(recoverWrapper(fn))
		select {
		case <-conf.ctx.Done():
			return
		default:
		}
		time.Sleep(conf.wait)
		if conf.noLimit {
			continue
		}
		if conf.max--; conf.max <= 0 {
			return
		}
	}
}
