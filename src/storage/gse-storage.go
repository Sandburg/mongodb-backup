package storage

import (
	"context"
	"os"

	"cloud.google.com/go/storage"
	log "github.com/sirupsen/logrus"

	firebase "firebase.google.com/go"
	"github.com/pkg/errors"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

type GseStorage struct {
	ctx    context.Context
	bucket *storage.BucketHandle
}

func New(bucketID string) GseStorage {
	ctx := context.Background()
	b := initbucket(ctx, bucketID)
	return GseStorage{
		ctx:    ctx,
		bucket: b,
	}
}

func (fs GseStorage) WriteFile(fileName string, content []byte) {
	wc := fs.bucket.Object(fileName).NewWriter(fs.ctx)

	if _, err := wc.Write(content); err != nil {
		log.Panic(errors.Wrapf(err, "unable to write data to bucket, file %q", fileName))
	}
	defer wc.Close()
}

func (fs GseStorage) GetFilenames() (names []string) {
	it := fs.bucket.Objects(fs.ctx, nil)
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Panic(err)
		}
		names = append(names, attrs.Name)
	}

	return names
}

func (fs GseStorage) DeleteFile(fileName string) {
	err := fs.bucket.Object(fileName).Delete(fs.ctx)
	if err != nil {
		log.Panic(errors.Wrapf(err, "error while deleting file: %s", fileName))
	}
}

func initbucket(ctx context.Context, bucketID string) *storage.BucketHandle {
	app := initApp(ctx)
	s, err := app.Storage(ctx)
	if err != nil {
		log.Panic(errors.Wrap(err, "error while inizializing GseStorage client"))
	}
	b, err := s.Bucket(bucketID)
	if err != nil {
		log.Panic(errors.Wrapf(err, "error while inizializing bucket: %s", bucketID))
	}
	return b
}

func initApp(ctx context.Context) *firebase.App {
	keyfile := os.Getenv("GCS_KEY_FILE_PATH")
	if keyfile == "" {
		log.Panic("path to keyfile must be set as environment variable GCS_KEY_FILE_PATH")
	}
	opt := option.WithCredentialsFile(os.Getenv("GCS_KEY_FILE_PATH"))
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Errorf("error initializing app: %v\n", err)
	}
	return app
}
