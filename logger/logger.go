package logger

import (
	"log"
	"os"
)

var (
	Log *log.Logger
)

// func init() {
// 	// set location of log file
// 	var logpath = "logs/logs.log"

// 	flag.Parse()
// 	var file, err1 = os.Create(logpath)

// 	if err1 != nil {
// 		panic(err1)
// 	}
// 	Log = log.New(file, "", log.LstdFlags|log.Lshortfile)
// 	Log.Println("LogFile : " + logpath)
// }
var (
	WarningLogger *log.Logger
	InfoLogger    *log.Logger
	ErrorLogger   *log.Logger
)

func init() {
	var logpath = "logs/logs.log"
	file, err := os.OpenFile(logpath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	InfoLogger = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarningLogger = log.New(file, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}
