package log

import (
	"log"
	"os"
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

var (
	selectedLogLevel = levelInfo
	staticData       = ""

	logStaticDataSet sync.Once
)

func Info(v ...any) {
	if selectedLogLevel <= levelInfo {
		log.New(os.Stdout, "INF: "+time.Now().UTC().String()+" | "+staticData, flags).Println(v...)
	}
}

func Warn(v ...any) {
	if selectedLogLevel <= levelWarn {
		log.New(os.Stdout, "WRN: "+time.Now().UTC().String()+" | "+staticData, flags).Println(v...)
	}
}

func Error(v ...any) {
	if selectedLogLevel <= levelError {
		log.New(os.Stdout, "ERR: "+time.Now().UTC().String()+" | "+staticData, flags).Println(v...)
	}
}

func Debug(v ...any) {
	if selectedLogLevel <= levelDebug {
		log.New(os.Stdout, "DBG: "+time.Now().UTC().String()+" | "+staticData, flags).Println(v...)
	}
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

func SetStaticData(data string) {
	if data == "" {
		return
	}
	logStaticDataSet.Do(func() {
		staticData = data + " | "
	})
}
