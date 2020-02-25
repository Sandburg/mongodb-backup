package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os/exec"

	"github.com/Sandburg/mongodb-backup/src/storage"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

func main() {
	godotenv.Load()
	c := initConfig()

	storage := storage.New(c.BucketID)

	cmd := exec.Command("mongodump", buildArgs(c)...)

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

func buildArgs(c Config) []string {
	args := []string{
		"--host",
		c.Host,
		"--port",
		c.Port,
		"--gzip",
		fmt.Sprintf("--archive=%s/%s", c.ArchiveDir, c.ArchiveName),
	}

	if c.DB != "" {
		args = append(args, c.DB)
	}

	if c.AuthPart != "" {
		args = append(args, c.AuthPart)
	}

	return args
}
