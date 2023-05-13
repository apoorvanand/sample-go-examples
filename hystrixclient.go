package main

import (
    "errors"
    "fmt"
    "github.com/afex/hystrix-go/hystrix"
    "net/http"
    "time"
)

func main() {
    hystrix.ConfigureCommand("my_api_command", hystrix.CommandConfig{
        Timeout:               1000, // Timeout in milliseconds
        MaxConcurrentRequests: 10,   // Max concurrent requests allowed
        ErrorPercentThreshold: 50,   // Error percentage threshold for circuit breaker
        RequestVolumeThreshold: 5,   // Minimum number of requests needed to trip circuit breaker
        SleepWindow:           5000, // Time window for circuit breaker to stay open
    })

    client := &http.Client{
        Timeout: time.Second * 5,
    }

    for i := 0; i < 10; i++ {
        err := hystrix.Do("my_api_command", func() error {
            req, _ := http.NewRequest("GET", "https://jsonplaceholder.typicode.com/todos/1", nil)
            resp, err := client.Do(req)
            if err != nil {
                return err
            }
            defer resp.Body.Close()
            if resp.StatusCode != http.StatusOK {
                return errors.New("API request failed")
            }
            // process response data here
            fmt.Println("API request succeeded")
            return nil
        }, func(err error) error {
            // fallback function for circuit breaker
            fmt.Println("Circuit breaker opened")
            return err
        })

        if err != nil {
            fmt.Println("API request failed:", err.Error())
        }

        time.Sleep(time.Second * 2)
    }
}
