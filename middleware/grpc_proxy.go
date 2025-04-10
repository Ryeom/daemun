package middleware

import (
	"context"
	"net/http"
	"time"

	logger "github.com/Ryeom/daemun/log"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

func GRPCProxyMiddleware(target string) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()

		// TODO : 운영환경에서 연결 풀이나 재사용 고려
		conn, err := grpc.DialContext(ctx, target, grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			c.String(http.StatusInternalServerError, "gRPC 연결 실패")
			logger.ServerLogger.Printf("gRPC 연결 실패 (%s): %v", target, err)
			c.Abort()
			return
		}
		defer conn.Close()

		// TODO : 아래 내용 서비스에 적용하기~
		// client := proto.NewYourServiceClient(conn)
		// grpcReq := &proto.YourRequest{
		//     // HTTP 요청의 파라미터, Body, Header 등을 기반으로 값 설정
		// }
		// grpcResp, err := client.YourMethod(ctx, grpcReq)
		// if err != nil {
		//     c.String(http.StatusInternalServerError, "gRPC 호출 실패")
		//     logger.ServerLogger.Printf("gRPC 호출 실패: %v", err)
		//     c.Abort()
		//     return
		// }
		// c.JSON(http.StatusOK, grpcResp)

		c.String(http.StatusOK, "gRPC Proxy 호출 성공: %s", c.Request.URL.Path)
		c.Abort()
	}
}
