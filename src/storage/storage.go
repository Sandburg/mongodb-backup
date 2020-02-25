package storage

import (
	"context"
	"os"

	log "github.com/sirupsen/logrus"

	"cloud.google.com/go/storage"
	firebase "firebase.google.com/go"
	"github.com/pkg/errors"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

type Storage struct {
	Ctx    context.Context
	Bucket *storage.BucketHandle
}

func New(bucketID string) Storage {
	ctx := context.Background()
	b := initBucket(ctx, bucketID)
	return Storage{
		Ctx:    ctx,
		Bucket: b,
	}
}

func (fs Storage) WriteFile(fileName string, content []byte) {
	wc := fs.Bucket.Object(fileName).NewWriter(fs.Ctx)

	if _, err := wc.Write(content); err != nil {
		log.Panic(errors.Wrapf(err, "unable to write data to bucket, file %q", fileName))
	}
	defer wc.Close()
}

func (fs Storage) GetFilenames() (names []string) {
	it := fs.Bucket.Objects(fs.Ctx, nil)
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

func (fs Storage) DeleteFile(fileName string) {
	err := fs.Bucket.Object(fileName).Delete(fs.Ctx)
	if err != nil {
		log.Panic(errors.Wrapf(err, "error while deleting file: %s", fileName))
	}
}

func initBucket(ctx context.Context, bucketID string) *storage.BucketHandle {
	app := initApp(ctx)
	s, err := app.Storage(ctx)
	if err != nil {
		log.Panic(errors.Wrap(err, "error while inizializing storage client"))
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
