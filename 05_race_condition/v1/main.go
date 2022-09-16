package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/cornelk/hashmap"
)

type Metadata struct {
	SuccessCount uint
}

type RemoteProperties struct {
	Headers               http.Header `json:"headers"`
	Url                   string      `json:"url"`
	ExpectedResponseCodes []int       `json:"expected_response_codes"`
}

type LockFreeJobCache struct {
	jobs  *hashmap.HashMap
	jobDB JobDB
}

type JobCache interface {
	GetAll() *JobsMap
	Set(j *Job) error
}

type JobDB interface {
	GetAll() ([]*Job, error)
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

func NewLockFreeJobCache(jobDB JobDB) *LockFreeJobCache {
	return &LockFreeJobCache{
		jobs:  hashmap.New(8), //nolint:gomnd
		jobDB: jobDB,
	}
}

func (c *LockFreeJobCache) Set(j *Job) error {
	if j == nil {
		return nil
	}

	c.jobs.Set(j.Id, j)
	return nil
}

func (c *LockFreeJobCache) GetAll() *JobsMap {
	jm := NewJobsMap()
	for el := range c.jobs.Iter() {
		jm.Jobs[el.Key.(string)] = el.Value.(*Job)
	}
	return jm
}

func (c *LockFreeJobCache) Start() {
	allJobs, _ := c.jobDB.GetAll()
	for _, j := range allJobs {
		j.StartWaiting(c)
	}
}

type Job struct {
	Id       string
	Metadata Metadata
	RemoteProperties
	lock sync.RWMutex
}

// Other type...
func (j *Job) StartWaiting(cache JobCache) {

	j.lock.Lock()
	defer j.lock.Unlock()

	jobRun := func() { j.Run(cache) }

	d := time.Duration(1) * time.Second
	t := time.AfterFunc(d, jobRun)

	// Calling stop method
	// w.r.to Timer1
	defer t.Stop()

	// Calling sleep method
	time.Sleep(2 * time.Second)

	t.Reset(d)
}

func (j *Job) Run(c JobCache) {
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
	go j.StartWaiting(c)
	j.lock.Unlock()

	fmt.Println("Run")
}

func (j *JobRunner) RemoteRun() {
	req, _ := http.NewRequest(http.MethodGet, "url", nil)
	j.setHeaders(req)
	// Do the request
	http.DefaultClient.Do(req.WithContext(context.Background()))
}

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

func (j *Job) Init(c JobCache) error {
	j.lock.Lock()
	defer j.lock.Unlock()

	c.Set(j)
	go j.Run(c)

	j.lock.Unlock()
	j.StartWaiting(c)
	j.lock.Lock()

	return nil
}

type MockDBGetAll struct {
	MockDB
	response []*Job
}

func (d *MockDBGetAll) GetAll() ([]*Job, error) {
	return d.response, nil
}

type MockDB struct{}

func (m *MockDB) GetAll() ([]*Job, error) {
	return nil, nil
}
func (m *MockDB) Get(id string) (*Job, error) {
	return nil, nil
}
func (m *MockDB) Delete(id string) error {
	return nil
}
func (m *MockDB) Save(job *Job) error {
	return nil
}
func (m *MockDB) Close() error {
	return nil
}

func NewMockCache() *LockFreeJobCache {
	return NewLockFreeJobCache(&MockDB{})
}

func parseTime(t *testing.T, value string) time.Time {
	now, err := time.Parse("2006-Jan-02 15:04", value)
	if err != nil {
		t.Fatal(err)
	}
	return now
}

var _ JobDB = (*MemoryDB)(nil)

type MemoryDB struct {
	m    map[string]*Job
	lock sync.RWMutex
}

func NewMemoryDB() *MemoryDB {
	return &MemoryDB{
		m: map[string]*Job{},
	}
}

func (m *MemoryDB) GetAll() (ret []*Job, _ error) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	for _, v := range m.m {
		ret = append(ret, v)
	}
	return
}

func (m *MemoryDB) Get(id string) (*Job, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	j, exist := m.m[id]
	if !exist {
		return nil, errors.New("NotFound")
	}
	return j, nil
}

func (m *MemoryDB) Delete(id string) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	if _, exists := m.m[id]; !exists {
		return errors.New("Doesn't exist") // Used for testing
	}
	delete(m.m, id)
	// log.Printf("After delete: %+v", m)
	return nil
}

func (m *MemoryDB) Save(j *Job) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.m[j.Id] = j
	return nil
}

func (m *MemoryDB) Close() error {
	return nil
}

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
