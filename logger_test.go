package logger

import (
	"bytes"
	"io"
	"log"
	"os"
	"strings"
	"sync"
	"testing"
)

// captureOutput captures the output of log, fmt or os.Stderr.WriteString
func captureOutput(f func()) string {
	reader, writer, err := os.Pipe()
	if err != nil {
		panic(err)
	}
	stdout := os.Stdout
	stderr := os.Stderr
	defer func() {
		os.Stdout = stdout
		os.Stderr = stderr
		log.SetOutput(os.Stderr)
	}()
	os.Stdout = writer
	os.Stderr = writer
	log.SetOutput(writer)
	out := make(chan string)
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		var buf bytes.Buffer
		wg.Done()
		_, _ = io.Copy(&buf, reader)
		out <- buf.String()
	}()
	wg.Wait()
	f()
	_ = writer.Close()
	return <-out
}

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

// TestPrintln test the print line method
func TestPrintln(t *testing.T) {
	captured := captureOutput(func() {
		Errorln(2, "test this method")
	})

	if !strings.Contains(captured, "go-logger/logger_test.go:go-logger.TestPrintln") {
		t.Fatalf("expected string: %s got: %s", "go-logger/logger_test.go:go-logger.TestPrintln", captured)
	}

	if !strings.Contains(captured, "test this method") {
		t.Fatalf("expected string: %s got: %s", "test this method", captured)
	}
}

// TestPrintf test the print fmt method
func TestPrintf(t *testing.T) {
	captured := captureOutput(func() {
		Errorfmt(2, "test this method: %s", "TestPrintf")
	})

	if !strings.Contains(captured, "go-logger/logger_test.go:go-logger.TestPrintf") {
		t.Fatalf("expected string: %s got: %s", "go-logger/logger_test.go:go-logger.TestPrintf", captured)
	}

	if !strings.Contains(captured, "test this method: TestPrintf") {
		t.Fatalf("expected string: %s got: %s", "test this method: TestPrintf", captured)
	}
}

// TestErrorln test the error line method
func TestErrorln(t *testing.T) {
	captured := captureOutput(func() {
		Errorln(2, "test this method")
	})

	if !strings.Contains(captured, "go-logger/logger_test.go:go-logger.TestErrorln") {
		t.Fatalf("expected string: %s got: %s", "go-logger/logger_test.go:go-logger.TestErrorln", captured)
	}

	if !strings.Contains(captured, "test this method") {
		t.Fatalf("expected string: %s got: %s", "test this method", captured)
	}
}

// TestErrorfmt test the error fmt method
func TestErrorfmt(t *testing.T) {
	captured := captureOutput(func() {
		Errorfmt(2, "test this method: %s", "Errorfmt")
	})

	if !strings.Contains(captured, "go-logger/logger_test.go:go-logger.TestErrorfmt") {
		t.Fatalf("expected string: %s got: %s", "go-logger/logger_test.go:go-logger.TestErrorfmt", captured)
	}

	if !strings.Contains(captured, "test this method: Errorfmt") {
		t.Fatalf("expected string: %s got: %s", "test this method: Errorfmt", captured)
	}
}

// TestData test the data method
func TestData(t *testing.T) {
	captured := captureOutput(func() {
		Data(2, WARN, "test this method", MakeError("another", "value"))
	})

	//2019/06/17 12:59:32 type="warn" file="go-logger/logger_test.go" method="go-logger.TestData.func1" line="188" message="test this method" another="value"

	// Check for warn
	if !strings.Contains(captured, `type="warn"`) {
		t.Fatalf("expected string: %s, got: %s", `type="warn"`, captured)
	}

	// Check for file
	if !strings.Contains(captured, `file="go-logger/logger_test.go"`) {
		t.Fatalf("expected string: %s, got: %s", `file="go-logger/logger_test.go"`, captured)
	}

	// Check for method
	if !strings.Contains(captured, `method="go-logger.TestData.func1"`) {
		t.Fatalf("expected string: %s, got: %s", `method="go-logger.TestData.func1"`, captured)
	}

	// Check for message
	if !strings.Contains(captured, `message="test this method"`) {
		t.Fatalf("expected string: %s, got: %s", `message="test this method"`, captured)
	}

	// Check for additional values
	if !strings.Contains(captured, `another="value"`) {
		t.Fatalf("expected string: %s, got: %s", `another="value"`, captured)
	}
}
