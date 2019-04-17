package main

import (
	"flag"

	"github.com/donaldguy/artsrv/art"
	"github.com/donaldguy/artsrv/requests"
	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.POST("/image", requests.RegisterHandler)
	r.POST("/image/:id/chunks", requests.UploadChunkHandler)

	r.GET("/image/:id", requests.DownloadHandler)

	return r
}

func main() {
	dbPath := flag.String("db", "/data/art.db", "Path where to create or read a BoltDB file")
	bind := flag.String("bind", ":4444", "a [<host>]:<port> expression for what IP, port to bind to")
	flag.Parse()

	art.InitBoltStorage(*dbPath)

	r := setupRouter()

	r.Run(*bind)
}
