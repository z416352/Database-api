package db_services

import (
	"context"
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
	if err != nil {
		return nil, err
	}

	return kline, nil
}
