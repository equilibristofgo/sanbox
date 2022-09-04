package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRemoteJobRunner(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, client")
	}))
	defer testServer.Close()

	mockRemoteJob := &Job{
		Metadata: Metadata{
			SuccessCount: 0,
		},
		RemoteProperties: RemoteProperties{
			Headers: nil,
			Url:     testServer.URL,
		},
	}

	cache := Cache{}
	// mockRemoteJob.Init(cache)
	// cache.Start(0, 2*time.Second) // Retain 1 minute

	var wg sync.WaitGroup
	wg.Add(1)
	mockRemoteJob.Run(cache, wg)

	mockRemoteJob.lock.RLock()
	assert.True(t, mockRemoteJob.Metadata.SuccessCount == 1)
	mockRemoteJob.lock.RUnlock()
}
