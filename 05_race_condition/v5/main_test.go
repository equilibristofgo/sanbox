package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
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

	mdb2b := NewMemoryDB()
	cache := NewLockFreeJobCache(mdb2b)

	mockRemoteJob.Init(cache)
	cache.Start()
	mockRemoteJob.Run(cache)

	mockRemoteJob.lock.RLock()
	assert.True(t, mockRemoteJob.Metadata.SuccessCount == 1)
	mockRemoteJob.lock.RUnlock()
}
