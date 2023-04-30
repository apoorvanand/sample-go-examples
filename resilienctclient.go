package main
import (
	"context"
	"errors"
	"math/rand"
	"net/http"
	"time"
)

const maxRetryAttempts = 5

func makeRequest(ctx context.Context, url string) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	var res *http.Response
	var attempt int

	for {
		// exponential backoff and jitter
		backoff := time.Duration(1<<attempt) * time.Second
		backoff += time.Duration(rand.Intn(1000)) * time.Millisecond
		time.Sleep(backoff)

		// send request
		res, err = client.Do(req)
		if err == nil {
			return res, nil
		}

		// check if retry is allowed
		if !errors.Is(err, context.Canceled) && attempt < maxRetryAttempts {
			attempt++
			continue
		}

		break
	}

	return nil, err
}
