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

func (m mockDBClient) GetMockedDatabases() ([]byte, error) {
	result := map[string]map[string][]string{}

	for dbName, db := range m.client {
		if result[dbName] == nil {
			result[dbName] = map[string][]string{}
		}
		for tblName, tbl := range db {
			if result[dbName][tblName] == nil {
				result[dbName][tblName] = []string{}
			}
			for _, fld := range tbl {
				result[dbName][tblName] = append(result[dbName][tblName], string(fld))
			}
		}
	}

	return json.Marshal(result)
}

func (m mockDBClient) SetMockedDatabases(data []byte) error {
	newClient := map[string]map[string][]string{}

	jsonErr := json.Unmarshal(data, &newClient)
	if jsonErr != nil {
		return jsonErr
	}

	for dbName, db := range newClient {
		if m.client[dbName] == nil {
			m.client[dbName] = map[string][][]byte{}
		}
		for tblName, tbl := range db {
			if m.client[dbName][tblName] == nil {
				m.client[dbName][tblName] = [][]byte{}
			}
			for _, fld := range tbl {
				m.client[dbName][tblName] = append(m.client[dbName][tblName], []byte(fld))
			}
		}
	}

	return nil
}

func (m mockDBClient) FindOne(database string, table string, matches map[string]interface{}, receiver interface{}) error {
	value := map[string]interface{}{}
	for _, v := range m.client[database][table] {
		json.Unmarshal(v, &value)
		matching := true
		for matchKey, matchValue := range matches {
			if value[matchKey] != matchValue {
				matching = false
			}
		}
		if matching {
			json.Unmarshal(v, &receiver)
		}
	}

	return nil
}

func (m mockDBClient) FindMany(database string, table string, matches map[string]interface{}) (ICursor, error) {
	cursor := cursorStruct{
		results: [][]byte{},
		current: -1,
	}

	value := map[string]interface{}{}
	for _, v := range m.client[database][table] {
		json.Unmarshal(v, &value)
		matching := true
		for matchKey, matchValue := range matches {
			if value[matchKey] != matchValue {
				matching = false
			}
		}
		if matching {
			cursor.results = append(cursor.results, v)
		}
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

func (m mockDBClient) InsertMany(database string, table string, provider []interface{}) error {
	if m.client[database] == nil {
		m.client[database] = make(map[string][][]byte)
	}

	if m.client[database][table] == nil {
		m.client[database][table] = [][]byte{}
	}

	for _, prov := range provider {
		value, marshalErr := json.Marshal(prov)
		if marshalErr != nil {
			return marshalErr
		}

		newTable := append(m.client[database][table], value)

		m.client[database][table] = newTable
	}

	return nil
}

func (m mockDBClient) UpdateOne(database string, table string, matches map[string]interface{}, provider interface{}) error {
	value := map[string]interface{}{}
	for k, v := range m.client[database][table] {
		json.Unmarshal(v, &value)
		matching := true
		for matchKey, matchValue := range matches {
			if value[matchKey] != matchValue {
				matching = false
			}
		}
		if matching {
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

func (m mockDBClient) DeleteOne(database string, table string, matches map[string]interface{}) error {
	value := map[string]interface{}{}
	for k, v := range m.client[database][table] {
		json.Unmarshal(v, &value)
		matching := true
		for matchKey, matchValue := range matches {
			if value[matchKey] != matchValue {
				matching = false
			}
		}
		if matching {
			m.client[database][table] = append(m.client[database][table][:k], m.client[database][table][k+1:]...)
			return nil
		}
	}

	return nil
}

func (m mockDBClient) DeleteMany(database string, table string, matches map[string]interface{}) error {
	value := map[string]interface{}{}

	toDelete := []int{}

	for k, v := range m.client[database][table] {
		json.Unmarshal(v, &value)
		matching := true
		for matchKey, matchValue := range matches {
			if value[matchKey] != matchValue {
				matching = false
			}
		}
		if matching {
			toDelete = append(toDelete, k)
		}
	}

	for index, key := range toDelete {
		i := key - index
		m.client[database][table] = append(m.client[database][table][:i], m.client[database][table][i+1:]...)
	}

	return nil
}

// MockDatabaseUtil is an implementation of IDatabaseUtil that stores data locally. Should not be used if not for testing purposes.
var MockDatabaseUtil = mockDBClient{
	client: mockDBClientFactory(),
}
