package storage

import (
	"context"
	"crypto/sha1"
	"fmt"
	"mime"
	"path/filepath"

	"cloud.google.com/go/firestore"
)

type StorageObject struct {
	Uid      string     `firestore:"uid"`
	Filename string     `firestore:"filename"`
	Filehash string     `firestore:"filehash"`
	MimeType string     `firestore:"mimeType,omitempty"`
	Metadata []Metadata `firestore:"metadata,omitempty"`
}

type Metadata struct {
	VisionMeta []VisionMeta `firestore:"visionMeta"`
}

type VisionMeta struct {
	Description string  `firestore:"description"`
	Score       float32 `firestore:"score"`
}

func NewStorageObject(uid, filename string) *StorageObject {
	return &StorageObject{
		Uid:      uid,
		Filename: filename,
		Filehash: hashFileName(filename),
		MimeType: mime.TypeByExtension(filepath.Ext(filename)),
	}
}

func (so *StorageObject) Persist(c *firestore.Client) error {
	key := fmt.Sprintf("Users/%s/Object/%s", so.Uid, so.Filehash)
	_, err := c.Doc(key).Set(context.Background(), so)
	return err
}

func hashFileName(s string) (b string) {
	h := sha1.New()
	h.Write([]byte(s))
	b = fmt.Sprintf("%x", h.Sum(nil))
	return
}
