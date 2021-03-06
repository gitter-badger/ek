// Package log provides improved logger
package log

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                     Copyright (c) 2009-2016 Essential Kaos                         //
//      Essential Kaos Open Source License <http://essentialkaos.com/ekol?en>         //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// DEBUG debug messages
// INFO info messages
// WARN warning messages
// ERROR error messages
// CRIT critical error messages
// AUX unskipable messages (separators, headers, etc...)
const (
	DEBUG = 0
	INFO  = 1
	WARN  = 2
	ERROR = 3
	CRIT  = 4
	AUX   = 99
)

// ////////////////////////////////////////////////////////////////////////////////// //

// Logger is a basic logger struct
type Logger struct {
	PrefixDebug bool // Prefix for debug messages
	PrefixInfo  bool // Prefix for info messages
	PrefixWarn  bool // Prefix for warning messages
	PrefixError bool // Prefix for error messages
	PrefixCrit  bool // Prefix for critical error messages

	file     string
	fd       *os.File
	w        *bufio.Writer
	level    int
	perms    os.FileMode
	useBufIO bool
}

// ////////////////////////////////////////////////////////////////////////////////// //

// Global is global logger struct
var Global = &Logger{
	PrefixWarn:  true,
	PrefixError: true,
	PrefixCrit:  true,

	level: INFO,
}

// ////////////////////////////////////////////////////////////////////////////////// //

// PrefixMap is map with messages prefixes
var PrefixMap = map[int]string{
	DEBUG: "[DEBUG]",
	INFO:  "[INFO]",
	WARN:  "[WARNING]",
	ERROR: "[ERROR]",
	CRIT:  "[CRITICAL]",
}

// TimeFormat contains format string for time in logs
var TimeFormat = "2006/01/02 15:04:05.000"

// ////////////////////////////////////////////////////////////////////////////////// //

var logLevelsNames = map[string]int{
	"debug":    0,
	"info":     1,
	"warn":     2,
	"warning":  2,
	"error":    3,
	"crit":     4,
	"critical": 4,
}

// ////////////////////////////////////////////////////////////////////////////////// //

// New creates new logger struct
func New(file string, perms os.FileMode) (*Logger, error) {
	logger := &Logger{
		PrefixWarn:  true,
		PrefixCrit:  true,
		PrefixError: true,

		level: INFO,
	}

	err := logger.Set(file, perms)

	if err != nil {
		return nil, err
	}

	return logger, nil
}

// Reopen close file descriptor for global logger and open it again
// Useful for log rotation
func Reopen() error {
	return Global.Reopen()
}

// MinLevel defines minimal logging level
func MinLevel(level interface{}) error {
	return Global.MinLevel(level)
}

// Set change global logger output target
func Set(file string, perms os.FileMode) error {
	return Global.Set(file, perms)
}

// EnableBufIO enable buffered I/O
func EnableBufIO(interval time.Duration) {
	Global.EnableBufIO(interval)
}

// Flush write buffered data to file
func Flush() error {
	return Global.Flush()
}

// Print write message to global logger output
func Print(level int, f string, a ...interface{}) (int, error) {
	return Global.Print(level, f, a...)
}

// Debug write debug message to global logger output
func Debug(f string, a ...interface{}) (int, error) {
	return Global.Debug(f, a...)
}

// Info write info message to global logger output
func Info(f string, a ...interface{}) (int, error) {
	return Global.Info(f, a...)
}

// Warn write warning message to global logger output
func Warn(f string, a ...interface{}) (int, error) {
	return Global.Warn(f, a...)
}

// Error write error message to global logger output
func Error(f string, a ...interface{}) (int, error) {
	return Global.Error(f, a...)
}

// Crit write critical message to global logger output
func Crit(f string, a ...interface{}) (int, error) {
	return Global.Crit(f, a...)
}

// Aux write unskipable message (for separators/headers)
func Aux(f string, a ...interface{}) (int, error) {
	return Global.Aux(f, a...)
}

// ////////////////////////////////////////////////////////////////////////////////// //

// Reopen close file descriptor and open again
// Useful for log rotation
func (l *Logger) Reopen() error {
	if l == nil {
		return errors.New("Logger is nil")
	}

	if l.fd == nil {
		return errors.New("Output file is not set")
	}

	if l.w != nil {
		l.w.Flush()
	}

	l.fd.Close()

	return l.Set(l.file, l.perms)
}

// MinLevel defines minimal logging level
func (l *Logger) MinLevel(level interface{}) error {
	if l == nil {
		return errors.New("Logger is nil")
	}

	levelCode := INFO

	switch level.(type) {

	case int:
		levelCode = level.(int)

	case int8:
		levelCode = int(level.(int8))

	case int16:
		levelCode = int(level.(int16))

	case int32:
		levelCode = int(level.(int32))

	case int64:
		levelCode = int(level.(int64))

	case uint:
		levelCode = int(level.(uint))

	case uint8:
		levelCode = int(level.(uint8))

	case uint16:
		levelCode = int(level.(uint16))

	case uint32:
		levelCode = int(level.(uint32))

	case uint64:
		levelCode = int(level.(uint64))

	case float32:
		levelCode = int(level.(float32))

	case float64:
		levelCode = int(level.(float64))

	case string:
		code, ok := logLevelsNames[strings.ToLower(level.(string))]

		if !ok {
			return errors.New("Unknown level " + level.(string))
		}

		levelCode = code

	default:
		return errors.New("Unexpected level type")
	}

	switch {
	case levelCode < DEBUG:
		levelCode = DEBUG
	case levelCode > CRIT:
		levelCode = CRIT
	}

	l.level = levelCode

	return nil
}

// EnableBufIO enable buffered I/O support
func (l *Logger) EnableBufIO(interval time.Duration) {
	l.useBufIO = true

	if l.fd != nil {
		l.w = bufio.NewWriter(l.fd)
	}

	go l.flushDaemon(interval)
}

// Set change logger output target
func (l *Logger) Set(file string, perms os.FileMode) error {
	fd, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_APPEND, perms)

	if err != nil {
		return err
	}

	// Flush data if writter exist
	if l.w != nil {
		l.w.Flush()
		l.w = nil
	}

	if l.fd != nil {
		l.fd.Close()
		l.fd = nil
	}

	l.fd, l.file, l.perms = fd, file, perms

	if l.useBufIO {
		l.w = bufio.NewWriter(l.fd)
	}

	return nil
}

// Print write message to logger output
func (l *Logger) Print(level int, f string, a ...interface{}) (int, error) {
	if l == nil {
		return -1, errors.New("Logger is nil")
	}

	if l.level > level {
		return 0, nil
	}

	var w io.Writer

	if l.fd == nil {
		switch level {
		case ERROR, CRIT:
			w = os.Stderr
		default:
			w = os.Stdout
		}
	} else {
		if l.w != nil {
			w = l.w
		} else {
			w = l.fd
		}
	}

	var showPrefixes bool

	switch {
	case level == DEBUG && l.PrefixDebug,
		level == INFO && l.PrefixInfo,
		level == WARN && l.PrefixWarn,
		level == ERROR && l.PrefixError,
		level == CRIT && l.PrefixCrit:
		showPrefixes = true
	}

	if f[len(f)-1:] != "\n" {
		f += "\n"
	}

	if showPrefixes {
		return fmt.Fprintf(w, "%s %s %s", getTime(), PrefixMap[level], fmt.Sprintf(f, a...))
	}

	return fmt.Fprintf(w, "%s %s", getTime(), fmt.Sprintf(f, a...))
}

// Flush write buffered data to file
func (l *Logger) Flush() error {
	if l == nil {
		return errors.New("Logger is nil")
	}

	if l.w == nil {
		return nil
	}

	return l.w.Flush()
}

// Debug write debug message to logger output
func (l *Logger) Debug(f string, a ...interface{}) (int, error) {
	if l == nil {
		return -1, errors.New("Logger is nil")
	}

	return l.Print(DEBUG, f, a...)
}

// Info write info message to logger output
func (l *Logger) Info(f string, a ...interface{}) (int, error) {
	if l == nil {
		return -1, errors.New("Logger is nil")
	}

	return l.Print(INFO, f, a...)
}

// Warn write warning message to logger output
func (l *Logger) Warn(f string, a ...interface{}) (int, error) {
	if l == nil {
		return -1, errors.New("Logger is nil")
	}

	return l.Print(WARN, f, a...)
}

// Error write error message to logger output
func (l *Logger) Error(f string, a ...interface{}) (int, error) {
	if l == nil {
		return -1, errors.New("Logger is nil")
	}

	return l.Print(ERROR, f, a...)
}

// Crit write critical message to logger output
func (l *Logger) Crit(f string, a ...interface{}) (int, error) {
	if l == nil {
		return -1, errors.New("Logger is nil")
	}

	return l.Print(CRIT, f, a...)
}

// Aux write unskipable message (for separators/headers)
func (l *Logger) Aux(f string, a ...interface{}) (int, error) {
	if l == nil {
		return -1, errors.New("Logger is nil")
	}

	return l.Print(AUX, f, a...)
}

// ////////////////////////////////////////////////////////////////////////////////// //

func (l *Logger) flushDaemon(interval time.Duration) {
	for _ = range time.NewTicker(interval).C {
		l.Flush()
	}
}

// ////////////////////////////////////////////////////////////////////////////////// //

func getTime() string {
	return "[ " + time.Now().Format(TimeFormat) + " ]"
}
