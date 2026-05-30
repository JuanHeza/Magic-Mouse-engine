package util

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

var (
	//Version holds version data from 'version' file
	Version VersionData
)

// VersionData version data read from 'version' file
type VersionData struct {
	Version   string `yaml:"version"`
	Commit    string `yaml:"commit"`
	BuildDate string `yaml:"build_date"`
}

func processFile(filePath string, vd *VersionData) error {
	newVar := filepath.Clean(filePath)
	f, err := os.Open(newVar)
	if err != nil {
		return err
	}
	defer func() {
		if err := f.Close(); err != nil {
			println("Error closing file:", err.Error())
		}
	}()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(vd)
	if err != nil {
		return err
	}
	return nil
}

// InitVersion initialize
func InitVersion() error {
	err := processFile("version", &Version)
	if err != nil {
		return err
	}
	return nil
}
