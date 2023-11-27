package db_services

import (
	"context"
	"errors"
	"time"

	api "github.com/z416352/Crawler/pkg/apiservice"
	"github.com/z416352/Crawler/pkg/logger"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InsertCrawlerData(klines []*api.BinanceAPI_Kline, collection *mongo.Collection) (*mongo.InsertManyResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	newValue := make([]interface{}, len(klines))
	defer cancel()

	for i := range klines {
		newValue[i] = klines[i]
	}

	result, err := collection.InsertMany(ctx, newValue)
	if err != nil {
		logger.DBLog.Error(err)
	}
	logger.DBLog.Debugf("Inserted '%d' K-lines data:", len(klines))

	return result, err
}

func GetAllData(collection *mongo.Collection) ([]*api.BinanceAPI_Kline, error) {
	var klines []*api.BinanceAPI_Kline

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.D{{}}

	// Perform the find operation
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		logger.DBLog.Errorf("Err: %v", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	// Iterate through the documents
	for cursor.Next(ctx) {
		var kline *api.BinanceAPI_Kline

		// Decode the document into the Person struct
		err := cursor.Decode(&kline)
		if err != nil {
			logger.DBLog.Errorf("Err: %v", err)
			return nil, err
		}

		klines = append(klines, kline)
	}

	return klines, nil
}

func GetOneNewestData(collection *mongo.Collection) (*api.BinanceAPI_Kline, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.D{{}}
	opts := options.FindOne().SetSort(bson.D{{"_id", -1}}) // Sort in descending order based on "_id" field

	// Perform the find operation
	var kline *api.BinanceAPI_Kline
	err := collection.FindOne(ctx, filter, opts).Decode(&kline)

	if err == mongo.ErrNoDocuments {
		return nil, errors.New("no documents found")
	} else if err != nil {
		return nil, err
	}

	return kline, nil
}

func GetMultiNewestData(collection *mongo.Collection, n int) ([]*api.BinanceAPI_Kline, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.D{{}}

	// Sort in descending order based on "_id" field and limit to n results
	opts := options.Find().SetSort(bson.D{{"_id", -1}}).SetLimit(int64(n))

	// Perform the find operation
	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var klines []*api.BinanceAPI_Kline

	for cursor.Next(ctx) {
		var kline *api.BinanceAPI_Kline
		if err := cursor.Decode(&kline); err != nil {
			return nil, err
		}
		klines = append(klines, kline)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return klines, nil
}
