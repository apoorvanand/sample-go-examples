package main

import (
	"fmt"
	"net/http"
	"sync"
//"time"
)

func main() {
	var wg sync.WaitGroup
	for i := 1; i <= 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 1; j <= 1000; j++ {
				resp, err := http.Get("http://example.com")
				if err != nil {
					fmt.Println("Error:", err)
				} else {
					fmt.Println(resp.Status)
					resp.Body.Close()
				}
			}
		}()
	}
	wg.Wait()
	fmt.Println("Load test complete")
}
