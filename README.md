# HA For Go

#### Add a little high availability to your Go functions or methods.

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
	max := ha.Max(5)

	// wait sometime before restart fn
	wait := ha.Wait(1 * time.Second)

	// context to stop restarting fn
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		// sleep 2 seconds
		time.Sleep(2 * time.Second)
		cancel()
	}()

	ha.Watch(fn, onStop, wait, max, ha.Context(ctx))
}
```
