package logger

import (
	"log"
	"os"
	"sync"
)

var logger *log.Logger
var once sync.Once

func GetLogger() *log.Logger {
	once.Do(func() {
		logger = log.New(
			os.Stdout,
			"authentication-deanery: ",
			log.Ldate|log.Ltime|log.LUTC|log.Lshortfile,
		)
	})

	return logger
}
