package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/salzr/ww2020/cloudrun/pkg/cloudrun/storage"
	"github.com/salzr/ww2020/cloudrun/pkg/version"
)

var svc *storage.Service

func init() {
	log.Println(version.Version)
	var err error
	svc, err = storage.Bootstrap()
	if err != nil {
		log.Fatalf("error bootstrapping storage service, %q", err)
	}
}

func main() {
	r := gin.Default()

	r.GET("/healthz", func(c *gin.Context) {
		c.String(http.StatusOK, version.Version)
	})

	r.POST("/upload", svc.HandleFileUpload)
	log.Fatal(r.Run())
}
