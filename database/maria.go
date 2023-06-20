package database

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/Ryeom/daemun/internal"
	"github.com/Ryeom/daemun/log"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
)

var MariaConnection map[string]*ConnectionInfo

func InitializeMaria() {
	MariaConnection = map[string]*ConnectionInfo{}
	c := ConnectionInfo{
		Host:     viper.GetString("maria.host"),
		Port:     viper.GetString("maria.port"),
		User:     viper.GetString("maria.user"),
		Password: viper.GetString("maria.password"),
		Database: viper.GetString("maria.database"),
	}
	err := c.newConnection()
	if err != nil {
		log.Logger.Error("MariaDB Error", err)
	}
}

type ConnectionInfo struct {
	Connection *sql.DB
	Host       string `json:"host"`
	Port       string `json:"port"`
	User       string `json:"user"`
	Password   string `json:"password"`
	Database   string `json:"database"`
}

func (c *ConnectionInfo) newConnection() error {
	log.Logger.Info(c.User, c.Password, c.Host, c.Port, c.Database)
	datasource := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", c.User, c.Password, c.Host, c.Port, c.Database)
	log.Logger.Info(datasource)
	mariaConnection, openErr := sql.Open("mysql", datasource)
	if openErr != nil {
		log.Logger.Error("Fail MariaDB Try Connect...", openErr)
		return errors.New("unable to connect MariaDB")
	}
	pingErr := mariaConnection.Ping()
	if pingErr != nil {
		log.Logger.Error("Fail MariaDB Ping...", pingErr)
		return errors.New("unable to ping MariaDB")
	}
	return nil
}

func (c *ConnectionInfo) selectEndpoints() map[string]string {
	if c.Connection == nil {
		log.Logger.Error("Disconnected MariaDB... connection is nil.")
		return nil
	}

	//tableName := ""
	if !internal.Contains([]any{""}, "") {
		//tableName = "t"
	}

	query := "SELECT * FROM ? WHERE "
	rows, queryErr := c.Connection.Query(query)
	if queryErr != nil {
		panic("")
	}
	defer rows.Close()
	result := map[string]string{}

	for rows.Next() {
		var e, h string
		scanErr := rows.Scan(&e, &h)
		fmt.Println(e, h)
		if scanErr != nil {
			log.Logger.Error("Fail Scan selected data...", scanErr)
		}
		result[e] = h
	}
	return result
}
