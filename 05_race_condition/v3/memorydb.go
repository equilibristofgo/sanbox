package main

import (
	"errors"
	"sync"
)

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
