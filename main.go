package main

import (
	"context"
	"coze-chat-proxy/common"
	"coze-chat-proxy/config"
	"coze-chat-proxy/logger"
	"coze-chat-proxy/router"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cancel := common.GetCancel()
	defer cancel()

	// Initialize HTTP server
	server := gin.New()
	server.Use(gin.Recovery())

	// 设置路由
	router.SetRouter(server)

	srv := &http.Server{
		Addr:    ":" + config.CONFIG.ServerPort,
		Handler: server,
	}

	// 提示服务启动
	logger.Logger.Info("HTTP server is running on port " + config.CONFIG.ServerPort)

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Logger.Fatal("failed to start HTTP server: " + err.Error())
		}
	}()

	// 等待中断信号
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// 收到信号后取消 context
	cancel()

	// 给 HTTP 服务器一些时间来关闭
	ctxShutDown, cancelShutDown := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutDown()

	if err := srv.Shutdown(ctxShutDown); err != nil {
		logger.Logger.Fatal("HTTP server Shutdown failed:" + err.Error())
	}
}
