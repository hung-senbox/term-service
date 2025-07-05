package main

import (
	"log"
	// "os"
	"term-service/logger"
	"term-service/pkg/config"
	"term-service/pkg/consul"
	"term-service/pkg/db"
	"term-service/pkg/router"

	"term-service/pkg/zap"
)

func main() {
	// filePath := os.Args[1]
	// if filePath == "" {
	// 	filePath = "configs/config.yaml"
	// }

	config.LoadConfig("configs/config.yaml")

	cfg := config.AppConfig

	logger.WriteLogData("info", map[string]any{"id": 123, "name": "Hung"})

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
