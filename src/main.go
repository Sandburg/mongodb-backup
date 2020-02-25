package main

import (
	"bytes"
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

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		log.Panicf("failed dumping mongodb - %s: %s", fmt.Sprint(err), stderr.String())
	}

	fmt.Printf(out.String())

	log.Infof("Finished mongodump")

	dumpData, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", c.ArchiveDir, c.ArchiveName))
	if err != nil {
		log.Panic(err)
	}

	storage.WriteFile(c.ArchiveName, dumpData)
}
