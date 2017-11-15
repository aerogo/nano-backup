package main

import (
	"fmt"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"runtime"
	"time"

	"github.com/fatih/color"
	"github.com/mholt/archiver"
)

const (
	interval        = 2 * time.Hour
	deleteThreshold = 48 * time.Hour
)

func main() {
	sourceDirectory, targetDirectory := setup()

	for {
		deleteOldFiles(targetDirectory)
		backup(sourceDirectory, targetDirectory)
		runtime.GC()
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

	color.Yellow("Creating backup %s", outFilePath)

	err := archiver.TarXZ.Make(outFilePath, []string{sourceDirectory})

	if err != nil {
		color.Red(err.Error())
	}

	color.Green("Finished.")
}

func deleteOldFiles(targetDirectory string) {
	filepath.Walk(targetDirectory, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		if time.Since(info.ModTime()) > deleteThreshold {
			color.Red("Deleting old backup %s", path)
			err := os.Remove(path)

			if err != nil {
				fmt.Println(err)
			}
		}

		return nil
	})
}
