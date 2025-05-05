package main

import (
	"fmt"
	"sync"
	"time"
)

type RateLimiter struct {
	// amount of requests per second
	doneReq         uint          // internal
	reqRate         uint          // 10 req
	reqRateDuration time.Duration // 1 sec
	lastReqTime     time.Time     // internal
	mu              sync.Mutex
}

func (r *RateLimiter) Allow() bool {
	reqTime := time.Now()
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.lastReqTime.Sub(reqTime).Seconds() < float64(r.reqRateDuration) {
		if r.doneReq >= r.reqRate {
			return false
		}
		r.doneReq++
	}
	r.lastReqTime = time.Now()
	r.doneReq = 0
	return true
}

func main() {

	t := time.Now()

	t2 := time.Now().Add(1 * time.Second)

	tokens := 10.0
	fmt.Println(t.Sub(t2).Seconds())

	tokens += t.Sub(t2).Seconds() * 0.05
	fmt.Println(tokens)
}
