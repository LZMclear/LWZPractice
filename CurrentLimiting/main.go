package main

import (
	"fmt"
	"go.uber.org/ratelimit"
	"time"
)

func main() {
	rl := ratelimit.New(100)
	prev := time.Now()
	for i := 0; i < 10; i++ {
		now := rl.Take()
		fmt.Println(i, now.Sub(prev))
		prev = now
	}
}
