package database

import "context"

// IDatabaseUtil is an abstraction for a utility that queries and mutates data in a database
type IDatabaseUtil interface {
	FindOne(database string, table string, field string, match interface{}, receiver interface{}) error
	FindMany(database string, table string, field string, match interface{}) (ICursor, error)
	FindAll(database string, table string) (ICursor, error)
	InsertOne(database string, table string, provider interface{}) error
	UpdateOne(database string, table string, field string, match interface{}, provider interface{}) error
	DeleteOne(database string, table string, field string, match interface{}) error
}

type ICursor interface {
	Decode(receiver interface{}) error
	Next(ctx context.Context) bool
	Close(ctx context.Context) error
}
