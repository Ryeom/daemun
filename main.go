package main

import (
	"daemun/log"
	"daemun/util"
	"errors"
	"github.com/spf13/viper"
	"os"
	"runtime"
	"strings"
)

func init() {

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
	}

	for k, v := range viper.AllKeys() {
		if strings.HasPrefix(v, ".") {
			continue
		}
		value := viper.GetString(v)
		if value == "" {
			continue
		}

	}
	// 기타 강제 설정
	ip := util.GetLocalIP()
	viper.SetDefault("gateway.current-ip", ip)

	return err
}
