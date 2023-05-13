package main

import (
    "errors"
    "fmt"
    "time"
)

type CircuitBreaker struct {
    failureThreshold int
    retryTimeout     time.Duration
    cooldownTimeout  time.Duration
    failures         int
    lastFailure      time.Time
    open             bool
}

func (cb *CircuitBreaker) AllowRequest() bool {
    if cb.open {
        timeSinceLastFailure := time.Now().Sub(cb.lastFailure)
        if timeSinceLastFailure >= cb.cooldownTimeout {
            cb.open = false
        } else {
            return false
        }
    }
    return true
}

func (cb *CircuitBreaker) ReportFailure() {
    cb.failures++
    if cb.failures >= cb.failureThreshold {
        cb.lastFailure = time.Now()
        cb.open = true
        go cb.autoRetry()
    }
}

func (cb *CircuitBreaker) autoRetry() {
    time.Sleep(cb.retryTimeout)
    cb.failures = 0
}

func main() {
    cb := &CircuitBreaker{
        failureThreshold: 2,
        retryTimeout:     time.Second * 3,
        cooldownTimeout:  time.Second * 10,
        failures:         0,
        lastFailure:      time.Time{},
        open:             false,
    }

    for i := 0; i < 5; i++ {
        if cb.AllowRequest() {
            // make API request here
            // simulate failure for first two requests
            if i < 2 {
                fmt.Println("API request failed")
                cb.ReportFailure()
            } else {
                fmt.Println("API request succeeded")
            }
        } else {
            fmt.Println("Circuit breaker open, request denied")
        }
        time.Sleep(time.Second * 2)
    }
}
