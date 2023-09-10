package main

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var requestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
	Name: prometheus.BuildFQName("", "http_server_requests", "seconds"),
}, []string{"status", "method", "uri"})

func NewMetricsMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()

		request := ctx.Request
		ctx.Next()

		status := ctx.Writer.Status()
		statusCode := strconv.Itoa(status)
		requestDuration.WithLabelValues(statusCode, request.Method, request.RequestURI).Observe(time.Now().Sub(start).Seconds())
	}

}
