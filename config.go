package logger

import "time"

// More info: https://docs.rapid7.com/insightidr/ports-used-by-insightidr/

// Package constants
const (
	LogEntriesPort = "10000"                   // 80, 514, 443, 10000
	LogEntriesURL  = "data.insight.rapid7.com" // "data.logentries.com" "data.insight.rapid7.com"  "eu.data.logs.insight.rapid7.com"
	MaxRetryDelay  = 2 * time.Minute
	RetryDelay     = 100 * time.Millisecond
)
