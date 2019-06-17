/*
Package main is an example package to show the use case of go-logger
*/
package main

import "github.com/mrz1836/go-logger"

// main
func main() {

	// Log using the Data() method
	logger.Data(2, logger.DEBUG, "testing the go-logger package")
}
