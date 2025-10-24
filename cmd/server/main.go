package main

import (
	"fmt"
	"log"
	"os"
	"time"

	// "os"

	"term-service/pkg/config"
	"term-service/pkg/consul"
	"term-service/pkg/db"
	"term-service/pkg/router"

	"term-service/pkg/zap"

	consulapi "github.com/hashicorp/consul/api"
)

func main() {
	filePath := os.Args[1]
	if filePath == "" {
		filePath = "configs/config.yaml"
	}

	config.LoadConfig(filePath)

	cfg := config.AppConfig

	// logger.WriteLogData("info", map[string]any{"id": 123, "name": "Hung"})

	//logger
	logger, err := zap.New(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	//consul
	consulConn := consul.NewConsulConn(logger, cfg)
	consulClient := consulConn.Connect()
	defer consulConn.Deregister()

	if err := waitPassing(consulClient, "go-main-service", 60*time.Second); err != nil {
		logger.Fatalf("Dependency not ready: %v", err)
	}

	//db
	db.ConnectMongoDB()

	// redis cache
	cacheClientRedis := db.InitRedisCache()
	defer cacheClientRedis.Close()

	r := router.SetupRouter(consulClient, cacheClientRedis, db.TermCollection, db.HolidayCollection)
	port := cfg.Server.Port
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to run server:", err)
	}
}

func waitPassing(cli *consulapi.Client, name string, timeout time.Duration) error {
	dl := time.Now().Add(timeout)
	for time.Now().Before(dl) {
		entries, _, err := cli.Health().Service(name, "", true, nil)
		if err == nil && len(entries) > 0 {
			return nil
		}
		time.Sleep(2 * time.Second)
	}
	return fmt.Errorf("%s not ready in consul", name)
}
