package main

type DataStore struct {
	data map[string]string
}

func NewDataStore() *DataStore {
	return &DataStore{data: make(map[string]string)}
}

func (dataStore *DataStore) Get(key string) string {
	return dataStore.data[key]
}

func (dataStore *DataStore) Set(key string, value string) {
	dataStore.data[key] = value
}
