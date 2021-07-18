package database

import (
	"context"
	"encoding/json"
)

func FakeDBClientFactory() map[string]map[string][][]byte {
	store := map[string]map[string][][]byte{}

	return store
}

type FakeDatabaseUtil struct {
	Client map[string]map[string][][]byte
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

func (m FakeDatabaseUtil) GetMockedDatabases() ([]byte, error) {
	result := map[string]map[string][]string{}

	for dbName, db := range m.Client {
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

func (m FakeDatabaseUtil) SetMockedDatabases(data []byte) error {
	newClient := map[string]map[string][]string{}

	jsonErr := json.Unmarshal(data, &newClient)
	if jsonErr != nil {
		return jsonErr
	}

	for dbName, db := range newClient {
		if m.Client[dbName] == nil {
			m.Client[dbName] = map[string][][]byte{}
		}
		for tblName, tbl := range db {
			if m.Client[dbName][tblName] == nil {
				m.Client[dbName][tblName] = [][]byte{}
			}
			for _, fld := range tbl {
				m.Client[dbName][tblName] = append(m.Client[dbName][tblName], []byte(fld))
			}
		}
	}

	return nil
}

func (m FakeDatabaseUtil) Connect(uri string) error {
	return nil
}

func (m FakeDatabaseUtil) FindOne(database string, table string, matches map[string]interface{}, receiver interface{}) error {
	value := map[string]interface{}{}
	for _, v := range m.Client[database][table] {
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

func (m FakeDatabaseUtil) FindMany(database string, table string, matches map[string]interface{}) (ICursor, error) {
	cursor := cursorStruct{
		results: [][]byte{},
		current: -1,
	}

	value := map[string]interface{}{}
	for _, v := range m.Client[database][table] {
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

func (m FakeDatabaseUtil) InsertOne(database string, table string, provider interface{}) error {
	if m.Client[database] == nil {
		m.Client[database] = make(map[string][][]byte)
	}

	if m.Client[database][table] == nil {
		m.Client[database][table] = [][]byte{}
	}

	value, marshalErr := json.Marshal(provider)
	if marshalErr != nil {
		return marshalErr
	}

	newTable := append(m.Client[database][table], value)

	m.Client[database][table] = newTable

	return nil
}

func (m FakeDatabaseUtil) InsertMany(database string, table string, provider []interface{}) error {
	if m.Client[database] == nil {
		m.Client[database] = make(map[string][][]byte)
	}

	if m.Client[database][table] == nil {
		m.Client[database][table] = [][]byte{}
	}

	for _, prov := range provider {
		value, marshalErr := json.Marshal(prov)
		if marshalErr != nil {
			return marshalErr
		}

		newTable := append(m.Client[database][table], value)

		m.Client[database][table] = newTable
	}

	return nil
}

func (m FakeDatabaseUtil) UpdateOne(database string, table string, matches map[string]interface{}, provider interface{}) error {
	value := map[string]interface{}{}
	for k, v := range m.Client[database][table] {
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
			m.Client[database][table][k] = value
			return nil
		}
	}

	return nil
}

func (m FakeDatabaseUtil) DeleteOne(database string, table string, matches map[string]interface{}) error {
	value := map[string]interface{}{}
	for k, v := range m.Client[database][table] {
		json.Unmarshal(v, &value)
		matching := true
		for matchKey, matchValue := range matches {
			if value[matchKey] != matchValue {
				matching = false
			}
		}
		if matching {
			m.Client[database][table] = append(m.Client[database][table][:k], m.Client[database][table][k+1:]...)
			return nil
		}
	}

	return nil
}

func (m FakeDatabaseUtil) DeleteMany(database string, table string, matches map[string]interface{}) error {
	value := map[string]interface{}{}

	toDelete := []int{}

	for k, v := range m.Client[database][table] {
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
		m.Client[database][table] = append(m.Client[database][table][:i], m.Client[database][table][i+1:]...)
	}

	return nil
}
