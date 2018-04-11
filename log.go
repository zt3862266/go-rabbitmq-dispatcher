package main

import (
	"github.com/natefinch/lumberjack"
	"log"
)

func SetlogFile(file string) {
	log.SetOutput(&lumberjack.Logger{
		Filename:  file,
		MaxSize:   500,
		MaxAge:    30,
		LocalTime: true,
		Compress:  true,
	})
	log.SetFlags(log.LstdFlags)
}

func INFO(format string, v ...interface{}) {
	log.SetPrefix("[INFO]	")
	log.Printf(format, v...)
}

func WARN(format string, v ...interface{}) {
	log.SetPrefix("[WARN]	")
	log.Printf(format, v...)
}
func ERROR(format string, v ...interface{}) {
	log.SetPrefix("[ERROR] ")
	log.Printf(format, v...)
}
