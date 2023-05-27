package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"golang.org/x/exp/slog"
)

const (
	keyPathDelimiter = "."
)

func New(path string) (*koanf.Koanf, error) {
	k := koanf.New(keyPathDelimiter)
	err := k.Load(file.Provider(path), yaml.Parser())
	if err != nil {
		return nil, err
	}
	err = buildDirs(k)
	if err != nil {
		return nil, err
	}
	return k, err
}

// getBasePath checks if the dataDir is
// set. Otherwise the current directory
// will be used.
func getBasepath(k *koanf.Koanf) string {
	ok := k.Exists("dataDir")
	if ok {
		return k.String("dataDir")
	}
	basepath, err := os.Getwd()
	if err != nil {
		slog.Error("could not get current working directory", "err", err)
		return ""
	}
	return basepath
}

// buildDirs creates the needed directories for
// the current API version to serve.
func buildDirs(k *koanf.Koanf) error {
	basepath := getBasepath(k)
	dataDirName := fmt.Sprintf("%s_data", k.String("name"))
	versionDir := filepath.Join(basepath, dataDirName, k.String("version"))

	// create the version dir (e.g. 0.1.0) which stores
	// the main database of the application.
	err := os.MkdirAll(versionDir, os.ModePerm)
	if err != nil {
		return err
	}
	err = k.Set("versionDir", versionDir)
	if err != nil {
		return err
	}

	// create the backup dir
	backupDir := filepath.Join(versionDir, "backup")
	err = os.MkdirAll(backupDir, os.ModePerm)
	if err != nil {
		return err
	}
	err = k.Set("backupDir", backupDir)
	if err != nil {
		return err
	}

	// create logging dir
	logsDir := filepath.Join(versionDir, "logs")
	err = os.MkdirAll(logsDir, os.ModePerm)
	if err != nil {
		return err
	}
	return k.Set("logsDir", logsDir)
}
