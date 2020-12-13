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
	r := gin.New()

	r.GET("/healthz", func(c *gin.Context) {
		c.String(http.StatusOK, version.Version)
	})

	r.GET("/version", func(c *gin.Context) {
		uinfo, _ := storage.ExtractUserInfo(c)
		c.JSON(http.StatusOK, gin.H{
			"version": version.Version,
			"uid":     uinfo.Sub,
			"code":    http.StatusOK,
		})
	})

	r.POST("/upload", svc.HandleFileUpload)

	log.Fatal(r.Run())
}
