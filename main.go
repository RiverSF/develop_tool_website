package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"develop_tools/internal/config"
	"develop_tools/internal/model"
	"develop_tools/internal/router"
	"develop_tools/pkg/logger"
	"develop_tools/pkg/path"
)

func main() {
	if err := logger.Init(); err != nil {
		log.Fatalf("fail to initLogger, err=%s", err.Error())
	}
	if err := config.Init(); err != nil {
		log.Fatalf("fail to initSetting, err=%s", err.Error())
	}
	if err := model.Init(); err != nil {
		log.Fatalf("fail to initDB, err=%s", err.Error())
	}
	defer model.Close()
	defer logger.Close()

	gin.SetMode(config.AppConfig.RunMode)

	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", config.ServerConfig.ServerPort),
		Handler:        getRouter(),
		ReadTimeout:    time.Duration(config.ServerConfig.ReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(config.ServerConfig.WriteTimeout) * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		logger.Info("listen and serve on 0.0.0.0:%d", config.ServerConfig.ServerPort)
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("fail to listenAndServe, err=%s", err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		log.Fatalf("server shutdown failed, err=%s", err.Error())
	}
}

func getRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())

	router.Init(r)

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})
	r.Static("/assets", path.Join("assets"))
	return r
}
