package logger

import (
	"context"
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// GormLoggerInterface is a logger interface to help work with GORM
type GormLoggerInterface interface {
	Error(context.Context, string, ...interface{})
	GetMode() GormLogLevel
	GetStackLevel() int
	Info(context.Context, string, ...interface{})
	SetMode(GormLogLevel) GormLoggerInterface
	SetStackLevel(level int)
	Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error)
	Warn(context.Context, string, ...interface{})
}

// GormLogLevel is the GORM log level
type GormLogLevel int

const (
	// Silent silent log level
	Silent GormLogLevel = iota + 1

	// Error error log level
	Error

	// Warn warn log level
	Warn

	// Info info log level
	Info
)

// This is the time cut-off for considering a query as "slow"
const slowQueryThreshold = 5 * time.Second

/*// Sourced file directory path
var gormSourceDir string

// On init, get the source file and store as a variable
func init() {
	_, file, _, _ := runtime.Caller(0) // nolint: dogsled // other variables not needed
	// compatible solution to get gorm source directory with various operating systems
	gormSourceDir = regexp.MustCompile(`gorm\.go`).ReplaceAllString(file, "")
}*/

// NewGormLogger will return a basic logger interface
func NewGormLogger(debugging bool, stackLevel int) GormLoggerInterface {
	logLevel := Warn
	if debugging {
		logLevel = Info
	}
	return &basicGormLogger{
		logLevel:   logLevel,
		stackLevel: stackLevel,
	}
}

// basicGormLogger is a basic implementation of the logger interface if no custom logger is provided
type basicGormLogger struct {
	logLevel   GormLogLevel // Log level (info, error, etc)
	stackLevel int          // How many files/functions to traverse upwards to record the file/line
}

// SetMode will set the log mode
func (l *basicGormLogger) SetMode(level GormLogLevel) GormLoggerInterface {
	newLogger := *l
	newLogger.logLevel = level
	return &newLogger
}

// SetStackLevel will set the stack level
func (l *basicGormLogger) SetStackLevel(level int) {
	l.stackLevel = level
}

// GetStackLevel will get the current stack level
func (l *basicGormLogger) GetStackLevel() int {
	return l.stackLevel
}

// GetMode will get the log mode
func (l *basicGormLogger) GetMode() GormLogLevel {
	return l.logLevel
}

// Info print information
func (l *basicGormLogger) Info(_ context.Context, message string, params ...interface{}) {
	if l.logLevel >= Info {
		displayLog(INFO, l.stackLevel, message, params...)
	}
}

// Warn print warn messages
func (l *basicGormLogger) Warn(_ context.Context, message string, params ...interface{}) {
	if l.logLevel >= Warn {
		displayLog(WARN, l.stackLevel, message, params...)
	}
}

// Error print error messages
func (l *basicGormLogger) Error(_ context.Context, message string, params ...interface{}) {
	if l.logLevel >= Error {
		displayLog(ERROR, l.stackLevel, message, params...)
	}
}

// Trace is for GORM/SQL tracing from datastore
func (l *basicGormLogger) Trace(_ context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if l.logLevel <= Silent {
		return
	}
	elapsed := time.Since(begin)
	switch {
	case err != nil && l.logLevel >= Error && (!strings.Contains(err.Error(), "record not found")):
		sql, rows := fc()
		Data(l.stackLevel, ERROR,
			"error executing query",
			MakeParameter("file", fileWithLineNum()),
			MakeParameter("error", err.Error()),
			MakeParameter("duration", fmt.Sprintf("%.3fms", float64(elapsed.Nanoseconds())/1e6)),
			MakeParameter("rows", rows),
			MakeParameter("sql", sql),
		)
	case elapsed > slowQueryThreshold && l.logLevel >= Warn:
		sql, rows := fc()
		Data(l.stackLevel, WARN,
			"warning executing query",
			MakeParameter("file", fileWithLineNum()),
			MakeParameter("slow_log", fmt.Sprintf("SLOW SQL >= %v", slowQueryThreshold)),
			MakeParameter("duration", fmt.Sprintf("%.3fms", float64(elapsed.Nanoseconds())/1e6)),
			MakeParameter("rows", rows),
			MakeParameter("sql", sql),
		)
	case l.logLevel == Info:
		sql, rows := fc()
		Data(l.stackLevel, INFO,
			"executing sql query",
			MakeParameter("file", fileWithLineNum()),
			MakeParameter("duration", fmt.Sprintf("%.3fms", float64(elapsed.Nanoseconds())/1e6)),
			MakeParameter("rows", rows),
			MakeParameter("sql", sql),
		)
	}
}

// displayLog will display a log using logger
func displayLog(level LogLevel, stackLevel int, message string, params ...interface{}) {
	var keyValues []KeyValue
	if len(params) > 0 {
		for index, val := range params {
			keyValues = append(keyValues, MakeParameter(fmt.Sprintf("param_%d", index), val))
		}
	}
	Data(stackLevel, level, message, keyValues...)
}

// fileWithLineNum return the file name and line number of the current file
// This is originally from GORM: https://github.com/go-gorm/gorm/blob/7837fb6fa001ef78bc76e66b48445dee7b2db37b/utils/utils.go#L23
// Copied method in order to not make GORM a dependency of the project for this tiny utility method
func fileWithLineNum() string {
	// the first & second caller usually from gorm internal, so set index start from 3
	for i := 2; i < 15; i++ {
		_, file, line, ok := runtime.Caller(i)
		if ok && (!strings.HasSuffix(file, "_test.go") &&
			!strings.Contains(file, "gorm.go") &&
			!strings.Contains(file, "callbacks.go") &&
			!strings.Contains(file, "finisher_api.go")) {
			return file + ":" + strconv.FormatInt(int64(line), 10)
		}
	}
	return ""
}
