package main

import (
	//"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	//"time"

	"github.com/SherClockHolmes/webpush-go"
)

// Subscription represents the push subscription for a user
type Subscription struct {
	Endpoint string `json:"endpoint"`
	Keys     struct {
		P256dh string `json:"p256dh"`
		Auth   string `json:"auth"`
	} `json:"keys"`
}

func main() {
	// Load the push subscription from a JSON file
	subscriptionFile, err := os.Open("subscription.json")
	if err != nil {
		panic(err)
	}
	defer subscriptionFile.Close()

	var subscription Subscription
	if err := json.NewDecoder(subscriptionFile).Decode(&subscription); err != nil {
		panic(err)
	}

	// Create the web push notification payload
	payload := []byte("Hello, world!")

	// Set the web push notification options
	options := &webpush.Options{
		Subscriber:      "example@example.com",
		TTL:             30 * 2,//time.Second,
		VAPIDPrivateKey: "vauNnbWNtXwl8O_8bJhEyjnKUrBTvCjfdoyIJgEZpbA=",//"VAPID_PRIVATE_KEY",
		VAPIDPublicKey:  "jIM6fZ4vyxRoYLGevEpvxEOyDgyeFuXxtmdYQD3sVdw=",//"VAPID_PUBLIC_KEY",
	}

	// Send the web push notification
	resp, err := webpush.SendNotification(payload, &webpush.Subscription{
		Endpoint: subscription.Endpoint,
		Keys: webpush.Keys{
			P256dh: subscription.Keys.P256dh,
			Auth:   subscription.Keys.Auth,
		},
	}, options)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		panic(fmt.Errorf("unexpected status code %d", resp.StatusCode))
	}
}
