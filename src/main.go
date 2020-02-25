package main

import (
	"fmt"
	"io/ioutil"
	"os/exec"

	"github.com/Sandburg/mongodb-backup/src/storage"
	log "github.com/sirupsen/logrus"
)

func main() {
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
