package route

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"net/url"

	"github.com/Ryeom/daemun/util/config"
)

type HttpResult struct {
	ResultCode interface{} `json:"resultCode"`
	ResultMsg  string      `json:"resultMsg"`
	ResultData interface{} `json:"resultData,omitempty"`
}

func Initialize(g *gin.Engine, appConfig *config.AppConfig) {
	g.GET("/daemun/healthCheck", func(c *gin.Context) {
		result := HttpResult{
			ResultCode: 200,
			ResultMsg:  "",
		}
		c.JSON(http.StatusOK, result)
	})
	for _, rc := range appConfig.Routes {
		routeConfig := rc

		targetURL, err := url.Parse("http://" + routeConfig.IP)
		if err != nil {
			log.Printf("잘못된 대상 URL %s (key: %s): %v", routeConfig.IP, routeConfig.Key, err)
			continue
		}
		//proxy := httputil.NewSingleHostReverseProxy(targetURL)

		//mux.HandleFunc("/"+routeConfig.Key, func(w http.ResponseWriter, r *http.Request) {
		//	proxy.ServeHTTP(w, r)
		//})
		log.Printf("라우트 등록: /%s -> %s", routeConfig.Key, targetURL.String())
	}

}
