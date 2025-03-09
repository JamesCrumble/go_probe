package main

import (
	"fmt"
	"go_probe/probe_kafka"
	"go_probe/probe_psql"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "net/http/pprof"
)

type Probe interface {
	Name() string
	Present(ctx *gin.Context) (any, error)
}

func ping(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func declareProbe(router gin.IRouter, probe Probe) {
	router.GET("/"+probe.Name(), ProbeHandler(probe))
}

func ProbeHandler(probe Probe) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		probeData, err := probe.Present(ctx)
		if err != nil {
			err = fmt.Errorf("cannot present %s probe: %s", probe.Name(), err.Error())
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		ctx.JSON(http.StatusOK, probeData)
	}
}

func configureRouting(router *gin.Engine) {
	router.GET("/ping", ping)

	probeGroup := router.Group("/probe")
	declareProbe(probeGroup, probe_psql.Realization())
	declareProbe(probeGroup, probe_kafka.Realization())

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
