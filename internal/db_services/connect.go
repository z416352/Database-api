package db_services

import (
	"Crawler/pkg/logger"
	"context"
	"fmt"
	"net"
	"os"
)

type MongodbInfoK8s struct {
	ServiceName string
	Namespace   string
	K8sDNSIp    string
}

type ConnectData struct {
	ip           string
	userName     string
	userPassword string
	port         string
	MongodbInfo  *MongodbInfoK8s
}

func (c *ConnectData) SetConnectData() *ConnectData {
	// c.ip = os.Getenv("IP")
	c.ip = c.getServiceIP()
	c.port = "27017"
	c.userName = os.Getenv("Name")
	c.userPassword = os.Getenv("Pass")

	if c.port == "" || c.userName == "" || c.userPassword == "" {
		logger.DBLog.Panicf("Set parameters. Enter in command line: Name=\"<UserName>\" Pass=\"<Password>\"")
	}

	return c
}

func (c *ConnectData) GetConnectURI() string {
	return fmt.Sprintf("mongodb://%s:%s@%s:%s/", c.userName, c.userPassword, c.ip, c.port)
}

func (c *ConnectData) getServiceIP() string {
	dnsName := c.MongodbInfo.ServiceName + "." + c.MongodbInfo.Namespace + ".svc.cluster.local"

	// 使用指定的 DNS 服务器地址进行查询
	resolver := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{}
			return d.DialContext(ctx, "udp", c.MongodbInfo.K8sDNSIp+":53")
		},
	}

	ips, err := resolver.LookupIPAddr(context.Background(), dnsName)
	if err != nil {
		fmt.Printf("无法解析 DNS 名称：%v\n", err)
		os.Exit(1)
	}

	// fmt.Printf("服务 %s 的 IP 地址列表：\n", c.MongodbInfo.ServiceName)

	return ips[0].String()
}
