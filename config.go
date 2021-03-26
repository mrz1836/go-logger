package logger

import "time"

// More info: https://docs.rapid7.com/insightidr/ports-used-by-insightidr/

// Package constants
const (
	LogEntriesPort         = "10000"               // 80, 514, 443, 10000
	LogEntriesTestEndpoint = "34.253.67.177"       // This is an IP for now, since Github Actions fails on resolving the domains
	LogEntriesURL          = "data.logentries.com" // "data.insight.rapid7.com" "eu.data.logs.insight.rapid7.com"
	MaxRetryDelay          = 2 * time.Minute
	RetryDelay             = 100 * time.Millisecond
)
