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

	// Last, use traditional Fatalf
	logger.Fatalf("fatal print line via %s - goodbye!", "Fatalf")
}
