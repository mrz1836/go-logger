/*
Package main is an example package to show the use case of go-logger
*/
package main

import "github.com/mrz1836/go-logger"

// main
func main() {
	// Log using the Data() method
	logger.Data(2, logger.DEBUG, "testing the go-logger package")

	// Use traditional Println
	logger.Println("regular print line")

	// Use traditional Printf
	logger.Printf("print line via %s", "Printf")

	// Use traditional Errorln
	logger.Errorln(1, "error print line")

	// Use traditional Errorfmt
	logger.Errorfmt(1, "print error line via %s", "Errorfmt")

	/*
		// Create another logger instance
		client, err := logger.NewLogEntriesClient("token-1234567", "us.data.logs.insight.rapid7.com", logger.LogEntriesPort)
		if err != nil {
			logger.Fatalf("error creating client %s", err.Error())
			return
		}

		// Start processing queue
		go client.ProcessQueue()

		// Fire a log using the new client
		client.Println("message for new log client")
	*/

	// Last, use traditional Fatalf
	logger.Fatalf("fatal print line via %s - goodbye!", "Fatalf")
}
