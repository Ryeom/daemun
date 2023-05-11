package internal

import (
	"errors"
	"fmt"
	"github.com/Ryeom/daemun/log"
	"github.com/spf13/viper"
	"os"
	"strings"
)

func InitializeSetting() error {
	viper.SetConfigName("settings")
	viper.SetConfigType("toml")
	configPath, err := os.Getwd()
	if err != nil {
		err = errors.New("Failed reading work directory. " + err.Error())
		return err
	}
	path := configPath
	viper.AddConfigPath(path)
	err = viper.ReadInConfig()
	if err != nil {
		err = errors.New("Failed reading configuration. " + err.Error())
		return err
	}

	if viper.GetString(log.ProjectName+".key") != "" {
		aeskey := viper.GetString(log.ProjectName + ".key")
		for _, v := range viper.AllKeys() {
			if strings.HasPrefix(v, log.ProjectName+".") {
				continue
			}
			originalValue, decErr := DecryptAES(viper.GetString(v), []byte(aeskey))
			if decErr != nil {
				log.Logger.Error("Failed reading configuration value Decrypt .. ", err)
				err = errors.New("Failed reading configuration value Decrypt. " + err.Error())
				return err
			}
			viper.Set(v, originalValue)
		}
	}
	// 기타 강제 설정
	ip := GetLocalIP()
	viper.SetDefault("gateway.current-ip", ip)

	return err
}

func PrintAllLog() {
	if IsOperationMode() {
		for _, v := range viper.AllKeys() {
			if strings.HasPrefix(v, ".") {
				continue
			}
			value := viper.GetString(v)
			if value == "" {
				continue
			} else {
				if log.Logger != nil {
					log.Logger.Info(v, ":", value)
				} else {
					fmt.Println(v, ":", value)
				}
			}

		}
	}
}

func GetMode() string {
	return viper.GetString("daemun.mode")
}
