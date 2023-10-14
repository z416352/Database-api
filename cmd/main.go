package main

import (
	// db_svc "Database/internal/db_services"

	gin_svc "Database/internal/gin_services"
)

func main() {

	go gin_svc.Gin_handler("8080")

	select {}

}
