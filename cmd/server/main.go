package main

import (
	"log"
	"os"
	"term-info-service/pkg/config"
	"term-info-service/pkg/consul"
	"term-info-service/pkg/db"
	"term-info-service/pkg/router"

	"term-info-service/pkg/zap"
)

func main() {
	filePath := os.Args[1]
	if filePath == "" {
		filePath = "configs/config.yaml"
	}

	config.LoadConfig(filePath)

	cfg := config.AppConfig

	//logger
	logger, err := zap.New(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	//consul
	consulConn := consul.NewConsulConn(logger, cfg)
	consulConn.Connect()
	defer consulConn.Deregister()

	//db
	db.ConnectMongoDB()

	r := router.SetupRouter(db.TermCollection)
	port := cfg.Server.Port
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to run server:", err)
	}
}
