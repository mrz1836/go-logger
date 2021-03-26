package logger

import "time"

// Package constants
const (
	LogEntriesPort = "443"                             // 80, 514, 443, 10000
	LogEntriesURL  = "us.data.logs.insight.rapid7.com" // "data.logentries.com" "data.insight.rapid7.com (US)"  "eu.data.logs.insight.rapid7.com"
	MaxRetryDelay  = 2 * time.Minute
	RetryDelay     = 100 * time.Millisecond
)
