package route

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/Ryeom/daemun/util/config"
)

type Router struct { // 커스텀 HTTP 라우터
	mux *http.ServeMux
}

func NewRouter(appConfig *config.AppConfig) *Router { //AppConfig에 정의된 라우트 정보를 바탕으로 라우터를 초기화
	mux := http.NewServeMux()

	// 각 라우트 설정마다 reverse proxy를 생성하여 등록 (TODO : 추후 변경)
	for _, rc := range appConfig.Routes {
		routeConfig := rc

		targetURL, err := url.Parse("http://" + routeConfig.IP)
		if err != nil {
			log.Printf("잘못된 대상 URL %s (key: %s): %v", routeConfig.IP, routeConfig.Key, err)
			continue
		}
		proxy := httputil.NewSingleHostReverseProxy(targetURL)

		mux.HandleFunc("/"+routeConfig.Key, func(w http.ResponseWriter, r *http.Request) {
			proxy.ServeHTTP(w, r)
		})
		log.Printf("라우트 등록: /%s -> %s", routeConfig.Key, targetURL.String())
	}

	// 매칭되지 않은 경로에 대한 기본 핸들러
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})

	return &Router{mux: mux}
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}
