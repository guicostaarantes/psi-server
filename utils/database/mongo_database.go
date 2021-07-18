package database

import (
	"context"
	"errors"

	"github.com/guicostaarantes/psi-server/utils/logging"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDatabaseUtil struct {
	Context     context.Context
	Client      *mongo.Client
	LoggingUtil logging.ILoggingUtil
}

func (m *MongoDatabaseUtil) Connect(uri string) error {
	client, connectErr := mongo.NewClient(options.Client().ApplyURI(uri))
	if connectErr != nil {
		return connectErr
	}

	connectErr = client.Connect(context.Background())
	if connectErr != nil {
		return connectErr
	}

	m.Client = client

	return nil
}

func (m *MongoDatabaseUtil) FindOne(database string, table string, matches map[string]interface{}, receiver interface{}) error {
	collection := m.Client.Database(database).Collection(table)

	filter := bson.D{}
	for k, v := range matches {
		filter = append(filter, primitive.E{Key: k, Value: v})
	}

	mongoErr := collection.FindOne(m.Context, filter).Decode(receiver)

	if mongoErr != nil && mongoErr.Error() != "mongo: no documents in result" {
		m.LoggingUtil.Error("d98dc14d", mongoErr)
		return errors.New("internal server error")
	}

	return nil
}

func (m *MongoDatabaseUtil) FindMany(database string, table string, matches map[string]interface{}) (ICursor, error) {
	collection := m.Client.Database(database).Collection(table)

	filter := bson.D{}
	for k, v := range matches {
		filter = append(filter, primitive.E{Key: k, Value: v})
	}

	cursor, mongoErr := collection.Find(m.Context, filter)

	if mongoErr != nil && mongoErr.Error() != "mongo: no documents in result" {
		m.LoggingUtil.Error("d98dc14d", mongoErr)
		return nil, errors.New("internal server error")
	}

	return cursor, nil
}

func (m *MongoDatabaseUtil) InsertOne(database string, table string, provider interface{}) error {
	collection := m.Client.Database(database).Collection(table)

	_, insertErr := collection.InsertOne(m.Context, provider)

	if insertErr != nil {
		m.LoggingUtil.Error("018778ce", insertErr)
		return errors.New("internal server error")
	}

	return nil
}

func (m *MongoDatabaseUtil) InsertMany(database string, table string, provider []interface{}) error {
	collection := m.Client.Database(database).Collection(table)

	if len(provider) == 0 {
		return nil
	}

	_, insertErr := collection.InsertMany(m.Context, provider)

	if insertErr != nil {
		m.LoggingUtil.Error("3789b465", insertErr)
		return errors.New("internal server error")
	}

	return nil
}

func (m *MongoDatabaseUtil) UpdateOne(database string, table string, matches map[string]interface{}, provider interface{}) error {
	collection := m.Client.Database(database).Collection(table)

	filter := bson.D{}
	for k, v := range matches {
		filter = append(filter, primitive.E{Key: k, Value: v})
	}

	update := bson.M{
		"$set": provider,
	}

	_, updateErr := collection.UpdateOne(m.Context, filter, update)

	if updateErr != nil {
		m.LoggingUtil.Error("cf9a49be", updateErr)
		return errors.New("internal server error")
	}

	return nil
}

func (m *MongoDatabaseUtil) DeleteOne(database string, table string, matches map[string]interface{}) error {
	collection := m.Client.Database(database).Collection(table)

	filter := bson.D{}
	for k, v := range matches {
		filter = append(filter, primitive.E{Key: k, Value: v})
	}

	_, deleteErr := collection.DeleteOne(m.Context, filter)

	if deleteErr != nil {
		m.LoggingUtil.Error("c6280e5a", deleteErr)
		return errors.New("internal server error")
	}

	return nil
}

func (m *MongoDatabaseUtil) DeleteMany(database string, table string, matches map[string]interface{}) error {
	collection := m.Client.Database(database).Collection(table)

	filter := bson.D{}
	for k, v := range matches {
		filter = append(filter, primitive.E{Key: k, Value: v})
	}

	_, deleteErr := collection.DeleteMany(m.Context, filter)

	if deleteErr != nil {
		m.LoggingUtil.Error("d802d913", deleteErr)
		return errors.New("internal server error")
	}

	return nil
}
