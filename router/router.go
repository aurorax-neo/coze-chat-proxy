package router

import (
	"coze-chat-proxy/middleware"
	"coze-chat-proxy/v1/chat"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SetRouter(router *gin.Engine) {

	router.GET("/", Index)
	v1Router := router.Group("/v1")
	v1Router.Use(middleware.V1Cors)
	v1Router.Use(middleware.V1Response)
	v1Router.Use(middleware.V1Auth)
	v1Router.POST("/chat/completions", chat.Completions)
}

func Index(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello world! This is coze-chat-proxy.",
	})
}
