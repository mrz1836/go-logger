package logger

import (
	"strings"
	"testing"
)

func FuzzFileTag(f *testing.F) {
	f.Add(0)
	f.Add(1)
	f.Add(2)
	f.Add(5)
	f.Add(-1)
	f.Add(100)

	f.Fuzz(func(t *testing.T, level int) {
		defer func() {
			if r := recover(); r != nil {
				t.Logf("FileTag panicked with level %d: %v", level, r)
			}
		}()

		result := FileTag(level)

		if result != "" && !strings.Contains(result, ":") {
			t.Errorf("FileTag(%d) should contain ':' separator, got: %s", level, result)
		}
	})
}

func FuzzFileTagComponents(f *testing.F) {
	f.Add(0)
	f.Add(1)
	f.Add(2)
	f.Add(5)
	f.Add(-1)
	f.Add(100)

	f.Fuzz(func(t *testing.T, level int) {
		defer func() {
			if r := recover(); r != nil {
				t.Logf("FileTagComponents panicked with level %d: %v", level, r)
			}
		}()

		result := FileTagComponents(level)

		if len(result) != 3 && len(result) != 0 {
			t.Errorf("FileTagComponents(%d) should return 3 components or empty, got %d components: %v", level, len(result), result)
		}
	})
}

func FuzzData(f *testing.F) {
	f.Add(2, uint8(0), "test message")
	f.Add(1, uint8(1), "")
	f.Add(5, uint8(2), "special chars: \n\t\r\"'")
	f.Add(3, uint8(3), "unicode: test rocket cafe")
	f.Add(0, uint8(0), strings.Repeat("a", 10000))

	f.Fuzz(func(t *testing.T, stackLevel int, logLevelRaw uint8, message string) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Data panicked: stackLevel=%d, logLevel=%d, message=%q, panic: %v",
					stackLevel, logLevelRaw, message, r)
			}
		}()

		logLevel := LogLevel(logLevelRaw % 4)

		oldImpl := GetImplementation()
		defer SetImplementation(oldImpl)

		testLogger := &testLoggerImpl{}
		SetImplementation(testLogger)

		Data(stackLevel, logLevel, message)

		if len(testLogger.messages) == 0 {
			t.Error("Data should have produced at least one log message")
		}

		for _, msg := range testLogger.messages {
			if !strings.Contains(msg, `type="`) {
				t.Error("Data output should contain type field")
			}
			if !strings.Contains(msg, `message="`) {
				t.Error("Data output should contain message field")
			}
		}
	})
}

func FuzzNoFileData(f *testing.F) {
	f.Add(uint8(0), "test message")
	f.Add(uint8(1), "")
	f.Add(uint8(2), "special chars: \n\t\r\"'")
	f.Add(uint8(3), "unicode: test rocket cafe")
	f.Add(uint8(0), strings.Repeat("x", 10000))

	f.Fuzz(func(t *testing.T, logLevelRaw uint8, message string) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("NoFileData panicked: logLevel=%d, message=%q, panic: %v",
					logLevelRaw, message, r)
			}
		}()

		logLevel := LogLevel(logLevelRaw % 4)

		oldImpl := GetImplementation()
		defer SetImplementation(oldImpl)

		testLogger := &testLoggerImpl{}
		SetImplementation(testLogger)

		NoFileData(logLevel, message)

		if len(testLogger.messages) == 0 {
			t.Error("NoFileData should have produced at least one log message")
		}

		for _, msg := range testLogger.messages {
			if !strings.Contains(msg, `type="`) {
				t.Error("NoFileData output should contain type field")
			}
			if !strings.Contains(msg, `message="`) {
				t.Error("NoFileData output should contain message field")
			}
			if strings.Contains(msg, `file="`) {
				t.Error("NoFileData output should not contain file field")
			}
		}
	})
}

func FuzzLogLevel_String(f *testing.F) {
	f.Add(uint8(0))
	f.Add(uint8(1))
	f.Add(uint8(2))
	f.Add(uint8(3))
	f.Add(uint8(255))

	f.Fuzz(func(t *testing.T, levelRaw uint8) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("LogLevel.String panicked with level %d: %v", levelRaw, r)
			}
		}()

		level := LogLevel(levelRaw)
		result := level.String()

		validLevels := []string{"debug", "info", "warn", "error", ""}
		found := false
		for _, valid := range validLevels {
			if result == valid {
				found = true
				break
			}
		}

		if !found {
			t.Errorf("LogLevel(%d).String() returned invalid level: %q", levelRaw, result)
		}
	})
}
