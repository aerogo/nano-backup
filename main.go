package main

import (
	"os"
	"os/user"
	"path"
	"time"

	"github.com/fatih/color"
	"github.com/mholt/archiver"
)

func main() {
	const interval = 1 * time.Hour

	user, err := user.Current()

	if err != nil {
		panic(err)
	}

	sourceDirectory := path.Join(user.HomeDir, ".aero/db/")
	targetDirectory := path.Join(user.HomeDir, ".aero/backups/")

	// Create directory in case it doesn't exist
	os.MkdirAll(targetDirectory, 0777)

	for {
		timestamp := time.Now().UTC().Format(time.RFC3339)
		outFileName := "db-" + timestamp + ".tar.xz"
		outFilePath := path.Join(targetDirectory, outFileName)

		color.Yellow(outFilePath)

		err := archiver.TarXZ.Make(outFilePath, []string{sourceDirectory})

		if err != nil {
			color.Red(err.Error())
		}

		color.Green(outFilePath)

		time.Sleep(interval)
	}
}
