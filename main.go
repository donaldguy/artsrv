package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/donaldguy/artsrv/art"
	"github.com/donaldguy/artsrv/requests"
	"github.com/gin-gonic/gin"

	prom "github.com/prometheus/client_golang/prometheus"
	"go.opencensus.io/exporter/prometheus"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/stats/view"
)

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.POST("/image", requests.RegisterHandler)
	r.POST("/image/:id/chunks", requests.UploadChunkHandler)

	r.GET("/image/:id", requests.DownloadHandler)

	return r
}

func setupMetrics(r *gin.Engine) {
	pe, err := prometheus.NewExporter(prometheus.Options{
		Registry: prom.DefaultGatherer.(*prom.Registry),
	})
	if err != nil {
		panic(fmt.Errorf("Initializing Prometheus: %v", err))
	}

	view.RegisterExporter(pe)

	err = view.Register(
		ochttp.ServerRequestCountView,
		ochttp.ServerRequestBytesView,
		ochttp.ServerResponseBytesView,
		ochttp.ServerLatencyView,
		ochttp.ServerRequestCountByMethod,
		ochttp.ServerResponseCountByStatusCode,
	)
	if err != nil {
		panic(fmt.Errorf("Registring ochttp metrics: %v", err))
	}

	r.GET("/metrics", gin.HandlerFunc(func(c *gin.Context) {
		pe.ServeHTTP(c.Writer, c.Request)
	}))
}

func main() {
	dbPath := flag.String("db", "/data/art.db", "Path where to create or read a BoltDB file")
	bind := flag.String("bind", ":4444", "a [<host>]:<port> expression for what IP, port to bind to")
	flag.Parse()

	art.InitBoltStorage(*dbPath)

	r := setupRouter()
	setupMetrics(r)

	http.ListenAndServe(
		*bind,
		&ochttp.Handler{Handler: r},
	)
}
