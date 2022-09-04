package main

import (
	"fmt"
	"net/http"
	"sync"
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

type Cache struct {
}

func (c *Cache) Set(j *Job) {
}

type Job struct {
	Metadata Metadata
	RemoteProperties
	lock sync.RWMutex
}

// Other type...
func (j *Job) StartWaiting(cache Cache, wg sync.WaitGroup) {

	j.lock.Lock()
	defer j.lock.Unlock()

	jobRun := func() { j.Run(cache, wg) }

	d := time.Duration(3) * time.Second
	t := time.AfterFunc(d, jobRun)

	// Calling stop method
	// w.r.to Timer1
	defer t.Stop()

	// Calling sleep method
	time.Sleep(10 * time.Second)

	t.Reset(d)
}

func (j *Job) Run(c Cache, wg sync.WaitGroup) {
	j.lock.RLock()
	jobRunner := &JobRunner{job: j, meta: j.Metadata}
	j.lock.RUnlock()

	meta := jobRunner.Run(c)

	j.lock.Lock()
	j.Metadata = meta
	j.lock.Unlock()

	j.lock.RLock()
	c.Set(j)
	j.lock.RUnlock()

	j.lock.Lock()
	go j.StartWaiting(c, wg)
	j.lock.Unlock()

	fmt.Println("Run")
	wg.Done()
}

func (j *JobRunner) RemoteRun() {
	req, _ := http.NewRequest(http.MethodGet, "url", nil)
	j.setHeaders(req)
}

type JobRunner struct {
	job  *Job
	meta Metadata
}

func (j *JobRunner) Run(c Cache) Metadata {
	j.job.lock.RLock()
	defer j.job.lock.RUnlock()

	j.RemoteRun()

	j.meta.SuccessCount++

	return j.meta
}

func (j *JobRunner) checkExpected(statusCode int) bool {
	// If no expected response codes passed, add 200 status code as expected
	if len(j.job.RemoteProperties.ExpectedResponseCodes) == 0 {
		j.job.RemoteProperties.ExpectedResponseCodes = append(j.job.RemoteProperties.ExpectedResponseCodes, 200)
	}
	for _, expected := range j.job.RemoteProperties.ExpectedResponseCodes {
		if expected == statusCode {
			return true
		}
	}

	return false
}

// setHeaders sets default and user specific headers to the http request
func (j *JobRunner) setHeaders(req *http.Request) {
	if j.job.RemoteProperties.Headers == nil {
		j.job.RemoteProperties.Headers = http.Header{}
	}
	// A valid assumption is that the user is sending something in json cause we're past 2017
	if j.job.RemoteProperties.Headers["Content-Type"] == nil {
		j.job.RemoteProperties.Headers["Content-Type"] = []string{"application/json"}
	}
	req.Header = j.job.RemoteProperties.Headers
}

func main() {

	var wg sync.WaitGroup

	c := Cache{}
	j := &Job{Metadata: Metadata{
		SuccessCount: 0,
	}}

	wg.Add(1)

	go j.Run(c, wg)

	wg.Wait()

}
