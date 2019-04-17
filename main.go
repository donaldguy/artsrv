package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/donaldguy/artsrv/art"
	"github.com/donaldguy/artsrv/requests"
	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	prom "github.com/prometheus/client_golang/prometheus"
	"go.opencensus.io/exporter/prometheus"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/stats/view"
)

func setupRoutes(r *gin.Engine) {
	r.POST("/image", requests.RegisterHandler)
	r.POST("/image/:id/chunks", requests.UploadChunkHandler)

	r.GET("/image/:id", requests.DownloadHandler)
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

func setupJSONLogs(r *gin.Engine) {
	zerolog.SetGlobalLevel(zerolog.WarnLevel)

	if gin.IsDebugging() {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Logger = log.Output(
			zerolog.ConsoleWriter{
				Out:     os.Stderr,
				NoColor: false,
			},
		)
	}

	r.Use(logger.SetLogger(logger.Config{
		SkipPath: []string{"/metrics"},
	}))
}

func main() {
	dbPath := flag.String("db", "/data/art.db", "Path where to create or read a BoltDB file")
	bind := flag.String("bind", ":4444", "a [<host>]:<port> expression for what IP, port to bind to")

	flag.Parse()

	art.InitBoltStorage(*dbPath)

	r := gin.New()

	setupJSONLogs(r)
	setupMetrics(r)

	setupRoutes(r)

	http.ListenAndServe(
		*bind,
		&ochttp.Handler{Handler: r},
	)
}
