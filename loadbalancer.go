package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
)

type loadBalancer struct {
	servers []*url.URL
	current int
	mutex   sync.Mutex
}

func (lb *loadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	lb.mutex.Lock()
	serverURL := lb.servers[lb.current]
	lb.current = (lb.current + 1) % len(lb.servers)
	lb.mutex.Unlock()

	proxy := httputil.NewSingleHostReverseProxy(serverURL)
	proxy.ServeHTTP(w, r)
}

func main() {
	serverUrls := []*url.URL{
		{
			Scheme: "http",
			Host:   "localhost:8080",
		},
		{
			Scheme: "http",
			Host:   "localhost:8081",
		},
	}

	loadBalancer := loadBalancer{
		servers: serverUrls,
	}
	fmt.Print("server running on port 8000")
	http.ListenAndServe(":8000", &loadBalancer)
	
}
