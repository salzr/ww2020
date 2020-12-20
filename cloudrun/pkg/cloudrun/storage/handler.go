package storage

import (
	"context"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

const stateUploaded = "uploaded"

func (s Service) HandleFileUpload(c *gin.Context) {
	uinfo, err := ExtractUserInfo(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	uid := uinfo.Sub

	form, _ := c.MultipartForm()
	files := form.File["upload[]"]
	uploaded := make([]string, 0)

	for _, f := range files {
		so := NewStorageObject(uid, f.Filename)
		filedst := filepath.Join(s.opts.TempDir, so.Filehash)
		err := c.SaveUploadedFile(f, filedst)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		sw := s.gcsClient.Bucket(s.opts.Bucket).Object(filepath.Join(uid, so.Filehash)).NewWriter(context.TODO())
		file, err := os.Open(filedst)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		_, err = io.Copy(sw, file)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		sw.Close()

		if err := so.Persist(s.fsClient); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		uploaded = append(uploaded, f.Filename)
		file.Close()
		os.Remove(filedst)
	}

	c.JSON(http.StatusOK, gin.H{"status": map[string]interface{}{
		"files": uploaded,
		"state": stateUploaded,
	}})
}
