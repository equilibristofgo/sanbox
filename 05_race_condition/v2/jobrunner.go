package main

import (
	"context"
	"net/http"
)

type JobRunner struct {
	job  *Job
	meta Metadata
}

func (j *JobRunner) Run(c JobCache) Metadata {
	j.job.lock.RLock()
	defer j.job.lock.RUnlock()

	j.RemoteRun()

	j.meta.SuccessCount++

	return j.meta
}

func (j *JobRunner) RemoteRun() {
	req, _ := http.NewRequest(http.MethodGet, "url", nil)
	j.setHeaders(req)
	// Do the request
	http.DefaultClient.Do(req.WithContext(context.Background()))
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
	j.job.lock.RLock()
	defer j.job.lock.RUnlock()
	if j.job.RemoteProperties.Headers == nil {
		j.job.RemoteProperties.Headers = http.Header{}
	}
	// A valid assumption is that the user is sending something in json cause we're past 2017
	if j.job.RemoteProperties.Headers["Content-Type"] == nil {
		j.job.RemoteProperties.Headers["Content-Type"] = []string{"application/json"}
	}
	req.Header = j.job.RemoteProperties.Headers
}
