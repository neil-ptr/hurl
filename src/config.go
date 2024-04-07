package src

import (
	"encoding/json"
	"errors"
	"flag"
	"io"
	"os"

	"github.com/joho/godotenv"
)

const (
	ENVKEY   = 0
	ENVVALUE = 1
)

type hurlConfigFile struct {
	EnvFilePath string `json:"env"`
}

type HurlConfig struct {
	Version        bool
	Verbose        bool
	BodyOutputPath string
}

func InitConfig() (HurlConfig, error) {
	version := flag.Bool("version", false, "print version")
	verbose := flag.Bool("v", false, "verbose output")
	bodyOutputPath := flag.String("o", "", "path to a file to output the response body")

	flag.Parse()

	configFile, err := os.Open("hurl.json")

	if !errors.Is(err, os.ErrNotExist) {
		if err != nil {
			return HurlConfig{}, err
		}

		configFileBytes, err := io.ReadAll(configFile)
		if err != nil {
			return HurlConfig{}, err
		}

		var hurlConfigFile hurlConfigFile
		err = json.Unmarshal(configFileBytes, &hurlConfigFile)
		if err != nil {
			return HurlConfig{}, err
		}

		err = godotenv.Load(hurlConfigFile.EnvFilePath)
		if err != nil {
			return HurlConfig{}, err
		}
	}

	return HurlConfig{
		*version,
		*verbose,
		*bodyOutputPath,
	}, nil
}
