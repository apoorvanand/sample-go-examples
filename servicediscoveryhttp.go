package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	consulapi "github.com/hashicorp/consul/api"
)

func main() {
	config := consulapi.DefaultConfig()
	client, err := consulapi.NewClient(config)
	if err != nil {
		log.Fatal(err)
	}

	port := 8080
	hostname := "localhost"

	err = client.Agent().ServiceRegister(&consulapi.AgentServiceRegistration{
		ID:      "service-1",
		Name:    "api",
		Address: hostname,
		Port:    port,
	})
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		services, _, err := client.Catalog().Service("api", "", nil)
		if err != nil {
			log.Fatal(err)
		}

		for _, service := range services {
			url := fmt.Sprintf("http://%s:%d", service.ServiceAddress, service.ServicePort)
			resp, err := http.Get(url)
			if err != nil {
				log.Printf("Error: %v", err)
				continue
			}
			fmt.Fprintf(w, "Service: %s, Response: %d\n", service.ServiceName, resp.StatusCode)
		}
	})

	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), nil))
}
