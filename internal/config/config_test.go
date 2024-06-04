package config

import (
	"reflect"
	"testing"
	"testing/fstest"
)

func TestGetConfig(t *testing.T) {
	tests := []struct {
		name       string
		configPath string
		got        fstest.MapFS
		want       ConfigFile
	}{
		{
			"should return empty config",
			"shazam.yml",
			fstest.MapFS{
				"shazam.yml": {
					Data: []byte("hello, world"),
				},
			},
			ConfigFile{
				"bash": {
					{
						Name: "bash",
						Files: []ConfigFiles{
							{
								Source:      "hello.txt",
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
			ConfigPath = tt.configPath
			got := GetConfig()
			eq := reflect.DeepEqual(got, tt.want)
			if !eq {
				t.Errorf("GetConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}