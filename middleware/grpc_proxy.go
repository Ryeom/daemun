package middleware

import (
	"context"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"google.golang.org/grpc"
)

// GRPCProxyMiddleware : HTTP 요청이 "/grpc/"로 시작 -> gRPC 호출 수행
func GRPCProxyMiddleware(target string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.URL.Path, "/grpc/") {
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()

				conn, err := grpc.DialContext(ctx, target, grpc.WithInsecure(), grpc.WithBlock())
				if err != nil {
					http.Error(w, "gRPC 연결 실패", http.StatusInternalServerError)
					log.Printf("gRPC 연결 실패: %v", err)
					return
				}
				defer conn.Close()

				// client := proto.NewYourServiceClient(conn)

				w.Header().Set("Content-Type", "text/plain")
				io.WriteString(w, "gRPC Proxy 호출 성공: "+r.URL.Path)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
