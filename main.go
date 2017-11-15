package main

import (
	"os"
	"os/user"
	"path"
	"time"

	"github.com/fatih/color"
	"github.com/mholt/archiver"
)

const interval = 1 * time.Hour

func main() {
	sourceDirectory, targetDirectory := setup()

	for {
		backup(sourceDirectory, targetDirectory)
		time.Sleep(interval)
	}
}

func setup() (string, string) {
	user, err := user.Current()

	if err != nil {
		panic(err)
	}

	sourceDirectory := path.Join(user.HomeDir, ".aero/db/")
	targetDirectory := path.Join(user.HomeDir, ".aero/backups/")

	// Create directory in case it doesn't exist
	os.MkdirAll(targetDirectory, 0777)

	return sourceDirectory, targetDirectory
}

func backup(sourceDirectory, targetDirectory string) {
	timestamp := time.Now().UTC().Format(time.RFC3339)
	outFileName := "db-" + timestamp + ".tar.xz"
	outFilePath := path.Join(targetDirectory, outFileName)

	color.Yellow(outFilePath)

	err := archiver.TarXZ.Make(outFilePath, []string{sourceDirectory})

	if err != nil {
		color.Red(err.Error())
	}

	color.Green(outFilePath)
}
