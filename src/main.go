package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"time"

	"github.com/Sandburg/mongodb-backup/src/storage"
	"github.com/joho/godotenv"
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

func main() {
	godotenv.Load()
	c := initConfig()

	storage := storage.New(c.BucketID)

	cmd := exec.Command(
		"mongodump",
		fmt.Sprintf("--host=%s", c.Host),
		fmt.Sprintf("--port=%s", c.Port),
		fmt.Sprintf("%s", c.AuthPart),
		fmt.Sprintf("%s", c.DB),
		fmt.Sprintf("--archive=%s/%s", c.ArchiveDir, c.ArchiveName),
	)

	log.Infof("Starting mongodump: %s", cmd.String())

	out, err := cmd.Output()
	if err != nil {
		log.Errorf("failed dumping mongodb: %s", string(out))
	}

	log.Info("Finished mongodump")

	dumpData, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", c.ArchiveDir, c.ArchiveName))
	if err != nil {
		log.Panic(err)
	}

	storage.WriteFile(c.ArchiveName, dumpData)
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
		auth = fmt.Sprintf("--username=%s --password=%s", os.Getenv("MONGODB_USER"), os.Getenv("MONGODB_PASSWORD"))
	}

	if os.Getenv("MONGODB_DB") != "" {
		db = fmt.Sprintf("--db=%s", os.Getenv("MONGODB_DB"))
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
