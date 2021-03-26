package logger

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
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
	assert.Equal(t, "debug", level.String())

	// Test for info
	level = 1
	assert.Equal(t, "info", level.String())

	// Test for warn
	level = 2
	assert.Equal(t, "warn", level.String())

	// Test for error
	level = 3
	assert.Equal(t, "error", level.String())

	// Test for empty
	level = 4
	assert.Equal(t, "", level.String())
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
	assert.Contains(t, fileTag, "go-logger/logger_test.go:go-logger.TestFileTag:")
}

// ExampleFileTag example using FileTag()
func ExampleFileTag() {
	// fileTag := FileTag(1)
	fileTag := "go-logger/logger_test.go:go-logger.ExampleFileTag:102"
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
	assert.NotEqual(t, 0, len(fileTagComps))
	assert.Equal(t, 3, len(fileTagComps))

	// Test the part
	assert.Equal(t, "testing/testing.go", fileTagComps[0])

	// Test the part
	assert.Equal(t, "testing.tRunner", fileTagComps[1])

	// Test the part // todo: this number changes frequently, maybe this is not the best test?
	assert.NotEqual(t, 0, fileTagComps[2])

	// Test the level: 1
	fileTagComps = FileTagComponents(1)
	assert.NotEqual(t, 0, len(fileTagComps))
	assert.Equal(t, 3, len(fileTagComps))

	// Test the part
	assert.Equal(t, "go-logger/logger_test.go", fileTagComps[0])

	// Test the part
	assert.Equal(t, "go-logger.TestFileTagComponents", fileTagComps[1])
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

	assert.Contains(t, captured, "go-logger/logger_test.go:go-logger.TestPrintln")
	assert.Contains(t, captured, "test this method")
}

// TestPrint test the print method
func TestPrint(t *testing.T) {
	captured := captureOutput(func() {
		Print("test this method")
	})

	assert.Contains(t, captured, "go-logger/logger_test.go:go-logger.TestPrint")
	assert.Contains(t, captured, "test this method")
}

// TestNoFilePrintln test the print line method
func TestNoFilePrintln(t *testing.T) {
	captured := captureOutput(func() {
		NoFilePrintln("test this method")
	})
	assert.Contains(t, captured, "test this method")
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

	assert.Contains(t, captured, "go-logger/logger_test.go:go-logger.TestPrintf")
	assert.Contains(t, captured, "test this method: TestPrintf")
}

// TestNoFilePrintf test the print fmt method
func TestNoFilePrintf(t *testing.T) {
	captured := captureOutput(func() {
		NoFilePrintf("test this method: %s", "TestPrintf")
	})

	assert.Contains(t, captured, "test this method: TestPrintf")
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

	assert.Contains(t, captured, "go-logger/logger_test.go:go-logger.TestErrorln")
	assert.Contains(t, captured, "test this method")
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

	assert.Contains(t, captured, "go-logger/logger_test.go:go-logger.TestErrorfmt")
	assert.Contains(t, captured, "test this method: Errorfmt")
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
		Data(2, WARN, "test this method", MakeParameter("another", "value"))
	})

	// 2019/06/17 12:59:32 type="warn" file="go-logger/logger_test.go" method="go-logger.TestData.func1" line="188" message="test this method" another="value"

	// Check for warn
	assert.Contains(t, captured, `type="warn"`)

	// Check for file
	assert.Contains(t, captured, `file="go-logger/logger_test.go"`)

	// Check for method
	assert.Contains(t, captured, `method="go-logger.TestData.func1"`)

	// Check for message
	assert.Contains(t, captured, `message="test this method"`)

	// Check for additional values
	assert.Contains(t, captured, `another="value"`)
}

// TestNoFileData test the NoFileData() method
func TestNoFileData(t *testing.T) {
	captured := captureOutput(func() {
		NoFileData(WARN, "test this method", MakeParameter("another", "value"))
	})

	// Check for warn
	assert.Contains(t, captured, `type="warn"`)

	// Check for message
	assert.Contains(t, captured, `message="test this method"`)

	// Check for additional values
	assert.Contains(t, captured, `another="value"`)
}

// BenchmarkData benchmarks the Data() method
func BenchmarkData(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Data(2, WARN, "test this method", MakeParameter("another", "value"))
	}
}

// TestLogPkg_Printf test log package Printf() method
func TestLogPkg_Printf(t *testing.T) {
	implementation = &logPkg{}

	captured := captureOutput(func() {
		implementation.Printf("test this method: %s", "TestPrintf")
	})

	assert.Contains(t, captured, "test this method: TestPrintf")
}

// BenchmarkLogPkg_Printf benchmarks the LogPkg_Printf() method
func BenchmarkLogPkg_Printf(b *testing.B) {
	implementation = &logPkg{}
	for i := 0; i < b.N; i++ {
		implementation.Printf("test this method: %s", "TestPrintf")
	}
}

// TestLogPkg_Println test log package LogPkg_Println() method
func TestLogPkg_Println(t *testing.T) {
	implementation = &logPkg{}

	captured := captureOutput(func() {
		implementation.Println("test this method: TestPrintln")
	})

	assert.Contains(t, captured, "test this method: TestPrintln")
}

// BenchmarkLogPkg_Println benchmarks the LogPkg_Println() method
func BenchmarkLogPkg_Println(b *testing.B) {
	implementation = &logPkg{}
	for i := 0; i < b.N; i++ {
		implementation.Println("test this method: TestPrintln")
	}
}

// TestFatalf will test the Fatalf() method
func TestFatalf(t *testing.T) {

	client, err := NewLogEntriesClient(testToken, LogEntriesTestEndpoint, LogEntriesPort)
	assert.NoError(t, err)
	assert.NotNil(t, client)

	SetImplementation(client)

	theImplementation := GetImplementation()
	assert.NotNil(t, theImplementation)

	if os.Getenv("EXIT_FUNCTION") == "1" {
		Fatalf("test %d", 1)
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestFatalf")
	cmd.Env = append(os.Environ(), "EXIT_FUNCTION=1")
	err = cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	t.Fatalf("process ran with err %v, want exit status 1", err)
}

// TestFatal will test the Fatal() method
func TestFatal(t *testing.T) {

	client, err := NewLogEntriesClient(testToken, LogEntriesTestEndpoint, LogEntriesPort)
	assert.NoError(t, err)
	assert.NotNil(t, client)

	SetImplementation(client)

	if os.Getenv("EXIT_FUNCTION") == "1" {
		Fatal("test exit")
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestFatal")
	cmd.Env = append(os.Environ(), "EXIT_FUNCTION=1")
	err = cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	t.Fatalf("process ran with err %v, want exit status 1", err)
}

// TestFatalln will test the Fatalln() method
func TestFatalln(t *testing.T) {

	client, err := NewLogEntriesClient(testToken, LogEntriesTestEndpoint, LogEntriesPort)
	assert.NoError(t, err)
	assert.NotNil(t, client)

	SetImplementation(client)

	if os.Getenv("EXIT_FUNCTION") == "1" {
		Fatalln("test exit")
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestFatalln")
	cmd.Env = append(os.Environ(), "EXIT_FUNCTION=1")
	err = cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	t.Fatalf("process ran with err %v, want exit status 1", err)
}

// TestPanic will test the Panic() method
func TestPanic(t *testing.T) {

	client, err := NewLogEntriesClient(testToken, LogEntriesTestEndpoint, LogEntriesPort)
	assert.NoError(t, err)
	assert.NotNil(t, client)

	SetImplementation(client)

	if os.Getenv("EXIT_FUNCTION") == "1" {
		Panic("test exit")
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestPanic")
	cmd.Env = append(os.Environ(), "EXIT_FUNCTION=1")
	err = cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	t.Fatalf("process ran with err %v, want exit status 1", err)
}

// TestPanicln will test the Panicln() method
func TestPanicln(t *testing.T) {

	client, err := NewLogEntriesClient(testToken, LogEntriesTestEndpoint, LogEntriesPort)
	assert.NoError(t, err)
	assert.NotNil(t, client)

	SetImplementation(client)

	if os.Getenv("EXIT_FUNCTION") == "1" {
		Panicln("test exit")
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestPanicln")
	cmd.Env = append(os.Environ(), "EXIT_FUNCTION=1")
	err = cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	t.Fatalf("process ran with err %v, want exit status 1", err)
}

// TestPanicf will test the Panicf() method
func TestPanicf(t *testing.T) {

	client, err := NewLogEntriesClient(testToken, LogEntriesTestEndpoint, LogEntriesPort)
	assert.NoError(t, err)
	assert.NotNil(t, client)

	SetImplementation(client)

	theImplementation := GetImplementation()
	assert.NotNil(t, theImplementation)

	if os.Getenv("EXIT_FUNCTION") == "1" {
		Panicf("test %d", 1)
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestPanicf")
	cmd.Env = append(os.Environ(), "EXIT_FUNCTION=1")
	err = cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	t.Fatalf("process ran with err %v, want exit status 1", err)
}
