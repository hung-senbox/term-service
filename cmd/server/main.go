package main

import (
	"log"
	"term-info-service/pkg/config"
	"term-info-service/pkg/consul"
	"term-info-service/pkg/db"
	"term-info-service/pkg/router"

	"term-info-service/pkg/zap"
)

func main() {
	config.LoadConfig("configs/config.yaml")

	cfg := config.AppConfig

	//logger
	logger, err := zap.New(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	//consul
	consulConn := consul.NewConsulConn(logger, cfg)
	// consulClient := consulConn.Connect()
	defer consulConn.Deregister()

	//db
	db.ConnectMongoDB()

	r := router.SetupRouter(db.TermCollection)

	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to run server:", err)
	}
}
