package route

import (
	"context"
	"net/http"
	"time"

	logger "github.com/Ryeom/daemun/log"
	mongoutil "github.com/Ryeom/daemun/mongo"
	"github.com/gin-gonic/gin"
)

func GetEndpointConfigsHandler(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	configs, err := mongoutil.GetAllEndpointConfigs(ctx, mongoutil.Client)
	if err != nil {
		logger.ServerLogger.Printf("EndpointConfigs 조회 실패: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, configs)
}

func CreateEndpointConfigHandler(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var config mongoutil.EndpointConfig
	if err := c.BindJSON(&config); err != nil {
		logger.ServerLogger.Printf("EndpointConfig 생성 - 디코딩 오류: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := mongoutil.CreateEndpointConfig(ctx, mongoutil.Client, config); err != nil {
		logger.ServerLogger.Printf("EndpointConfig 생성 실패: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusCreated)
}

func UpdateEndpointConfigHandler(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	id := c.Param("id")
	var updateData map[string]interface{}
	if err := c.BindJSON(&updateData); err != nil {
		logger.ServerLogger.Printf("EndpointConfig 수정 - 디코딩 오류: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := mongoutil.UpdateEndpointConfig(ctx, mongoutil.Client, id, updateData); err != nil {
		logger.ServerLogger.Printf("EndpointConfig 수정 실패 (id: %s): %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

func DeleteEndpointConfigHandler(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	id := c.Param("id")
	if err := mongoutil.DeleteEndpointConfig(ctx, mongoutil.Client, id); err != nil {
		logger.ServerLogger.Printf("EndpointConfig 삭제 실패 (id: %s): %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}
