package storage

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	stateUploaded = "uploaded"
)

func (s Service) HandleFileUpload(c *gin.Context) {
	form, _ := c.MultipartForm()
	files := form.File["upload[]"]
	uploaded := make([]string, 0)

	for _, f := range files {
		id := uuid.New()
		filedst := filepath.Join(s.opts.TempDir, fmt.Sprintf("%s-%s", id, f.Filename))
		err := c.SaveUploadedFile(f, filedst)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		uploaded = append(uploaded, f.Filename)
		os.Remove(filedst)
	}

	c.JSON(http.StatusOK, gin.H{"status": map[string]interface{}{
		"files": uploaded,
		"state": stateUploaded,
	}})
}
