package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/fs"
	"github.com/knadh/koanf/v2"
	"github.com/naivary/apitmpl/configs"
	"golang.org/x/exp/slog"
)

const (
	keyPathDelimiter = "."
	DefaultCfgFile   = ""
)

func defaultConfig() (*koanf.Koanf, error) {
	k := koanf.New(keyPathDelimiter)
	err := k.Load(fs.Provider(configs.Fs, "default.yaml"), yaml.Parser())
	if err != nil {
		return nil, err
	}
	if err = buildDirs(k); err != nil {
		return nil, err
	}
	return k, nil
}

func customConfig(path string) (*koanf.Koanf, error) {
	k := koanf.New(keyPathDelimiter)
	err := k.Load(file.Provider(path), yaml.Parser())
	if err != nil {
		return nil, err
	}
	if err = buildDirs(k); err != nil {
		return nil, err
	}
	return k, err

}

func New(path string) (*koanf.Koanf, error) {
	if path == "" {
		return defaultConfig()
	}
	return customConfig(path)
}

// getBasePath checks if the dataDir is
// set. Otherwise the current directory
// will be used.
func getBasepath(k *koanf.Koanf) string {
	if k.Exists("dataDir") {
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
	if err := os.MkdirAll(versionDir, os.ModePerm); err != nil {
		return err
	}
	if err := k.Set("versionDir", versionDir); err != nil {
		return err
	}

	// create the backup dir
	backupDir := filepath.Join(versionDir, "backup")
	if err := os.MkdirAll(backupDir, os.ModePerm); err != nil {
		return err
	}
	if err := k.Set("backupDir", backupDir); err != nil {
		return err
	}

	// create logging dir
	logsDir := filepath.Join(versionDir, "logs")
	if err := os.MkdirAll(logsDir, os.ModePerm); err != nil {
		return err
	}
	return k.Set("logsDir", logsDir)
}
