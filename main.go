package main

import (
	"fmt"
	"time"

	"github.com/Albitech-llc/logger-service/logger"
)

func main() {
	logServ := logger.NewService()
	defer logServ.Close()

	fmt.Println("Start")
	for i := range 10 {
		logServ.LogInfo("Main", fmt.Sprintf("INFO %d", i))

		logServ.LogWarning("Main", fmt.Sprintf("WWWWWW %d", i))

		logServ.LogError("Main", fmt.Sprintf("EEEEE %d", i))
	}

	time.Sleep(15 * time.Second)
	fmt.Println("Step 2")
	for i := range 3 {
		logServ.LogInfo("Main 2", fmt.Sprintf("INFO %d", i))

		logServ.LogWarning("Main 2", fmt.Sprintf("WWWWWW %d", i))

		logServ.LogError("Main 2", fmt.Sprintf("EEEEE %d", i))
	}

	fmt.Println("End")
	time.Sleep(15 * time.Second)

}
