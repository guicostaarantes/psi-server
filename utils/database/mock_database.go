package database

import (
	"context"
	"encoding/json"
)

func FakeDBClientFactory() map[string][][]byte {
	store := map[string][][]byte{}

	return store
}

type FakeDatabaseUtil struct {
	Client map[string][][]byte
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
	result := map[string][]string{}

	for tblName, tbl := range m.Client {
		if result[tblName] == nil {
			result[tblName] = []string{}
		}
		for _, fld := range tbl {
			result[tblName] = append(result[tblName], string(fld))
		}
	}

	return json.Marshal(result)
}

func (m FakeDatabaseUtil) SetMockedDatabases(data []byte) error {
	newClient := map[string][]string{}

	jsonErr := json.Unmarshal(data, &newClient)
	if jsonErr != nil {
		return jsonErr
	}

	for tblName, tbl := range m.Client {
		if m.Client[tblName] == nil {
			m.Client[tblName] = [][]byte{}
		}
		for _, fld := range tbl {
			m.Client[tblName] = append(m.Client[tblName], []byte(fld))
		}
	}

	return nil
}

func (m FakeDatabaseUtil) Connect(uri string) error {
	return nil
}

func (m FakeDatabaseUtil) FindOne(table string, matches map[string]interface{}, receiver interface{}) error {
	value := map[string]interface{}{}
	for _, v := range m.Client[table] {
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

func (m FakeDatabaseUtil) FindMany(table string, matches map[string]interface{}) (ICursor, error) {
	cursor := cursorStruct{
		results: [][]byte{},
		current: -1,
	}

	value := map[string]interface{}{}
	for _, v := range m.Client[table] {
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

func (m FakeDatabaseUtil) InsertOne(table string, provider interface{}) error {
	if m.Client == nil {
		m.Client = make(map[string][][]byte)
	}

	if m.Client[table] == nil {
		m.Client[table] = [][]byte{}
	}

	value, marshalErr := json.Marshal(provider)
	if marshalErr != nil {
		return marshalErr
	}

	newTable := append(m.Client[table], value)

	m.Client[table] = newTable

	return nil
}

func (m FakeDatabaseUtil) InsertMany(table string, provider []interface{}) error {
	if m.Client == nil {
		m.Client = make(map[string][][]byte)
	}

	if m.Client[table] == nil {
		m.Client[table] = [][]byte{}
	}

	for _, prov := range provider {
		value, marshalErr := json.Marshal(prov)
		if marshalErr != nil {
			return marshalErr
		}

		newTable := append(m.Client[table], value)

		m.Client[table] = newTable
	}

	return nil
}

func (m FakeDatabaseUtil) UpdateOne(table string, matches map[string]interface{}, provider interface{}) error {
	value := map[string]interface{}{}
	for k, v := range m.Client[table] {
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
			m.Client[table][k] = value
			return nil
		}
	}

	return nil
}

func (m FakeDatabaseUtil) DeleteOne(table string, matches map[string]interface{}) error {
	value := map[string]interface{}{}
	for k, v := range m.Client[table] {
		json.Unmarshal(v, &value)
		matching := true
		for matchKey, matchValue := range matches {
			if value[matchKey] != matchValue {
				matching = false
			}
		}
		if matching {
			m.Client[table] = append(m.Client[table][:k], m.Client[table][k+1:]...)
			return nil
		}
	}

	return nil
}

func (m FakeDatabaseUtil) DeleteMany(table string, matches map[string]interface{}) error {
	value := map[string]interface{}{}

	toDelete := []int{}

	for k, v := range m.Client[table] {
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
		m.Client[table] = append(m.Client[table][:i], m.Client[table][i+1:]...)
	}

	return nil
}
