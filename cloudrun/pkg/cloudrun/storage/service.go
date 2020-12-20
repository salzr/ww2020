package storage

import (
	"context"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/storage"
)

type Service struct {
	gcsClient *storage.Client
	fsClient  *firestore.Client
	opts      *Options
}

func NewStorageService(opts *Options) (*Service, error) {
	gcsClient, err := storage.NewClient(context.Background())
	if err != nil {
		return nil, err
	}

	dbClient, err := firestore.NewClient(context.Background(), opts.ProjectId)
	if err != nil {
		return nil, err
	}

	return &Service{gcsClient: gcsClient, fsClient: dbClient, opts: opts}, nil
}

func (s *Service) validate(name string) error {
	bkt := s.gcsClient.Bucket(name)
	_, err := bkt.Attrs(context.Background())
	if err != nil {
		return err
	}
	return err
}
