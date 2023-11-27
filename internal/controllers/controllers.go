package controllers

import (
	// api "github.com/z416352/Crawler/Crawler/pkg/apiservice"
	// "github.com/z416352/Crawler/Crawler/pkg/logger"
	"context"
	"fmt"
	"net/http"
	"strconv"

	api "github.com/z416352/Crawler/pkg/apiservice"
	"github.com/z416352/Crawler/pkg/logger"

	"github.com/z416352/Database-api/internal/db_services"
	db_svc "github.com/z416352/Database-api/internal/db_services"
	"github.com/z416352/Database-api/pkg/responses"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var target *api.BinanceCrawlTarget
var client *mongo.Client

// database_map[<Symbol>]
var database_map map[string]*mongo.Database

// collection_map[<Symbol>][<Collection>]
var collection_map map[string]map[string]*mongo.Collection

func init() {
	target = new(api.BinanceCrawlTarget).GetCrawlTarget()

	client = db_services.ConnectDB()
	database_map = db_svc.GetDatabase(client, target.Symbol_list)
	collection_map = db_svc.GetCollection(
		database_map,
		db_svc.Transfer_TF_Coll(target.TimeFrame_list).([]string),
	)
}

func InsertData() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request api.BinanceTypeRequestDetail

		// Call BindJSON to bind the received JSON
		if err := c.BindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, responses.UserResponse{
				Status:  http.StatusBadRequest,
				Message: err.Error(),
			})
			return
		}

		klines := request.Klines
		symbol := request.Symbol
		collection := db_svc.Transfer_TF_Coll(request.Timeframe).(string)

		result, err := db_svc.InsertCrawlerData(klines, collection_map[symbol][collection])
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.UserResponse{
				Status:  http.StatusInternalServerError,
				Message: err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, responses.UserResponse{
			Status:  http.StatusOK,
			Message: "success",
			Data:    map[string]interface{}{"data": result},
		})
	}
}

func GetData() gin.HandlerFunc {
	return func(c *gin.Context) {
		symbol := c.Param("symbol")
		timeframe := c.Param("timeframe")
		num, _ := strconv.Atoi(c.Param("num"))

		var kline *api.BinanceAPI_Kline
		var klines []*api.BinanceAPI_Kline
		var err error
		if num == 1 {
			kline, err = db_svc.GetOneNewestData(collection_map[symbol][db_svc.Transfer_TF_Coll(timeframe).(string)])
		} else {
			klines, err = db_svc.GetMultiNewestData(collection_map[symbol][db_svc.Transfer_TF_Coll(timeframe).(string)], num)
		}

		if err != nil {
			logger.DBLog.Errorf("Err: %v", err)

			c.JSON(http.StatusInternalServerError, responses.UserResponse{
				Status:  http.StatusInternalServerError,
				Message: "error",
				Data:    map[string]interface{}{"err": err.Error()},
			})
			return
		}

		if num == 1 {
			klines = append(klines, kline)
		}

		c.JSON(http.StatusOK, responses.UserResponse{
			Status:  http.StatusOK,
			Message: fmt.Sprintf("success"),
			Data:    map[string]interface{}{"klines": klines},
		})

	}
}

func GetDBExist() gin.HandlerFunc {
	return func(c *gin.Context) {
		dbname := c.Param("dbname")

		if err := db_svc.CheckDBExist(client, dbname); err != nil {
			c.JSON(http.StatusNotFound, responses.UserResponse{
				Status:  http.StatusNotFound,
				Message: fmt.Sprint(err),
			})
			return
		}

		collections, err := database_map[dbname].ListCollectionNames(context.TODO(), bson.M{})
		timeframes := db_svc.Transfer_TF_Coll(collections).([]string)

		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.UserResponse{
				Status:  http.StatusInternalServerError,
				Message: fmt.Sprint(err),
			})
			return
		}

		c.JSON(http.StatusOK, responses.UserResponse{
			Status:  http.StatusOK,
			Message: fmt.Sprintf("success"),
			Data:    map[string]interface{}{"timeframes": timeframes},
		})
	}
}
