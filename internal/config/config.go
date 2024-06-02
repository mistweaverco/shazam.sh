package config

import (
	"os"

	"github.com/charmbracelet/log"
	"gopkg.in/yaml.v3"
)

var Config ConfigFile
var ConfigPath string
var Flags ConfigFlags

type ConfigFlags struct {
	DotfilesPath   string
	DryRun         bool
	Only           string
	Path           string
	PullInExisting bool
	Root           string
}

type ConfigFiles struct {
	Source      string `yaml:"source"`
	Destination string `yaml:"destination"`
}

type ConfigNode struct {
	Name  string        `yaml:"name"`
	Files []ConfigFiles `yaml:"files"`
}

type ConfigFile map[string][]ConfigNode

func GetConfig() ConfigFile {
	var cfg ConfigFile
	var ps = string(os.PathSeparator)

	configPath := ConfigPath
	if Flags.DotfilesPath != "" {
		configPath = Flags.DotfilesPath + ps + ConfigPath
	}

	file, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatal(err)
	}

	err = yaml.Unmarshal(file, &cfg)
	if err != nil {
		log.Fatal("Error parsing "+ConfigPath, "error", err)
	}

	return cfg
}