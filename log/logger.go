package log

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

type (
	level int

	Level string
)

const (
	flags = 0
)

const (
	LevelInfo  Level = "info"
	LevelWarn  Level = "warn"
	LevelError Level = "error"
	LevelDebug Level = "debug"
)

const (
	levelInfo level = iota
	levelDebug
	levelWarn
	levelError
)

const (
	timeFormat = "2006-01-02 15:04:05"
)

var (
	selectedLogLevel = levelInfo
	staticData       = "| "

	logStaticDataSet sync.Once
)

func newLoggerWithCaller(lvl string, noCallerAttached bool) *log.Logger {
	callStackSkip := 3
	withCaller := ""
	if !noCallerAttached {
		if _, caller, line, ok := runtime.Caller(callStackSkip); ok {
			withCaller = fmt.Sprintf("caller:%s:%d | ", filepath.Base(caller), line)
		}
	}
	return log.New(os.Stdout, "["+time.Now().UTC().Format(timeFormat)+"] "+lvl+" "+staticData+withCaller, flags)
}

type T interface {
	Info(v ...any)
	Warn(v ...any)
	Error(v ...any)
	Debug(v ...any)
}

type internalLog struct {
	noCallerAttached bool
}

func (i internalLog) Info(v ...any) {
	if selectedLogLevel <= levelInfo {
		newLoggerWithCaller("INF", i.noCallerAttached).Println(v...)
	}
}

func (i internalLog) Warn(v ...any) {
	if selectedLogLevel <= levelWarn {
		newLoggerWithCaller("WRN", i.noCallerAttached).Println(v...)
	}
}

func (i internalLog) Error(v ...any) {
	if selectedLogLevel <= levelError {
		newLoggerWithCaller("ERR", i.noCallerAttached).Println(v...)
	}
}

func (i internalLog) Debug(v ...any) {
	if selectedLogLevel <= levelDebug {
		newLoggerWithCaller("DBG", i.noCallerAttached).Println(v...)
	}
}

func WithNoCaller() T {
	return &internalLog{
		noCallerAttached: true,
	}
}

func Info(v ...any) {
	internalLog{}.Info(v...)
}

func Warn(v ...any) {
	internalLog{}.Warn(v...)
}

func Error(v ...any) {
	internalLog{}.Error(v...)
}

func Debug(v ...any) {
	internalLog{}.Debug(v...)
}

func SetLogLevel(lvl string) {
	switch Level(lvl) {
	case LevelInfo:
		selectedLogLevel = levelInfo
	case LevelWarn:
		selectedLogLevel = levelWarn
	case LevelError:
		selectedLogLevel = levelError
	case LevelDebug:
		selectedLogLevel = levelDebug
	}
}

type Setter interface {
	Add(key, value string) Setter
	Set()
}

type set struct {
	setBuffer *strings.Builder
}

func GetStaticDataSetter() Setter {
	return &set{
		setBuffer: &strings.Builder{},
	}
}

func (s *set) Add(key, value string) Setter {
	s.setBuffer.WriteString(fmt.Sprintf("%s:%s | ", key, value))
	return s
}

func (s *set) Set() {
	logStaticDataSet.Do(func() {
		data := s.setBuffer.String()
		if data == "" {
			return
		}
		staticData += data
	})
}
