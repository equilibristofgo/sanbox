package main

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
