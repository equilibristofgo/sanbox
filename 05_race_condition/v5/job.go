package main

import (
	"fmt"
	"sync"
	"time"
)

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
