package src

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"io"
	"os"
	"strings"
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

func getPathOfNearestConfigFile() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	pathArr := strings.Split(dir, "/")
	lenFromRoot := len(pathArr)

	for 0 <= lenFromRoot {
		path := strings.Join(pathArr[:lenFromRoot], "/")

		configFilePath := fmt.Sprintf("%s/hurl.json", path)

		configFile, _ := os.Open(configFilePath)

		if configFile != nil {
			return configFilePath, nil
		}

		lenFromRoot--
	}

	return "", nil
}

func InitConfig() (HurlConfig, error) {
	version := flag.Bool("version", false, "print version")
	verbose := flag.Bool("v", false, "verbose output")
	bodyOutputPath := flag.String("o", "", "path to a file to output the response body")

	flag.Parse()

	// go up until hit hurl.json or hit root dir
	configFileDir, err := getPathOfNearestConfigFile()
	if err != nil {
		return HurlConfig{}, err
	}

	hurlJson, err := os.Open(configFileDir)

	if !errors.Is(err, os.ErrNotExist) {
		if err != nil {
			return HurlConfig{}, err
		}

		configFileBytes, err := io.ReadAll(hurlJson)
		if err != nil {
			return HurlConfig{}, err
		}

		var hurlConfigFile hurlConfigFile
		err = json.Unmarshal(configFileBytes, &hurlConfigFile)
		if err != nil {
			return HurlConfig{}, err
		}

		// relative paths
		envFilePath := hurlConfigFile.EnvFilePath
		if !strings.HasPrefix(hurlConfigFile.EnvFilePath, "/") {
			if strings.HasPrefix(hurlConfigFile.EnvFilePath, "./") {
				envFilePath := strings.Replace(hurlConfigFile.EnvFilePath, "./", "", 1)
				hurlConfigFile.EnvFilePath = envFilePath
			}

			configFileDirPath := strings.Split(configFileDir, "/")
			configFileDirPath = configFileDirPath[:len(configFileDirPath)-1]
			configFileDir = strings.Join(configFileDirPath, "/")

			envFilePath = fmt.Sprintf("/%s/%s", configFileDir, hurlConfigFile.EnvFilePath)
		}

		err = godotenv.Load(envFilePath)
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
