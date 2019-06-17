package logger

import (
	"strings"
	"testing"
)

// TestLogLevel_String test the log level to string method
func TestLogLevel_String(t *testing.T) {

	// Set the level
	var level LogLevel = 0

	// Test for debug
	if level.String() != "debug" {
		t.Fatalf("expected string to be: %s, got: %s", "debug", level.String())
	}

	// Test for info
	level = 1
	if level.String() != "info" {
		t.Fatalf("expected string to be: %s, got: %s", "info", level.String())
	}

	// Test for warn
	level = 2
	if level.String() != "warn" {
		t.Fatalf("expected string to be: %s, got: %s", "warn", level.String())
	}

	// Test for error
	level = 3
	if level.String() != "error" {
		t.Fatalf("expected string to be: %s, got: %s", "error", level.String())
	}
}

// TestFileTag test file tag method
func TestFileTag(t *testing.T) {

	// File tag
	fileTag := FileTag(1)
	if !strings.Contains(fileTag, "go-logger/logger_test.go:go-logger.TestFileTag:") {
		t.Fatalf("expected file tag: %s, got: %s", "go-logger/logger_test.go:go-logger.TestFileTag:", fileTag)
	}
}

// TestFileTagComponents test the file tag components method
func TestFileTagComponents(t *testing.T) {

	// Test the level: 2
	fileTagComps := FileTagComponents(2)
	if len(fileTagComps) == 0 || len(fileTagComps) != 3 {
		t.Fatal("expected file tag components to have 3 components")
	}

	// Test the part
	if fileTagComps[0] != "testing/testing.go" {
		t.Fatalf("expected component: %s, got: %s", "testing/testing.go", fileTagComps[0])
	}

	// Test the part
	if fileTagComps[1] != "testing.tRunner" {
		t.Fatalf("expected component: %s, got: %s", "testing.tRunner", fileTagComps[1])
	}

	// Test the part
	if fileTagComps[2] != "865" {
		t.Fatalf("expected component: %s, got: %s", "865", fileTagComps[2])
	}

	// Test the level: 1
	fileTagComps = FileTagComponents(1)
	if len(fileTagComps) == 0 || len(fileTagComps) != 3 {
		t.Fatal("expected file tag components to have 3 components")
	}

	// Test the part
	if fileTagComps[0] != "go-logger/logger_test.go" {
		t.Fatalf("expected component: %s, got: %s", "go-logger/logger_test.go", fileTagComps[0])
	}

	// Test the part
	if fileTagComps[1] != "go-logger.TestFileTagComponents" {
		t.Fatalf("expected component: %s, got: %s", "go-logger.TestFileTagComponents", fileTagComps[1])
	}
}
