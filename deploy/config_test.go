package deploy_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/a-novel/golib/deploy"
)

func TestLoadConfig(t *testing.T) {
	configFiles := map[string][]byte{
		"item1":     []byte("item1: foo"),
		"item1Alt":  []byte("item1: qux"),
		"item2":     []byte("item2: bar"),
		"bothItems": []byte("item1: one\nitem2: two"),
	}

	type config struct {
		Item1 string `yaml:"item1"`
		Item2 string `yaml:"item2"`
	}

	testCases := []struct {
		name string

		env   string
		files []deploy.ConfigFile

		expect config
	}{
		{
			name: "DevEnv",
			env:  deploy.DevENV,
			files: []deploy.ConfigFile{
				deploy.GlobalConfig(configFiles["bothItems"]),
				// Override field in first entry.
				deploy.DevConfig(configFiles["item1"]),
				// Ignored.
				deploy.StagingConfig(configFiles["item1Alt"]),
			},
			expect: config{
				Item1: "foo",
				Item2: "two",
			},
		},
		{
			name: "StagingEnv",
			env:  deploy.StagingEnv,
			files: []deploy.ConfigFile{
				deploy.GlobalConfig(configFiles["bothItems"]),
				// Override field in first entry.
				deploy.StagingConfig(configFiles["item1Alt"]),
				// Ignored.
				deploy.DevConfig(configFiles["item1"]),
			},
			expect: config{
				Item1: "qux",
				Item2: "two",
			},
		},
		{
			name: "ProdEnv",
			env:  deploy.ProdENV,
			files: []deploy.ConfigFile{
				deploy.GlobalConfig(configFiles["bothItems"]),
				// Override field in first entry.
				deploy.ProdConfig(configFiles["item1Alt"]),
				// Ignored.
				deploy.DevConfig(configFiles["item1"]),
			},
			expect: config{
				Item1: "qux",
				Item2: "two",
			},
		},

		{
			name: "NoDefaultValue",
			env:  deploy.DevENV,
			files: []deploy.ConfigFile{
				deploy.DevConfig(configFiles["item1"]),
			},
			expect: config{Item1: "foo"},
		},
		{
			name: "FileOrderMatters",
			env:  deploy.DevENV,
			files: []deploy.ConfigFile{
				deploy.DevConfig(configFiles["item1"]),
				deploy.GlobalConfig(configFiles["bothItems"]),
			},
			expect: config{
				Item1: "one",
				Item2: "two",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			deploy.ENV = tc.env

			cfg := deploy.LoadConfig[config](tc.files...)

			require.Equal(t, tc.expect, *cfg)
		})
	}
}
