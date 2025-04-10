package route

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/Ryeom/daemun/util/config"
	"github.com/gin-gonic/gin"
)

type HttpResult struct {
	ResultCode interface{} `json:"resultCode"`
	ResultMsg  string      `json:"resultMsg"`
	ResultData interface{} `json:"resultData,omitempty"`
}

func Initialize(g *gin.Engine, appConfig *config.AppConfig) {
	/* 1. health check */
	g.GET("/daemun/healthCheck", func(c *gin.Context) {
		result := HttpResult{
			ResultCode: 200,
			ResultMsg:  "OK",
		}
		c.JSON(http.StatusOK, result)
	})

	/* 2. others (proxy route) */
	for _, rc := range appConfig.Routes {
		routeConfig := rc

		targetURL, err := url.Parse("http://" + routeConfig.IP)
		if err != nil {
			log.Printf("잘못된 대상 URL %s (key: %s): %v", routeConfig.IP, routeConfig.Key, err)
			continue
		}

		proxy := httputil.NewSingleHostReverseProxy(targetURL)

		endpointPath := "/" + routeConfig.Key
		g.Any(endpointPath, gin.WrapH(proxy))

		log.Printf("라우트 등록: %s -> %s", endpointPath, targetURL.String())
	}

	/* 2. self request */
	g.GET("/daemun/endpoints", GetEndpointConfigsHandler)
	g.POST("/daemun/endpoints", CreateEndpointConfigHandler)
	g.PUT("/daemun/endpoints", UpdateEndpointConfigHandler)
	g.DELETE("/daemun/endpoints", DeleteEndpointConfigHandler)

	g.GET("/daemun/limit-policy", GetLimitPolicyHandler)
	g.POST("/daemun/limit-policy", CreateLimitPolicyHandler)
	g.PUT("/daemun/limit-policy", UpdateLimitPolicyHandler)
	g.DELETE("/daemun/limit-policy", DeleteLimitPolicyHandler)
}
