package main

import (
	"github.com/Sandburg/mongodb-backup/src/backup"

	"github.com/Sandburg/mongodb-backup/src/config"
	"github.com/Sandburg/mongodb-backup/src/storage"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	c := config.NewConfig()
	s := storage.New(c.BucketID)

	b := backup.NewBackup(c, s)

	b.DumpDbToFile()
	b.UploadToBucket()
	b.DeleteOldBackups()
}
