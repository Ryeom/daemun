package route

import (
	"context"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	logger "github.com/Ryeom/daemun/log"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

type DestinationConfig struct {
	Protocol string
	URL      string
}

func dispatchHTTP(c *gin.Context, dest DestinationConfig) {
	parsedURL, err := url.Parse(dest.URL)
	if err != nil {
		c.String(http.StatusInternalServerError, "잘못된 대상 URL")
		logger.ServerLogger.Printf("잘못된 대상 URL %s: %v", dest.URL, err)
		return
	}
	proxy := httputil.NewSingleHostReverseProxy(parsedURL)
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		logger.ServerLogger.Printf("Reverse proxy error: %v", err)
		http.Error(w, "Proxy Error", http.StatusBadGateway)
	}

	proxy.ServeHTTP(c.Writer, c.Request)
}

func dispatchGRPC(c *gin.Context, dest DestinationConfig) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	// TODO : 연결 재사용/풀링 고려하기
	conn, err := grpc.DialContext(ctx, dest.URL, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		c.String(http.StatusInternalServerError, "gRPC 연결 실패")
		logger.ServerLogger.Printf("gRPC 연결 실패 (%s): %v", dest.URL, err)
		return
	}
	defer conn.Close()

	// TODO: 실제 gRPC 클라이언트를 사용하여 HTTP 요청을 gRPC 요청으로 변환 후 호출
	// 예를 들어, proto 파일로부터 생성된 클라이언트를 사용하여 요청을 전송하고 응답을 처리합니다.
	// client := proto.NewYourServiceClient(conn)
	// grpcReq := &proto.YourRequest{...} // 변환 로직 구현 필요
	// grpcResp, err := client.YourMethod(ctx, grpcReq)
	// if err != nil { ... }

	c.String(http.StatusOK, "gRPC 호출 성공: %s", c.Request.URL.Path)
}

func DispatchRequest(c *gin.Context, dest DestinationConfig) {
	if dest.Protocol == "grpc" {
		dispatchGRPC(c, dest)
	} else {
		dispatchHTTP(c, dest)
	}
}

func GinDispatchMiddleware(dest DestinationConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		DispatchRequest(c, dest)
		c.Abort()
	}
}
