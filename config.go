package logger

import "time"

// Package constants
const (
	LogEntriesPort = "10000"                           // 80, 514, 10000
	LogEntriesURL  = "us.data.logs.insight.rapid7.com" // "data.logentries.com"  "eu.data.logs.insight.rapid7.com"
	MaxRetryDelay  = 2 * time.Minute
	RetryDelay     = 100 * time.Millisecond
)
