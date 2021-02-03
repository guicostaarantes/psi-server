package database

import "context"

// IDatabaseUtil is an abstraction for a utility that queries and mutates data in a database
type IDatabaseUtil interface {
	FindOne(database string, table string, matches map[string]interface{}, receiver interface{}) error
	FindMany(database string, table string, matches map[string]interface{}) (ICursor, error)
	InsertOne(database string, table string, provider interface{}) error
	UpdateOne(database string, table string, matches map[string]interface{}, provider interface{}) error
	DeleteOne(database string, table string, matches map[string]interface{}) error
}

// ICursor is an abstraction of an entity for navigating multiple results of a database query
type ICursor interface {
	Decode(receiver interface{}) error
	Next(ctx context.Context) bool
	Close(ctx context.Context) error
}
