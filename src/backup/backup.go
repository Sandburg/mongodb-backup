package backup

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os/exec"

	"github.com/Sandburg/mongodb-backup/src/config"
	"github.com/labstack/gommon/log"
)

type Storage interface {
	WriteFile(fn string, content []byte)
	GetFilenames() (names []string)
	DeleteFile(fn string)
}

type Backup struct {
	c config.Config
	s Storage
}

func NewBackup(c config.Config, s Storage) *Backup {
	return &Backup{c, s}
}

func (b Backup) DumpDbToFile() {
	cmd := exec.Command("mongodump", b.buildArgs()...)

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

}

func (b Backup) UploadToBucket() {
	dumpData, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", b.c.ArchiveDir, b.c.ArchiveName))
	if err != nil {
		log.Panic(err)
	}

	log.Infof("Start uploading file")

	b.s.WriteFile(b.c.ArchiveName, dumpData)

	log.Infof("Finished uploading file")
}

func (b Backup) DeleteOldBackups() {
	fn := b.s.GetFilenames()
	if len(fn) > b.c.MaxBackupItems {
		log.Infof("Deleting oldest backup: %s", fn[0])
		b.s.DeleteFile(fn[0])
	}
}

func (b Backup) buildArgs() []string {
	args := []string{"--host", b.c.Host, "--port", b.c.Port, "--gzip", fmt.Sprintf("--archive=%s/%s", b.c.ArchiveDir, b.c.ArchiveName)}

	if b.c.DB != "" {
		args = append(args, b.c.DB)
	}

	if b.c.AuthPart != "" {
		args = append(args, b.c.AuthPart)
	}

	return args
}
