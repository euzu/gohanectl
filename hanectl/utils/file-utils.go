package utils

import (
	"github.com/rs/zerolog/log"
	"io"
	"os"
	"path/filepath"
)

func FileExists(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}
	return true
}

func FileNotExists(filename string) bool {
	return !FileExists(filename)
}

func ChangeWorkingDir(workingdir string) {
	dir := workingdir
	if workingdir == "" {
		var err error
		dir, err = filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			log.Panic().Msgf("Cant change working directory to binary path")
		}
	}
	if err := os.Chdir(dir); err != nil {
		log.Panic().Msgf("Cant change working directory: %s", dir)
	}
	newDir, err := os.Getwd()
	if err != nil {
		log.Panic().Msgf("Cant change working directory: %s", newDir)
	}
	log.Info().Msgf("Current working directory: %s", newDir)
}

func CheckedClose(c io.Closer) {
	_ = c.Close()
}