package main

import (
	"errors"
	"fmt"
	"github.com/Ryeom/daemun/database"
	"github.com/Ryeom/daemun/internal"
	"github.com/Ryeom/daemun/log"
	"github.com/Ryeom/daemun/router"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"os"
	"strconv"
)

func init() {
	err := internal.InitializeSetting()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}
func setArguments(mode string) error {
	if len(os.Args) != 1 {
		err := errors.New("Process mode is not exists. arguments length : " + strconv.Itoa(len(os.Args)))
		return err
	}
	//if mode != os.Args[1] {
	//	err := errors.New("Process mode is not match. : " + strconv.Itoa(len(os.Args)))
	//}
	return nil
}

func main() {
	//runtime.GOMAXPROCS(runtime.NumCPU())
	log.InitializeApplicationLog()
	database.InitializeMaria()
	e := echo.New()
	e.Use(middleware.CORS())
	e.Use(middleware.Recover())
	e.Use(middleware.LoggerWithConfig(log.GetCustomLogConfig()))
	router.Initialize(e)

	log.Logger.Fatal(e.Start(":8080"))
}
