package database

import (
	"context"
	"encoding/json"
)

func mockDBClientFactory() map[string]map[string][][]byte {
	store := map[string]map[string][][]byte{}

	return store
}

type mockDBClient struct {
	client map[string]map[string][][]byte
}

type cursorStruct struct {
	results [][]byte
	current int
}

func (c *cursorStruct) Decode(receiver interface{}) error {
	json.Unmarshal(c.results[c.current], &receiver)
	return nil
}

func (c *cursorStruct) Next(ctx context.Context) bool {
	if len(c.results) == c.current+1 {
		return false
	}
	c.current = c.current + 1
	return true
}

func (c *cursorStruct) Close(ctx context.Context) error {
	c.results = [][]byte{}
	c.current = -1
	return nil
}

func (m mockDBClient) FindOne(database string, table string, field string, match interface{}, receiver interface{}) error {
	value := map[string]interface{}{}
	for _, v := range m.client[database][table] {
		json.Unmarshal(v, &value)
		if value[field] == match {
			json.Unmarshal(v, &receiver)
			return nil
		}
	}

	return nil
}

func (m mockDBClient) FindMany(database string, table string, field string, match interface{}) (ICursor, error) {
	cursor := cursorStruct{
		results: [][]byte{},
		current: -1,
	}

	value := map[string]interface{}{}
	for _, v := range m.client[database][table] {
		json.Unmarshal(v, &value)
		if value[field] == match {
			cursor.results = append(cursor.results, v)
		}
	}

	return &cursor, nil
}

func (m mockDBClient) FindAll(database string, table string) (ICursor, error) {
	cursor := cursorStruct{
		results: [][]byte{},
		current: -1,
	}

	value := map[string]interface{}{}
	for _, v := range m.client[database][table] {
		json.Unmarshal(v, &value)
		cursor.results = append(cursor.results, v)
	}

	return &cursor, nil
}

func (m mockDBClient) InsertOne(database string, table string, provider interface{}) error {
	if m.client[database] == nil {
		m.client[database] = make(map[string][][]byte)
	}

	if m.client[database][table] == nil {
		m.client[database][table] = [][]byte{}
	}

	value, marshalErr := json.Marshal(provider)
	if marshalErr != nil {
		return marshalErr
	}

	newTable := append(m.client[database][table], value)

	m.client[database][table] = newTable

	return nil
}

func (m mockDBClient) UpdateOne(database string, table string, field string, match interface{}, provider interface{}) error {
	value := map[string]interface{}{}
	for k, v := range m.client[database][table] {
		json.Unmarshal(v, &value)
		if value[field] == match {
			value, marshalErr := json.Marshal(provider)
			if marshalErr != nil {
				return marshalErr
			}
			m.client[database][table][k] = value
			return nil
		}
	}

	return nil
}

func (m mockDBClient) DeleteOne(database string, table string, field string, match interface{}) error {
	value := map[string]interface{}{}
	for k, v := range m.client[database][table] {
		json.Unmarshal(v, &value)
		if value[field] == match {
			newTable := append(m.client[database][table][:k], m.client[database][table][k+1:]...)
			m.client[database][table] = newTable
			return nil
		}
	}

	return nil
}

// MockDatabaseUtil is an implementation of IDatabaseUtil that stores data locally. Should not be used if not for testing purposes.
var MockDatabaseUtil = mockDBClient{
	client: mockDBClientFactory(),
}
