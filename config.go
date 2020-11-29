package logger

import "time"

// Package constants
const (
	LogEntriesPort = "10000"               // 80, 514, 10000
	LogEntriesURL  = "data.logentries.com" // "us.data.logs.insight.rapid7.com"  "eu.data.logs.insight.rapid7.com"
	MaxRetryDelay  = 2 * time.Minute
	RetryDelay     = 100 * time.Millisecond
)
