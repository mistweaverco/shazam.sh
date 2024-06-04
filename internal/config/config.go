package config

import (
	"os"

	"github.com/charmbracelet/log"
	"gopkg.in/yaml.v3"
)

var ps = string(os.PathSeparator)

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
	ConfigPath string
	DataReader ConfigDataReader
	File       ConfigFile
	Flags      ConfigFlags
}

func (c Config) GetConfigFile() ConfigFile {
	configPath := c.ConfigPath
	if c.Flags.DotfilesPath != "" {
		configPath = c.Flags.DotfilesPath + ps + c.ConfigPath
	}

	file, err := c.DataReader(configPath)
	if err != nil {
		log.Fatal(err)
	}

	err = yaml.Unmarshal(file, &c.File)
	if err != nil {
		log.Fatal("Error parsing "+configPath, "error", err)
	}
	return c.File
}

func (c Config) GetConfigFlags() ConfigFlags {
	return c.Flags
}

func NewConfig(cfg Config) Config {
	return cfg
}