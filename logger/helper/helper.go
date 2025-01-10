package helper

import (
	"fmt"
	"os"
)

// Helper function to write log data to the fallback file
func WriteToFile(file *os.File, log []byte) {
	//timestamp := time.Now().Format(time.RFC3339)
	// logEntry := map[string]interface{}{
	// 	"timestamp": timestamp,
	// 	"log":       json.RawMessage(log),
	// }
	// logData, err := json.Marshal(logEntry)
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "Failed to serialize log for file: %v\n", err)
	// 	return
	// }
	//defer file.Close()
	_, writeErr := file.Write(append(log, '\n')) // Write log and add newline
	if writeErr != nil {
		fmt.Fprintf(os.Stderr, "Failed to write log to file: %v\n", writeErr)
	}
}
