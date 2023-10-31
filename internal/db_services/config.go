package db_services

import (
	"context"
	"fmt"
	"slices"

	"github.com/z416352/Crawler/pkg/logger"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDB() *mongo.Client {
	c := new(ConnectData)
	c.MongodbInfo = &MongodbInfoK8s{
		ServiceName: "mongodb-service",
		Namespace:   "mongodb",
		K8sDNSIp:    "10.96.0.10",
	}

	c.SetConnectData()

	// 设置连接选项
	clientOptions := options.Client().ApplyURI(c.GetConnectURI())

	// 连接到 MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		logger.DBLog.Panic(err)
	}

	// 检查连接
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		logger.DBLog.Panic(err)
	}

	logger.DBLog.Infof("You successfully connected to MongoDB!")

	return client
}

func CheckDBExist(client *mongo.Client, databaseNames ...string) error {
	// List all databases in MongoDB.
	databaseList, err := client.ListDatabaseNames(context.TODO(), bson.M{})
	if err != nil {
		return err
	}

	for _, dbName := range databaseNames {
		if slices.Contains(databaseList, dbName) {
			continue
		}
		return fmt.Errorf("Database '%s' does not exist.", dbName)
	}

	return nil
}

func CheckCollectionExist(db *mongo.Database, collectionNames ...string) error {
	// List collections
	collectionList, err := db.ListCollectionNames(context.TODO(), bson.M{})
	if err != nil {
		return err
	}

	for _, collection := range collectionNames {
		if slices.Contains(collectionList, collection) {
			continue
		}
		return fmt.Errorf("Collection '%s' does not exist.", collection)
	}

	return nil
}

func GetDatabase(client *mongo.Client, dbNames []string) map[string]*mongo.Database {
	db_clients := make(map[string]*mongo.Database)

	for _, dbName := range dbNames {
		db_clients[dbName] = client.Database(dbName)
	}

	return db_clients
}

// Getting database collections
func GetCollection(database map[string]*mongo.Database, collection_list []string) map[string]map[string]*mongo.Collection {
	// collection_map[<Symbol>][<Collection>]
	var collection = map[string]map[string]*mongo.Collection{}

	// Check colcollection exist
	for symbol, db_client := range database {
		collection[symbol] = map[string]*mongo.Collection{}

		for _, collectionName := range collection_list {
			collection[symbol][collectionName] = db_client.Collection(collectionName)
		}
	}

	for symbol, timeframe_map := range collection {
		coll_list := []string{}
		for timeframe, _ := range timeframe_map {
			coll_list = append(coll_list, timeframe)
		}
		logger.DBLog.Infof("Connected  ->  Database: '%s', Collection: '%v'", symbol, coll_list)
	}

	logger.DBLog.Infof("Connect %d Database.", len(collection))

	return collection
}
