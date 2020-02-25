package main

import (
	"fmt"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

type Config struct {
	BucketID    string
	Host        string
	Port        string
	ArchiveDir  string
	ArchiveName string
	AuthPart    string
	DB          string
}

func initConfig() Config {
	bucketID := os.Getenv("GCS_BUCKET")
	if bucketID == "" {
		log.Panic("bucket id can not be empty, set it as environment variable GCS_BUCKET")
	}

	var (
		host        = "localhost"
		port        = "27017"
		dir         = "/tmp"
		archiveName = fmt.Sprintf("backup-%s.tar.gz", time.Now().Format(time.RFC3339))
		auth        string
		db          string
	)

	if os.Getenv("MONGODB_HOST") != "" {
		host = os.Getenv("MONGODB_HOST")
	}

	if os.Getenv("MONGODB_PORT") != "" {
		port = os.Getenv("MONGODB_PORT")
	}

	if os.Getenv("BACKUP_DIR") != "" {
		dir = os.Getenv("BACKUP_DIR")
	}

	if os.Getenv("MONGODB_USER") != "" && os.Getenv("MONGODB_PASSWORD") != "" {
		auth = fmt.Sprintf(`--username="%s" --password="%s"`, os.Getenv("MONGODB_USER"), os.Getenv("MONGODB_PASSWORD"))
	}

	if os.Getenv("MONGODB_DB") != "" {
		db = fmt.Sprintf(`--db="%s"`, os.Getenv("MONGODB_DB"))
	}

	return Config{
		BucketID:    bucketID,
		Host:        host,
		Port:        port,
		ArchiveDir:  dir,
		ArchiveName: archiveName,
		AuthPart:    auth,
		DB:          db,
	}
}
