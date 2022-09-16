package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"testing"
	"time"
)

type Metadata struct {
	SuccessCount uint
}

type RemoteProperties struct {
	Headers               http.Header `json:"headers"`
	Url                   string      `json:"url"`
	ExpectedResponseCodes []int       `json:"expected_response_codes"`
}

type JobsMap struct {
	Jobs map[string]*Job
	Lock sync.RWMutex
}

func NewJobsMap() *JobsMap {
	return &JobsMap{
		Jobs: map[string]*Job{},
		Lock: sync.RWMutex{},
	}
}

type JobCache interface {
	GetAll() *JobsMap
	Set(j *Job) error
}

type JobDB interface {
	GetAll() ([]*Job, error)
}

type MockDBGetAll struct {
	MockDB
	response []*Job
}

func (d *MockDBGetAll) GetAll() ([]*Job, error) {
	return d.response, nil
}

func parseTime(t *testing.T, value string) time.Time {
	now, err := time.Parse("2006-Jan-02 15:04", value)
	if err != nil {
		t.Fatal(err)
	}
	return now
}

var _ JobDB = (*MemoryDB)(nil)

func main() {

	mdb2b := NewMemoryDB()
	c := NewLockFreeJobCache(mdb2b)

	j := &Job{Metadata: Metadata{
		SuccessCount: 0,
	}}

	j.Init(c)
	c.Start()

	j.Run(c)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	done := make(chan bool, 1)
	fmt.Println("awaiting signal")
	<-done
	fmt.Println("exiting")

}
