package main

import (
	"fmt"
	"time"

	"golang.org/x/net/context"

	"github.com/AlexanderChen1989/ha"
)

func main() {
	fn := func() {
		var i *int
		*i = 100
	}
	onStop := ha.OnStop(func(err error) {
		fmt.Println(err)
	})
	max := ha.Max(5)
	wait := ha.Wait(1 * time.Second)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(2 * time.Second)
		cancel()
	}()
	ha.Watch(fn, onStop, wait, max, ha.Context(ctx))
}
