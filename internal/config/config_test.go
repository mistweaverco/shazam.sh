package config

import (
	"reflect"
	"testing"
)

func TestConfigFile(t *testing.T) {
	tests := []struct {
		name string
		got  ConfigFile
		want ConfigFile
	}{
		{
			"Test GetConfig",
			NewConfig(Config{
				DataReader: func(path string) ([]byte, error) {
					var yamlContents = `
configuration:
  - name: bash
    files:
      - source: .bashrc
        destination: .bashrc`
					return []byte(yamlContents), nil
				},
			}).GetConfigFile(),
			ConfigFile{
				"configuration": {
					{
						Name: "bash",
						Files: []ConfigFiles{
							{
								Source:      ".bashrc",
								Destination: ".bashrc",
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			eq := reflect.DeepEqual(tt.got, tt.want)
			if !eq {
				t.Errorf("GetConfig() = %v, want %v", tt.got, tt.want)
			}
		})
	}
}

func TestConfigFlags(t *testing.T) {
	tests := []struct {
		name string
		got  ConfigFlags
		want ConfigFlags
	}{
		{
			"Test GetConfig",
			NewConfig(Config{
				Flags: ConfigFlags{
					DotfilesPath:   "",
					DryRun:         false,
					Only:           "",
					Path:           "",
					PullInExisting: false,
					Root:           "",
				},
			}).GetConfigFlags(),
			ConfigFlags{
				DotfilesPath:   "",
				DryRun:         false,
				Only:           "",
				Path:           "",
				PullInExisting: false,
				Root:           "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			eq := reflect.DeepEqual(tt.got, tt.want)
			if !eq {
				t.Errorf("GetConfig() = %v, want %v", tt.got, tt.want)
			}
		})
	}
}