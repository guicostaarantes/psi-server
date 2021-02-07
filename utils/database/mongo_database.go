package database

import (
	"context"
	"errors"
	"log"
	"os"
	"strings"

	"github.com/guicostaarantes/psi-server/utils/logging"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func connect() *mongo.Client {
	uri := []string{
		"mongodb://",
		os.Getenv("PSI_MONGO_USERNAME"),
		":",
		os.Getenv("PSI_MONGO_PASSWORD"),
		"@",
		os.Getenv("PSI_MONGO_HOST"),
		":",
		os.Getenv("PSI_MONGO_PORT"),
	}

	client, err := mongo.NewClient(options.Client().ApplyURI(strings.Join(uri, "")))
	if err != nil {
		log.Fatal(err)
	}

	err = client.Connect(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	return client
}

type mongoClient struct {
	context          context.Context
	client           *mongo.Client
	loggingUtil      logging.ILoggingUtil
	noDocumentsError string
}

func (m mongoClient) FindOne(database string, table string, matches map[string]interface{}, receiver interface{}) error {
	collection := m.client.Database(database).Collection(table)

	filter := bson.D{}
	for k, v := range matches {
		filter = append(filter, primitive.E{Key: k, Value: v})
	}

	mongoErr := collection.FindOne(m.context, filter).Decode(receiver)

	if mongoErr != nil && mongoErr.Error() != m.noDocumentsError {
		m.loggingUtil.Error("d98dc14d", mongoErr)
		return errors.New("internal server error")
	}

	return nil
}

func (m mongoClient) FindMany(database string, table string, matches map[string]interface{}) (ICursor, error) {
	collection := m.client.Database(database).Collection(table)

	filter := bson.D{}
	for k, v := range matches {
		filter = append(filter, primitive.E{Key: k, Value: v})
	}

	cursor, mongoErr := collection.Find(m.context, filter)

	if mongoErr != nil && mongoErr.Error() != m.noDocumentsError {
		m.loggingUtil.Error("d98dc14d", mongoErr)
		return nil, errors.New("internal server error")
	}

	return cursor, nil
}

func (m mongoClient) InsertOne(database string, table string, provider interface{}) error {
	collection := m.client.Database(database).Collection(table)

	_, insertErr := collection.InsertOne(m.context, provider)

	if insertErr != nil {
		m.loggingUtil.Error("018778ce", insertErr)
		return errors.New("internal server error")
	}

	return nil
}

func (m mongoClient) InsertMany(database string, table string, provider []interface{}) error {
	collection := m.client.Database(database).Collection(table)

	_, insertErr := collection.InsertMany(m.context, provider)

	if insertErr != nil {
		m.loggingUtil.Error("3789b465", insertErr)
		return errors.New("internal server error")
	}

	return nil
}

func (m mongoClient) UpdateOne(database string, table string, matches map[string]interface{}, provider interface{}) error {
	collection := m.client.Database(database).Collection(table)

	filter := bson.D{}
	for k, v := range matches {
		filter = append(filter, primitive.E{Key: k, Value: v})
	}

	update := bson.M{
		"$set": provider,
	}

	_, updateErr := collection.UpdateOne(m.context, filter, update)

	if updateErr != nil {
		m.loggingUtil.Error("cf9a49be", updateErr)
		return errors.New("internal server error")
	}

	return nil
}

func (m mongoClient) DeleteOne(database string, table string, matches map[string]interface{}) error {
	collection := m.client.Database(database).Collection(table)

	filter := bson.D{}
	for k, v := range matches {
		filter = append(filter, primitive.E{Key: k, Value: v})
	}

	_, deleteErr := collection.DeleteOne(m.context, filter)

	if deleteErr != nil {
		m.loggingUtil.Error("c6280e5a", deleteErr)
		return errors.New("internal server error")
	}

	return nil
}

func (m mongoClient) DeleteMany(database string, table string, matches map[string]interface{}) error {
	collection := m.client.Database(database).Collection(table)

	filter := bson.D{}
	for k, v := range matches {
		filter = append(filter, primitive.E{Key: k, Value: v})
	}

	_, deleteErr := collection.DeleteMany(m.context, filter)

	if deleteErr != nil {
		m.loggingUtil.Error("d802d913", deleteErr)
		return errors.New("internal server error")
	}

	return nil
}

// MongoDatabaseUtil is an implementation of IQueryUtil that uses MongoDB via go.mongodb.org/mongo-driver/mongo
var MongoDatabaseUtil = mongoClient{
	context:          context.Background(),
	client:           connect(),
	loggingUtil:      logging.PrintLogUtil,
	noDocumentsError: "mongo: no documents in result",
}
