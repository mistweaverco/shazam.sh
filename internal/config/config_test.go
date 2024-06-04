package config

import (
	"reflect"
	"testing"
)

func TestGetConfig(t *testing.T) {
	tests := []struct {
		name string
		got  ConfigFile
		want ConfigFile
	}{
		{
			"should return a valid config",
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