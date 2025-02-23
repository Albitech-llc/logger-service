package main

import (
	"fmt"
	"time"

	"github.com/Elbitech-llc/logger-service/logger"
	"github.com/Elbitech-llc/logger-service/pkg/caching"
)

func main() {
	logServ := logger.NewService()

	cfg := logger.LoadConfig()

	rdb, _, err := caching.InitializeRedis(cfg.Redis.Host, cfg.Redis.Port, cfg.Redis.LogsDB)
	if err != nil {
		fmt.Printf("Failed to initialize Redis: %v", err)
	}
	defer func() {
		if err := rdb.Close(); err != nil {
			fmt.Printf("Failed to close Redis: %v", err)
		}
	}()

	defer logServ.Close(rdb) // Ensure service resources are released

	fmt.Println("Start")
	for i := range 2000 {
		logServ.LogInfo("Main", fmt.Sprintf("INFO %d", i))

		logServ.LogWarning("Main", fmt.Sprintf("WWWWWW %d", i))

		logServ.LogError("Main", fmt.Sprintf("EEEEE %d", i))
	}

	time.Sleep(15 * time.Second)
	fmt.Println("Step 2")
	for i := range 3000 {
		logServ.LogInfo("Main 2", fmt.Sprintf("INFO %d", i))

		logServ.LogWarning("Main 2", fmt.Sprintf("WWWWWW %d", i))

		logServ.LogError("Main 2", fmt.Sprintf("EEEEE %d", i))
	}

	fmt.Println("End")
	time.Sleep(15 * time.Second)

}
