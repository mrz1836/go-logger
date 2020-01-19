package logger

import "time"

// Package constants
const (
	LogEntriesPort = "10000"
	LogEntriesURL  = "data.logentries.com"
	//LogEntriesURL = "us.data.logs.insight.rapid7.com"
	MaxRetryDelay = 2 * time.Minute
	RetryDelay    = 100 * time.Millisecond
)
