package logger

import (
	"log"
	"os"
	"strings"
)

var AccessLogger *log.Logger
var SecurityLogger *log.Logger
var ServerLogger *log.Logger
var TraceLogger *log.Logger

func InitLoggers() error {
	if _, err := os.Stat("log"); os.IsNotExist(err) {
		if err := os.Mkdir("log", 0755); err != nil {
			return err
		}
	}

	accessFile, err := os.OpenFile("log/access.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	AccessLogger = log.New(accessFile, "ACCESS: ", log.Ldate|log.Ltime|log.Lshortfile)

	securityFile, err := os.OpenFile("log/security.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	SecurityLogger = log.New(securityFile, "SECURITY: ", log.Ldate|log.Ltime|log.Lshortfile)

	serverFile, err := os.OpenFile("log/server.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	ServerLogger = log.New(serverFile, "SERVER: ", log.Ldate|log.Ltime|log.Lshortfile)

	traceFile, err := os.OpenFile("log/trace.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	TraceLogger = log.New(traceFile, "TRACE: ", log.Ldate|log.Ltime|log.Lshortfile)

	var banner = `
` + strings.Repeat("▅", 150) + `
     _                                                           
    | |                                     _               _    
  _ | | ____ ____ ____  _   _ ____      ___| |_  ____  ____| |_  
 / || |/ _  ) _  |    \| | | |  _ \    /___)  _)/ _  |/ ___)  _) 
( (_| ( (/ ( ( | | | | | |_| | | | |  |___ | |_( ( | | |   | |__ 
 \____|\____)_||_|_|_|_|\____|_| |_|  (___/ \___)_||_|_|    \___)
` + strings.Repeat("░", 150)
	AccessLogger.Println(banner)
	SecurityLogger.Println(banner)
	ServerLogger.Println(banner)
	TraceLogger.Println(banner)

	return nil
}
