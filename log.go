package log

import (
	"fmt"
	"os"
	"sync"
	"time"
)

// Level ...
type Level int

const (
	debug Level = 0
	info  Level = 1
	warn  Level = 2
	err   Level = 3
	fatal Level = 4
)

// Logger ...
type Logger struct {
	filename   string
	mutex      *sync.Mutex
	level      Level
	timeformat string
}

// NewLogger ...
func NewLogger(filename string, timeformat string) *Logger {
	l := &Logger{
		filename:   filename,
		mutex:      &sync.Mutex{},
		level:      debug,
		timeformat: timeformat,
	}
	return l
}

// TimeFormat ...
func (l *Logger) TimeFormat(timeformat string) {
	l.timeformat = timeformat
}

// SetLevel ...
func (l *Logger) SetLevel(level string) {
	switch level {
	case "DEBUG":
		l.level = debug
		break
	case "INFO":
		l.level = info
		break
	case "WARN":
		l.level = warn
		break
	case "ERROR":
		l.level = err
		break
	case "FATAL":
		l.level = fatal
		break
	default:
		l.level = debug
	}
}

// Debug ...
func (l *Logger) Debug(msg string) {
	l.log(debug, msg)
}

// Debugf ...
func (l *Logger) Debugf(msg string, values ...interface{}) {
	message := fmt.Sprintf(msg, values...)
	l.log(debug, message)
}

// Info ...
func (l *Logger) Info(msg string) {
	l.log(info, msg)
}

// Infof ...
func (l *Logger) Infof(msg string, values ...interface{}) {
	message := fmt.Sprintf(msg, values...)
	l.log(info, message)
}

// Warn ...
func (l *Logger) Warn(msg string) {
	l.log(warn, msg)
}

// Warnf ...
func (l *Logger) Warnf(msg string, values ...interface{}) {
	message := fmt.Sprintf(msg, values...)
	l.log(warn, message)
}

// Error ...
func (l *Logger) Error(msg string) {
	l.log(err, msg)
}

// Errorf ...
func (l *Logger) Errorf(msg string, values ...interface{}) {
	message := fmt.Sprintf(msg, values...)
	l.log(err, message)
}

// Fatal ...
func (l *Logger) Fatal(msg string) {
	l.log(fatal, msg)
}

// Fatalf ...
func (l *Logger) Fatalf(msg string, values ...interface{}) {
	message := fmt.Sprintf(msg, values...)
	l.log(fatal, message)
}

func (l *Logger) log(level Level, msg string) {
	if level < l.level {
		return
	}

	l.mutex.Lock()
	defer l.mutex.Unlock()

	fmtMsg := fmt.Sprintf("%s %s: %s%s\n", time.Now().Format(l.timeformat), level.toString(), msg)

	logFile, err := os.OpenFile(l.filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
	if err != nil {
		os.Stderr.WriteString((fmt.Sprintf("go-log: %s%s\n", fmtMsg)))
		return
	}
	defer func() {
		err := logFile.Close()
		if err != nil {
			os.Stderr.WriteString(fmt.Sprintf("Closing failed on log file: %s\n", l.filename))
			return
		}
	}()

	_, err = logFile.WriteString(fmtMsg)
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("Writing failed to log file: %s\n", l.filename))
		return
	}
}

// toString ...
func (level Level) toString() string {
	levels := [...]string{
		"DEBUG",
		"INFO",
		"WARN",
		"ERROR",
		"FATAL",
	}
	if level < debug || level > fatal {
		return levels[debug]
	}
	return levels[level]
}
