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
}

func (c *ConnectData) SetConnectData() {
	if _, ok := os.LookupEnv("TEST"); ok {
		return
	}

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
}

func (c *ConnectData) GetConnectURI() string {
	if test_mod := os.Getenv("TEST"); test_mod == "1" {
		if uri, ok := os.LookupEnv("mongodb_URI"); ok {
			return uri
		}else{
			logger.DBLog.Panicf("Test mod. 'mongodb_URI' pram missing.")
		}

	}

	return fmt.Sprintf("mongodb+srv://%s:%s@%s/?retryWrites=true&w=majority", c.userName, c.userPassword, c.mongoDB_Atlas_Info)
}
