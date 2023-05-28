package database

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/Ryeom/daemun/internal"
	"github.com/Ryeom/daemun/log"
	_ "github.com/go-sql-driver/mysql"
)

var MariaConnection map[string]*MariaConnectionInfo

func InitializeMaria() {
	var err error
	MariaConnection = map[string]*MariaConnectionInfo{}
	MariaConnection[""], err = newMariaConnection("", "", "", "", "")
	if err != nil {
		log.Logger.Error("MariaDB Error", err)
	}
}

type MariaConnectionInfo struct {
	Connection *sql.DB
	Host       string `json:"host"`
	Port       string `json:"port"`
	User       string `json:"user"`
	Password   string `json:"password"`
	Database   string `json:"database"`
}

func newMariaConnection(host, port, user, password, database string) (*MariaConnectionInfo, error) {
	log.Logger.Info(host, port, user, password, database)
	datasource := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, host, port, database)
	log.Logger.Info(datasource)
	mariaConnection, openErr := sql.Open("mysql", datasource)
	if openErr != nil {
		log.Logger.Error("Fail MariaDB Try Connect...", openErr)
		return nil, errors.New("unable to connect MariaDB")
	}
	pingErr := mariaConnection.Ping()
	if pingErr != nil { // 처음 Ping 이 실패했다...
		log.Logger.Error("Fail MariaDB Ping...", pingErr)
		return nil, errors.New("unable to ping MariaDB")
	}
	return &MariaConnectionInfo{mariaConnection, host, port, user, password, database}, nil
}

func (m *MariaConnectionInfo) selectEndpoints() map[string]string {
	if m.Connection == nil {
		return nil
	}

	//tableName := ""
	if !internal.Contains([]string{""}, "") {
		//tableName = "t"
	}

	query := ""
	rows, queryErr := m.Connection.Query(query)
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
