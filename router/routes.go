package router

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Ryeom/daemun/internal"
	"github.com/Ryeom/daemun/log"
	"github.com/labstack/echo"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	BASIC = iota
	CERT
	UNCERT
)

func Initialize(e *echo.Echo) {
	apis := e.Group("/gw/admin")
	{
		gatewayRoute(apis)
	}
	list := setSolutionGroup()
	for k, v := range list {
		if len(v) < 1 {
			log.Logger.Info("Undefined Functions ... ", k)
			continue
		}
		log.Logger.Info("API Initialize ... ", k)
		e.Group(k, v...)
	}
}

func setSolutionGroup() map[string][]echo.MiddlewareFunc {
	switch 0 {
	case BASIC:
		return map[string][]echo.MiddlewareFunc{}
	case CERT:
		return map[string][]echo.MiddlewareFunc{}
	case UNCERT:
		return map[string][]echo.MiddlewareFunc{}
	default:
		return nil
	}
}
func gatewayRoute(g *echo.Group) {
	g.GET("/health-check", healthCheck)
}

func healthCheck(c echo.Context) error {
	result := HttpResult{
		Code: "",
		Msg:  "[GW]" + internal.HostName() + " OK",
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
	Result        DestinationResult/* 목적지 반환 결과 */ `json:"result"`
	Code          string `json:"code"`
	Msg           string `json:"msg"`
	TransactionId string
	RouterNumber  int/* gateway number */ `json:"transactionId"`
	EncKey        string/* 암호화키 */ `json:"routerNumber"`
	RequestTime   time.Time/* 요청 시간 */ `json:"encKey"`
	ResponseTime  time.Time/* 응답 시간 */ `json:"requestTime"`
	ClientIP      string/* 요청지 IP */ `json:"responseTime"`
	/* 소요 시간 */

}

type DestinationResult struct {
	ResultCode interface{} `json:"resultCode"` /* 반환 코드 */
	ResultMsg  string      `json:"resultMsg"`
	ResultData interface{} `json:"resultData,omitempty"`
	Version    string      /* 목적지 version */
}
