OneShot (一发)
---
OneShot provide a high performance way for pause and resume goroutine once.

## Example

```go
package main

import (
	"github.com/cocotyty/oneshot"
	"log"
	"time"
)

func main() {
	s := &oneshot.Shot{}
	go func() {
		time.Sleep(time.Second)
		s.WakeUp()
	}()
	log.Println("wait")
	s.Wait()
	log.Println("finished")
}
```