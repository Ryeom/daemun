package router

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Ryeom/daemun/internal"
	"github.com/labstack/echo"
	"io/ioutil"
	"net/http"
)

func Initialize(e *echo.Echo) {
	apis := e.Group("/daemun")
	{
		route(apis)
	}
}
func route(g *echo.Group) {
	g.GET("/healthCheck", healthCheck)
}

func healthCheck(c echo.Context) error {
	result := HttpResult{
		ResultCode: 200,
		ResultMsg:  "[GW]" + internal.HostName() + " OK",
	}

	j, err := json.Marshal(c.Request().Header)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(j))
	rawBody, _ := ioutil.ReadAll(c.Request().Body)
	c.Request().Body = ioutil.NopCloser(bytes.NewBuffer(rawBody))

	fmt.Println(string(rawBody))
	return c.JSON(http.StatusOK, result)
}

type HttpResult struct {
	ResultCode interface{} `json:"resultCode"`
	ResultMsg  string      `json:"resultMsg"`
	ResultData interface{} `json:"resultData,omitempty"`
	Version    string

	RespBody      interface{}
	TransactionId string
	RouterNumber  int
	EncKey        string
}
