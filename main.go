package main

import (
	"errors"
	"github.com/Ryeom/daemun/internal"
	"github.com/Ryeom/daemun/log"
	"github.com/spf13/viper"
	"os"
	"runtime"
	"strconv"
)

func init() {

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

}

func initializeSetting() error {
	viper.SetConfigName("settings")
	viper.SetConfigType("toml")
	configPath, err := os.Getwd()
	if err != nil {
		log.Logger.Error("work directory ", err)
		err = errors.New("Failed reading work directory. " + err.Error())
		return err
	}
	path := configPath
	viper.AddConfigPath(path)
	log.Logger.Info("Reading configuration from", path)

	err = viper.ReadInConfig()
	if err != nil {
		log.Logger.Error("Failed reading configuration .. ", err)
		err = errors.New("Failed reading configuration. " + err.Error())
		return err
	}

	os := runtime.GOOS
	mark := false
	if os == "darwin" {
		mark = true
		log.Logger.Error(mark, "")
	}
	//
	//for k, v := range viper.AllKeys() {
	//	if strings.HasPrefix(v, ".") {
	//		continue
	//	}
	//	value := viper.GetString(v)
	//	if value == "" {
	//		continue
	//	}
	//
	//}
	// 기타 강제 설정
	ip := internal.GetLocalIP()
	viper.SetDefault("gateway.current-ip", ip)

	return err
}
