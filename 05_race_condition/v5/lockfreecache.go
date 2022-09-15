package main

import "github.com/cornelk/hashmap"

type LockFreeJobCache struct {
	jobs  *hashmap.HashMap
	jobDB JobDB
}

func NewMockCache() *LockFreeJobCache {
	return NewLockFreeJobCache(&MockDB{})
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
