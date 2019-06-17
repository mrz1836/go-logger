package logger

import (
	"bytes"
	"fmt"
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
	var level LogLevel

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

// ExampleLogLevel_String example using level.String()
func ExampleLogLevel_String() {
	var level LogLevel
	fmt.Println(level.String())
	// Output:debug
}

// BenchmarkLogLevel_String benchmarks the level.String() method
func BenchmarkLogLevel_String(b *testing.B) {
	var level LogLevel
	for i := 0; i < b.N; i++ {
		_ = level.String()
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

// ExampleFileTag example using FileTag()
func ExampleFileTag() {
	fileTag := FileTag(1)
	fmt.Println(fileTag)
	// Output:go-logger/logger_test.go:go-logger.ExampleFileTag:102
}

// BenchmarkFileTag benchmarks the FileTag() method
func BenchmarkFileTag(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = FileTag(1)
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

// ExampleFileTagComponents example using FileTagComponents()
func ExampleFileTagComponents() {
	fileTag := FileTagComponents(1)
	fmt.Println(fileTag[0])
	// Output:go-logger/logger_test.go
}

// BenchmarkFileTaComponents benchmarks the FileTagComponents() method
func BenchmarkFileTagComponents(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = FileTagComponents(1)
	}
}

// TestPrintln test the print line method
func TestPrintln(t *testing.T) {
	captured := captureOutput(func() {
		Println("test this method")
	})

	if !strings.Contains(captured, "go-logger/logger_test.go:go-logger.TestPrintln") {
		t.Fatalf("expected string: %s got: %s", "go-logger/logger_test.go:go-logger.TestPrintln", captured)
	}

	if !strings.Contains(captured, "test this method") {
		t.Fatalf("expected string: %s got: %s", "test this method", captured)
	}
}

// BenchmarkPrintln benchmarks the Println() method
func BenchmarkPrintln(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Println("something")
	}
}

// TestPrintf test the print fmt method
func TestPrintf(t *testing.T) {
	captured := captureOutput(func() {
		Printf("test this method: %s", "TestPrintf")
	})

	if !strings.Contains(captured, "go-logger/logger_test.go:go-logger.TestPrintf") {
		t.Fatalf("expected string: %s got: %s", "go-logger/logger_test.go:go-logger.TestPrintf", captured)
	}

	if !strings.Contains(captured, "test this method: TestPrintf") {
		t.Fatalf("expected string: %s got: %s", "test this method: TestPrintf", captured)
	}
}

// BenchmarkPrintf benchmarks the Printf() method
func BenchmarkPrintf(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Printf("test this method: %s", "TestPrintf")
	}
}

// TestErrorln test the Errorln() method
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

// BenchmarkErrorln benchmarks the Errorln() method
func BenchmarkErrorln(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Errorln(2, "test this method")
	}
}

// TestErrorfmt test the Errorfmt() method
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

// BenchmarkErrorfmt benchmarks the Errorfmt() method
func BenchmarkErrorfmt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Errorfmt(2, "test this method: %s", "Errorfmt")
	}
}

// TestData test the Data() method
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

// BenchmarkData benchmarks the Data() method
func BenchmarkData(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Data(2, WARN, "test this method", MakeError("another", "value"))
	}
}

// TestLogPkg_Printf test log package Printf() method
func TestLogPkg_Printf(t *testing.T) {
	implementation = &logPkg{}

	captured := captureOutput(func() {
		implementation.Printf("test this method: %s", "TestPrintf")
	})

	if !strings.Contains(captured, "test this method: TestPrintf") {
		t.Fatalf("expected string: %s got: %s", "test this method: TestPrintf", captured)
	}
}

// BenchmarkLogPkg_Printf benchmarks the LogPkg_Printf() method
func BenchmarkLogPkg_Printf(b *testing.B) {
	implementation = &logPkg{}
	for i := 0; i < b.N; i++ {
		implementation.Printf("test this method: %s", "TestPrintf")
	}
}

// TestLogPkg_Println test log package LogPkg_Prinln() method
func TestLogPkg_Println(t *testing.T) {
	implementation = &logPkg{}

	captured := captureOutput(func() {
		implementation.Println("test this method: TestPrintln")
	})

	if !strings.Contains(captured, "test this method: TestPrintln") {
		t.Fatalf("expected string: %s got: %s", "test this method: TestPrintln", captured)
	}
}

// BenchmarkLogPkg_Println benchmarks the LogPkg_Println() method
func BenchmarkLogPkg_Println(b *testing.B) {
	implementation = &logPkg{}
	for i := 0; i < b.N; i++ {
		implementation.Println("test this method: TestPrintln")
	}
}

// TestMakeError test making an error struct and MakeError() method
func TestMakeError(t *testing.T) {
	err := MakeError("myKey", "myValue")
	if err.Key() != "myKey" {
		t.Fatalf("expected value: %s, got: %s", "myKey", err.Key())
	}
	if err.Value() != "myValue" {
		t.Fatalf("expected value: %s, got: %s", "myValue", err.Value())
	}
	if err.Error() != `{"key":"myKey","value":"myValue"}` {
		t.Fatalf("expected value: %s, got: %s", `{"key":"myKey","value":"myValue"}`, err.Error())
	}
}

// BenchmarkMakeError benchmarks the MakeError() method
func BenchmarkMakeError(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = MakeError("myKey", "myValue")
	}
}
