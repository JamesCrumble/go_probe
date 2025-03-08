package main

import (
	"context"
	"fmt"
	"go_probe/probe_kafka"
	"go_probe/probe_psql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	_ "net/http/pprof"
)

func ping(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
func psqlProbe(ctx *gin.Context) {
	probeData, err := probe_psql.Probe()
	if err != nil {
		err = fmt.Errorf("cannot extract data from postgresql: %s", err.Error())
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, probeData)
}
func kafkaProbe(ctx *gin.Context) {
	kafkaCtx, cancel := context.WithTimeout(ctx, time.Second*probe_kafka.TIMEOUT_SECONDS)
	defer cancel()

	probeData, err := probe_kafka.Probe(kafkaCtx)
	if err != nil {
		err = fmt.Errorf("cannot consume from kafka: %s", err.Error())
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, probeData)
}

func configureRouting(router *gin.Engine) {
	router.GET("/ping", ping)

	probeGroup := router.Group("/probe")
	probeGroup.GET("/psql", psqlProbe)
	probeGroup.GET("/kafka", kafkaProbe)
}

func runServer() {
	router := gin.Default()
	configureRouting(router)

	router.Run("0.0.0.0:4000")
}

func main() {
	go http.ListenAndServe("localhost:8080", nil) // prof
	runServer()
}
