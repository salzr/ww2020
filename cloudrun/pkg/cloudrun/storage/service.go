package storage

import (
	"context"

	"cloud.google.com/go/storage"
)

type Service struct {
	gcsClient *storage.Client
	opts      *Options
}

func NewStorageService(opts *Options) (*Service, error) {
	client, err := storage.NewClient(context.Background())
	if err != nil {
		return nil, err
	}
	return &Service{gcsClient: client, opts: opts}, nil
}

func (s *Service) validate(name string) error {
	bkt := s.gcsClient.Bucket(name)
	_, err := bkt.Attrs(context.Background())
	if err != nil {
		return err
	}
	return err
}
