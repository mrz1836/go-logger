package logger

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"
)

type testError struct {
	message string
}

func (e *testError) Error() string {
	return e.message
}

type testLoggerImpl struct {
	messages []string
}

func (t *testLoggerImpl) Fatal(v ...interface{}) { t.messages = append(t.messages, fmt.Sprint(v...)) }
func (t *testLoggerImpl) Fatalf(format string, v ...interface{}) {
	t.messages = append(t.messages, fmt.Sprintf(format, v...))
}
func (t *testLoggerImpl) Fatalln(v ...interface{}) { t.messages = append(t.messages, fmt.Sprint(v...)) }
func (t *testLoggerImpl) Panic(v ...interface{})   { t.messages = append(t.messages, fmt.Sprint(v...)) }
func (t *testLoggerImpl) Panicf(s string, v ...interface{}) {
	t.messages = append(t.messages, fmt.Sprintf(s, v...))
}
func (t *testLoggerImpl) Panicln(v ...interface{}) { t.messages = append(t.messages, fmt.Sprint(v...)) }
func (t *testLoggerImpl) Print(v ...interface{})   { t.messages = append(t.messages, fmt.Sprint(v...)) }
func (t *testLoggerImpl) Printf(format string, v ...interface{}) {
	t.messages = append(t.messages, fmt.Sprintf(format, v...))
}
func (t *testLoggerImpl) Println(v ...interface{}) { t.messages = append(t.messages, fmt.Sprint(v...)) }

func FuzzNewGormLogger(f *testing.F) {
	f.Add(true, 1)
	f.Add(false, 2)
	f.Add(true, 0)
	f.Add(false, -1)
	f.Add(true, 100)

	f.Fuzz(func(t *testing.T, debugging bool, stackLevel int) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("NewGormLogger panicked: debugging=%v, stackLevel=%d, panic: %v",
					debugging, stackLevel, r)
			}
		}()

		logger := NewGormLogger(debugging, stackLevel)

		if logger == nil {
			t.Error("NewGormLogger should return non-nil logger")
			return
		}

		basicLogger, ok := logger.(*basicGormLogger)
		if !ok {
			t.Error("NewGormLogger should return *basicGormLogger")
			return
		}

		if basicLogger.stackLevel != stackLevel {
			t.Errorf("Stack level mismatch: expected %d, got %d", stackLevel, basicLogger.stackLevel)
		}

		expectedLogLevel := Warn
		if debugging {
			expectedLogLevel = Info
		}

		if basicLogger.logLevel != expectedLogLevel {
			t.Errorf("Log level mismatch: expected %v, got %v", expectedLogLevel, basicLogger.logLevel)
		}

		if logger.GetStackLevel() != stackLevel {
			t.Errorf("GetStackLevel mismatch: expected %d, got %d", stackLevel, logger.GetStackLevel())
		}

		if logger.GetMode() != expectedLogLevel {
			t.Errorf("GetMode mismatch: expected %v, got %v", expectedLogLevel, logger.GetMode())
		}
	})
}

func FuzzBasicGormLogger_SetMode(f *testing.F) {
	f.Add(uint8(1))
	f.Add(uint8(2))
	f.Add(uint8(3))
	f.Add(uint8(4))
	f.Add(uint8(0))
	f.Add(uint8(255))

	f.Fuzz(func(t *testing.T, levelRaw uint8) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("basicGormLogger.SetMode panicked: level=%d, panic: %v", levelRaw, r)
			}
		}()

		originalLogger := NewGormLogger(true, 2)
		level := GormLogLevel(levelRaw % 5)

		newLogger := originalLogger.SetMode(level)

		if newLogger == nil {
			t.Error("SetMode should return non-nil logger")
			return
		}

		if newLogger == originalLogger {
			t.Error("SetMode should return a new logger instance, not modify the original")
		}

		if newLogger.GetMode() != level {
			t.Errorf("New logger mode mismatch: expected %v, got %v", level, newLogger.GetMode())
		}

		if originalLogger.GetMode() != Info {
			t.Errorf("Original logger mode should still be Info, got %v", originalLogger.GetMode())
		}
	})
}

func FuzzBasicGormLogger_SetStackLevel(f *testing.F) {
	f.Add(0)
	f.Add(1)
	f.Add(5)
	f.Add(-1)
	f.Add(100)
	f.Add(-100)

	f.Fuzz(func(t *testing.T, stackLevel int) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("basicGormLogger.SetStackLevel panicked: stackLevel=%d, panic: %v",
					stackLevel, r)
			}
		}()

		logger := NewGormLogger(true, 2)
		originalLevel := logger.GetStackLevel()

		logger.SetStackLevel(stackLevel)

		if logger.GetStackLevel() != stackLevel {
			t.Errorf("Stack level mismatch: expected %d, got %d", stackLevel, logger.GetStackLevel())
		}

		if logger.GetStackLevel() == originalLevel && stackLevel != originalLevel {
			t.Error("SetStackLevel should have changed the stack level")
		}
	})
}

func FuzzBasicGormLogger_Info(f *testing.F) {
	f.Add("test message")
	f.Add("")
	f.Add("unicode: test rocket")
	f.Add("special\nchars\ttabs")
	f.Add(strings.Repeat("long message ", 1000))
	f.Add("null byte: \x00 embedded")

	f.Fuzz(func(t *testing.T, message string) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("basicGormLogger.Info panicked: message=%q, panic: %v", message, r)
			}
		}()

		oldImpl := GetImplementation()
		defer SetImplementation(oldImpl)

		testLogger := &testLoggerImpl{}
		SetImplementation(testLogger)

		logger := NewGormLogger(true, 3)
		ctx := context.Background()

		logger.Info(ctx, message)

		if len(testLogger.messages) == 0 {
			t.Error("Info should have produced at least one log message")
		}

		silentLogger := logger.SetMode(Silent)
		testLogger.messages = nil
		silentLogger.Info(ctx, message)

		if len(testLogger.messages) != 0 {
			t.Error("Info should not log when mode is Silent")
		}
	})
}

func FuzzBasicGormLogger_Warn(f *testing.F) {
	f.Add("warning message")
	f.Add("")
	f.Add("unicode warning: test")
	f.Add("multiline\nwarning\nmessage")

	f.Fuzz(func(t *testing.T, message string) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("basicGormLogger.Warn panicked: message=%q, panic: %v", message, r)
			}
		}()

		oldImpl := GetImplementation()
		defer SetImplementation(oldImpl)

		testLogger := &testLoggerImpl{}
		SetImplementation(testLogger)

		logger := NewGormLogger(true, 3)
		ctx := context.Background()

		logger.Warn(ctx, message)

		if len(testLogger.messages) == 0 {
			t.Error("Warn should have produced at least one log message")
		}

		silentLogger := logger.SetMode(Silent)
		testLogger.messages = nil
		silentLogger.Warn(ctx, message)

		if len(testLogger.messages) != 0 {
			t.Error("Warn should not log when mode is Silent")
		}
	})
}

func FuzzBasicGormLogger_Error(f *testing.F) {
	f.Add("error message")
	f.Add("")
	f.Add("unicode error: test")
	f.Add("error\nwith\nnewlines")

	f.Fuzz(func(t *testing.T, message string) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("basicGormLogger.Error panicked: message=%q, panic: %v", message, r)
			}
		}()

		oldImpl := GetImplementation()
		defer SetImplementation(oldImpl)

		testLogger := &testLoggerImpl{}
		SetImplementation(testLogger)

		logger := NewGormLogger(true, 3)
		ctx := context.Background()

		logger.Error(ctx, message)

		if len(testLogger.messages) == 0 {
			t.Error("Error should have produced at least one log message")
		}

		silentLogger := logger.SetMode(Silent)
		testLogger.messages = nil
		silentLogger.Error(ctx, message)

		if len(testLogger.messages) != 0 {
			t.Error("Error should not log when mode is Silent")
		}
	})
}

func FuzzBasicGormLogger_Trace(f *testing.F) {
	f.Add("SELECT * FROM users", int64(5), "test error")
	f.Add("", int64(0), "")
	f.Add("INSERT INTO table VALUES (?)", int64(1), "duplicate key")
	f.Add("UPDATE users SET name=?", int64(10), "record not found")
	f.Add(strings.Repeat("SELECT ", 1000), int64(999999), "connection failed")

	f.Fuzz(func(t *testing.T, sql string, rows int64, errMsg string) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("basicGormLogger.Trace panicked: sql=%q, rows=%d, errMsg=%q, panic: %v",
					sql, rows, errMsg, r)
			}
		}()

		oldImpl := GetImplementation()
		defer SetImplementation(oldImpl)

		testLogger := &testLoggerImpl{}
		SetImplementation(testLogger)

		logger := NewGormLogger(true, 3)
		ctx := context.Background()

		begin := time.Now().Add(-time.Millisecond * 100)

		var err error
		if errMsg != "" {
			err = &testError{message: errMsg}
		}

		fc := func() (string, int64) {
			return sql, rows
		}

		logger.Trace(ctx, begin, fc, err)

		if logger.GetMode() == Silent && len(testLogger.messages) != 0 {
			t.Error("Trace should not log when mode is Silent")
		}

		slowBegin := time.Now().Add(-SlowQueryThreshold - time.Second)
		testLogger.messages = nil
		logger.Trace(ctx, slowBegin, fc, nil)

		if logger.GetMode() >= Warn && len(testLogger.messages) == 0 {
			t.Error("Trace should log slow queries when mode >= Warn")
		}

		testLogger.messages = nil
		infoLogger := logger.SetMode(Info)
		infoLogger.Trace(ctx, time.Now(), fc, nil)

		if len(testLogger.messages) == 0 {
			t.Error("Trace should log all queries when mode is Info")
		}
	})
}

func FuzzFileWithLineNum(f *testing.F) {
	f.Add(1)

	f.Fuzz(func(t *testing.T, _ int) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("fileWithLineNum panicked: %v", r)
			}
		}()

		result := fileWithLineNum()

		if result == "" {
			return
		}

		parts := strings.Split(result, ":")
		if len(parts) != 2 {
			t.Errorf("fileWithLineNum should return format 'file:line', got: %q", result)
			return
		}

		fileName := parts[0]
		if strings.HasSuffix(fileName, "_test.go") {
			t.Error("fileWithLineNum should skip test files")
		}

		if strings.Contains(fileName, "gorm.go") {
			t.Error("fileWithLineNum should skip gorm.go file")
		}

		if strings.Contains(fileName, "callbacks.go") {
			t.Error("fileWithLineNum should skip callbacks.go file")
		}

		if strings.Contains(fileName, "finisher_api.go") {
			t.Error("fileWithLineNum should skip finisher_api.go file")
		}
	})
}

func FuzzDisplayLog(f *testing.F) {
	f.Add(uint8(0), 2, "test message")
	f.Add(uint8(1), 3, "")
	f.Add(uint8(2), 1, "unicode: test")
	f.Add(uint8(3), 5, "message\nwith\nnewlines")
	f.Add(uint8(0), -1, "negative stack level")

	f.Fuzz(func(t *testing.T, levelRaw uint8, stackLevel int, message string) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("displayLog panicked: level=%d, stackLevel=%d, message=%q, panic: %v",
					levelRaw, stackLevel, message, r)
			}
		}()

		oldImpl := GetImplementation()
		defer SetImplementation(oldImpl)

		testLogger := &testLoggerImpl{}
		SetImplementation(testLogger)

		level := LogLevel(levelRaw % 4)

		displayLog(level, stackLevel, message)

		if len(testLogger.messages) == 0 {
			t.Error("displayLog should have produced at least one log message")
		}

		for _, msg := range testLogger.messages {
			if !strings.Contains(msg, `type="`) {
				t.Error("displayLog output should contain type field")
			}
			if !strings.Contains(msg, `message="`) {
				t.Error("displayLog output should contain message field")
			}
		}
	})
}
