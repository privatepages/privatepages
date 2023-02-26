package app

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/sync/errgroup"

	"package/main/internal/config"
	"package/main/internal/controllers"
)

var ctx, cancel = context.WithCancel(context.Background())
var group, groupCtx = errgroup.WithContext(ctx)
var conf *config.Config

func init() {
	conf = config.Cfg
}

// Run like main func
func Run() {

	log.Info("Starting app")

	gin.SetMode(gin.TestMode) // DebugMode, ReleaseMode

	r := gin.New()

	r.Use(cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		MaxAge: 12 * time.Hour,
	}))

	r.Use(gin.LoggerWithConfig(gin.LoggerConfig{SkipPaths: []string{"/metrics"}}))

	r.POST("/upload", controllers.Upload)
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	r.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})
	r.GET("/liveness", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})
	r.GET("/readiness", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	server := &http.Server{
		Addr:    conf.HTTPListenIPPort,
		Handler: r,
		// BaseContext: ctx,
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 2 * time.Second,
	}

	group.Go(func() error {
		signalChannel := make(chan os.Signal, 1)
		defer close(signalChannel)
		signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
		select {
		case sig := <-signalChannel:
			log.Errorf("Received signal: %s", sig)
			if err := server.Shutdown(ctx); err != nil {
				log.Errorf("Received an error while shutting down the server: %s", err)
			}
			cancel()
		case <-groupCtx.Done():
			log.Error("Closing signal goroutine")
			if err := server.Shutdown(ctx); err != nil {
				log.Errorf("Received an error while shutting down the server: %s", err)
			}
			return groupCtx.Err()
		}
		return nil
	})

	group.Go(func() error {
		log.Infof("Starting web server on %s", conf.HTTPListenIPPort)
		err := server.ListenAndServe()
		return err
	})

	err := group.Wait()
	if err != nil {
		if errors.Is(err, context.Canceled) {
			log.Error("Context was canceled")
		} else {
			log.Errorf("Received error: %v\n", err)
		}
	} else {
		log.Error("Sucsessfull finished")
	}
}
