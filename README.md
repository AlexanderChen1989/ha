# HA For Go

#### Add a little high availability to your Go functions or methods.
## Simple Usage
```go
ha.Watch(func() {...})
ha.Watch(
	func() {...},
	ha.RestartDelay(time.Second),
	// ha.MaxRestart(5),
	// ha.CancelCtx(ctx),
 )
```
## Caution!
Panic in goroutine may panic runtime! Watch is only watch current goroutine.
```go
	Watch(func() {
		go func() {
			panic("No!!") // Will panic runtime
		}()
	})

	// you can do this
	// func will be restart when panic happened
	go Watch(func() {
		panic("No!!")
	})
```


## Example
```go
package main

import (
	"fmt"
	"time"

	"golang.org/x/net/context"

	"github.com/AlexanderChen1989/ha"
)

func main() {
	// a function that shoud be watched
	fn := func() {
		var i *int
		*i = 100
	}

	// when fn stop
	onStop := ha.OnStop(func(err error) {
		fmt.Println(err)
	})

	// try to restart fn at most 5 times
	max := ha.MaxRestart(5)

	// wait sometime before restart fn
	delay := ha.RestartDelay(1 * time.Second)

	// context to stop restarting fn
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		// sleep 2 seconds
		time.Sleep(2 * time.Second)
		cancel()
	}()

	ha.Watch(fn, onStop, delay, max, ha.CancelCtx(ctx))
}


```
