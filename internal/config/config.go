package config

import (
	"os"

	"github.com/charmbracelet/log"
	"gopkg.in/yaml.v3"
)

type ConfigFiles struct {
	Source      string `yaml:"source"`
	Destination string `yaml:"destination"`
}

type ConfigNode struct {
	Name  string        `yaml:"name"`
	Files []ConfigFiles `yaml:"files"`
}

type ConfigFile map[string][]ConfigNode
type ConfigDataReader func(path string) ([]byte, error)

type ConfigFlags struct {
	DotfilesPath   string
	DryRun         bool
	Only           string
	Path           string
	PullInExisting bool
	Root           string
}

type Config struct {
	File       ConfigFile
	Flags      ConfigFlags
	DataReader ConfigDataReader
}

var ConfigPath string
var Flags ConfigFlags

func (c Config) GetConfigFile() ConfigFile {
	return c.File
}

func (c Config) GetConfigFlags() ConfigFlags {
	return c.Flags
}

func NewConfig(cfg Config) Config {
	var ps = string(os.PathSeparator)

	configPath := ConfigPath
	if Flags.DotfilesPath != "" {
		configPath = Flags.DotfilesPath + ps + ConfigPath
	}

	file, err := cfg.DataReader(configPath)
	if err != nil {
		log.Fatal(err)
	}

	err = yaml.Unmarshal(file, &cfg.File)
	if err != nil {
		log.Fatal("Error parsing "+ConfigPath, "error", err)
	}

	return cfg
}