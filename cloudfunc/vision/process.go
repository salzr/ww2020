package vision

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/functions/metadata"
	vision "cloud.google.com/go/vision/apiv1"
	"github.com/salzr/ww2020/cloudrun/pkg/cloudrun/storage"
	visionpb "google.golang.org/genproto/googleapis/cloud/vision/v1"
)

const (
	logInfoPattern         = "%s [GoogleCloudStorageEvent] %s\n"
	objectLookupKeyPattern = "Users/%s/Object/%s"

	logError = "E"
	logInfo  = "I"
	logWarn  = "W"
)

type GCSEvent struct {
	Bucket string `json:"bucket"`
	Name   string `json:"name"`
}

type Logger struct {
	pattern string
}

var (
	visionClient *vision.ImageAnnotatorClient
	fstoreClient *firestore.Client

	logger *Logger
)

func init() {
	var err error
	logger = newLogger()
	visionClient, err = vision.NewImageAnnotatorClient(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fstoreClient, err = firestore.NewClient(context.Background(), os.Getenv("PROJECT_ID"))
	if err != nil {
		log.Fatal(err)
	}
}

func ProcessEvent(ctx context.Context, e GCSEvent) error {
	meta, err := metadata.FromContext(ctx)
	if err != nil {
		return logger.LogWithError(fmt.Sprintf("metadata.FromContext: %v", err))
	}
	logger.Info(fmt.Sprintf("eventType=[%s] eventID=[%s]", meta.EventType, meta.EventID))
	image := &visionpb.Image{
		Source: &visionpb.ImageSource{
			GcsImageUri: fmt.Sprintf("gs://%s/%s", e.Bucket, e.Name),
		},
	}
	annotations, err := visionClient.DetectLabels(ctx, image, nil, 10)
	if err != nil {
		return logger.LogWithError(fmt.Sprintf("visionClient.DetectLabels: %v", err.Error()))
	}
	vmeta := make([]storage.VisionMeta, 0)
	for _, a := range annotations {
		vmeta = append(vmeta, storage.VisionMeta{Description: a.Description, Score: a.Score})
	}
	if len(vmeta) > 0 {
		doc := fstoreClient.Doc(extractLookupKey(e.Name))
		wr, err := doc.Update(context.Background(), []firestore.Update{
			{Path: "metadata", Value: storage.Metadata{
				VisionMeta: vmeta,
			}},
		})
		if err != nil {
			logger.LogWithError(fmt.Sprintf("firestore.Update: %v", err.Error()))
		}
		logger.Info(fmt.Sprintf("%s: %s", extractLookupKey(e.Name), wr.UpdateTime))
	}
	return nil
}

func extractLookupKey(path string) (key string) {
	parts := strings.Split(path, "/")
	user, obj := parts[0], parts[1]
	key = fmt.Sprintf(objectLookupKeyPattern, user, obj)
	return
}

func newLogger() *Logger {
	return &Logger{logInfoPattern}
}

func (l *Logger) Info(m string) {
	log.Printf(logInfoPattern, logInfo, m)
}

func (l *Logger) Error(m string) {
	log.Printf(logInfoPattern, logError, m)
}

func (l *Logger) Warn(m string) {
	log.Printf(logInfoPattern, logWarn, m)
}

func (l *Logger) LogWithError(m string) error {
	l.Error(m)
	return errors.New(m)
}
