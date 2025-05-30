package main

import (
	"fmt"
	"go_probe/probe_kafka"
	"go_probe/probe_psql"
	"go_probe/probe_redis"
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
			err = fmt.Errorf("cannot present %s probe: %+v", probe.Name(), err)
			ctx.AbortWithError(http.StatusInternalServerError, err)
		}
		ctx.JSON(http.StatusOK, probeData)
	}
}

func configureRouting(router *gin.Engine) {
	router.GET("/ping", ping)

	probeGroup := router.Group("/probe")
	declareProbe(probeGroup, probe_psql.Realization())
	declareProbe(probeGroup, probe_kafka.Realization())
	declareProbe(probeGroup, probe_redis.Realization())

}

func runServer() {
	router := gin.Default()
	configureRouting(router)

	router.Run(config.ApiAddr())
}

func main() {
	if config.PROF_ENABLE {
		go http.ListenAndServe(config.ProfAddr(), nil) // prof
	}
	runServer()
}
