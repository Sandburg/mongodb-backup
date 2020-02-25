package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
)

type Config struct {
	BucketID       string
	Host           string
	Port           string
	ArchiveDir     string
	ArchiveName    string
	AuthPart       string
	DB             string
	MaxBackupItems int64
}

func initConfig() Config {
	bucketID := os.Getenv("GCS_BUCKET")
	if bucketID == "" {
		log.Panic("bucket id can not be empty, set it as environment variable GCS_BUCKET")
	}

	var (
		host           = "localhost"
		port           = "27017"
		dir            = "/tmp"
		archiveName    = fmt.Sprintf("backup-%d.gz", time.Now().Unix())
		auth           string
		db             string
		maxBackupItems int64 = 10
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
		auth = fmt.Sprintf("--username=%s --password=%s", os.Getenv("MONGODB_USER"), os.Getenv("MONGODB_PASSWORD"))
	}

	if os.Getenv("MONGODB_DB") != "" {
		db = fmt.Sprintf("--db=%s", os.Getenv("MONGODB_DB"))
	}

	if os.Getenv("MAX_BACKUP_ITEMS") != "" {
		maxBackupItems, _ = strconv.ParseInt(os.Getenv("MONGODB_DB"), 10, 64)
	}

	return Config{
		BucketID:       bucketID,
		Host:           host,
		Port:           port,
		ArchiveDir:     dir,
		ArchiveName:    archiveName,
		AuthPart:       auth,
		DB:             db,
		MaxBackupItems: maxBackupItems,
	}
}
