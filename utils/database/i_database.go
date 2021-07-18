package database

import "context"

// IDatabaseUtil is an abstraction for a utility that queries and mutates data in a database
type IDatabaseUtil interface {
	Connect(uri string) error
	FindOne(table string, matches map[string]interface{}, receiver interface{}) error
	FindMany(table string, matches map[string]interface{}) (ICursor, error)
	InsertOne(table string, provider interface{}) error
	InsertMany(table string, provider []interface{}) error
	UpdateOne(table string, matches map[string]interface{}, provider interface{}) error
	DeleteOne(table string, matches map[string]interface{}) error
	DeleteMany(table string, matches map[string]interface{}) error
}

// ICursor is an abstraction of an entity for navigating multiple results of a database query
type ICursor interface {
	Decode(receiver interface{}) error
	Next(ctx context.Context) bool
	Close(ctx context.Context) error
}
