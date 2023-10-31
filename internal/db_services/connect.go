package db_services

import (
	"fmt"
	"os"

	"github.com/z416352/Crawler/pkg/logger"
)

type MongodbInfoK8s struct {
	ServiceName string
	Namespace   string
	K8sDNSIp    string
}

type ConnectData struct {
	userName           string
	userPassword       string
	mongoDB_Atlas_Info string
	MongodbInfo        *MongodbInfoK8s
}

func (c *ConnectData) SetConnectData() *ConnectData {
	c.userName = os.Getenv("MONGODB_USERNAME")
	c.userPassword = os.Getenv("MONGODB_PASSWORD")
	c.mongoDB_Atlas_Info = os.Getenv("Atlas_Info")

	if c.userName == "" || c.userPassword == "" {
		params := []string{}

		if c.userName == "" {
			params = append(params, "MONGODB_USERNAME")
		}

		if c.userPassword == "" {
			params = append(params, "MONGODB_PASSWORD")
		}

		if c.mongoDB_Atlas_Info == "" {
			params = append(params, "Atlas_Info")
		}

		logger.DBLog.Panicf("Parameters Error. %v aren't set.", params)
	}

	return c
}

func (c *ConnectData) GetConnectURI() string {
	return fmt.Sprintf("mongodb+srv://%s:%s@%s/?retryWrites=true&w=majority", c.userName, c.userPassword, c.mongoDB_Atlas_Info)
}
