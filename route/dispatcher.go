package route

import (
	"context"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"google.golang.org/grpc"

	logger "github.com/Ryeom/daemun/log"
	// "daemun/proto" // 실제 proto 패키지 경로로 변경 필요
)

type DestinationConfig struct {
	Protocol string // "http" 또는 "grpc"
	URL      string // 예: "http://localhost:9001" 또는 "localhost:50051" (grpc)
}

func dispatchHTTP(w http.ResponseWriter, r *http.Request, dest DestinationConfig) {
	parsedURL, err := url.Parse(dest.URL)
	if err != nil {
		http.Error(w, "잘못된 대상 URL", http.StatusInternalServerError)
		logger.ServerLogger.Printf("잘못된 대상 URL: %v", err)
		return
	}
	proxy := httputil.NewSingleHostReverseProxy(parsedURL)
	proxy.ServeHTTP(w, r)
}

func dispatchGRPC(w http.ResponseWriter, r *http.Request, dest DestinationConfig) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, dest.URL, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		http.Error(w, "gRPC 연결 실패", http.StatusInternalServerError)
		logger.ServerLogger.Printf("gRPC 연결 실패: %v", err)
		return
	}
	defer conn.Close()

	// TODO: gRPC 클라이언트로 HTTP 요청을 gRPC 요청으로 변환하여 호출
	// client := proto.NewYourServiceClient(conn)
	// grpcReq := &proto.YourRequest{ ... }
	// grpcResp, err := client.YourMethod(ctx, grpcReq)
	// if err != nil { ... }
	// w.Header().Set("Content-Type", "application/json")
	// io.WriteString(w, grpcResp.String())

	// 예시: 단순 성공 메시지 반환
	w.Header().Set("Content-Type", "text/plain")
	io.WriteString(w, "gRPC 호출 성공: "+r.URL.Path)
}

// DispatchRequest는 DestinationConfig의 Protocol에 따라 HTTP 또는 gRPC 호출을 수행합니다.
func DispatchRequest(w http.ResponseWriter, r *http.Request, dest DestinationConfig) {
	if dest.Protocol == "grpc" {
		dispatchGRPC(w, r, dest)
	} else {
		dispatchHTTP(w, r, dest)
	}
}
