package log

import (
	"log"
	"os"
	"time"
)

const (
	flags = 0
)

func Info(v ...any) {
	log.New(os.Stdout, "INF: "+time.Now().UTC().String()+" | ", flags).Println(v...)
}

func Warn(v ...any) {
	log.New(os.Stdout, "WRN: "+time.Now().UTC().String()+" | ", flags).Println(v...)
}

func Error(v ...any) {
	log.New(os.Stdout, "ERR: "+time.Now().UTC().String()+" | ", flags).Println(v...)
}
