package src

import (
	"bufio"
	"encoding/json"
	"flag"
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
	Verbose        bool
	BodyOutputPath string
}

func readEnvironmentVariables(path string) error {
	envFile, err := os.Open(path)
	if err != nil {
		return err
	}

	sc := bufio.NewScanner(envFile)
	for sc.Scan() {
		environmentVariable := strings.SplitN(sc.Text(), "=", 2)

		key := environmentVariable[ENVKEY]
		value := environmentVariable[ENVVALUE]

		os.Setenv(key, value)
	}

	return nil
}

func InitConfig() (HurlConfig, error) {
	verbose := flag.Bool("v", false, "verbose output")
	bodyOutputPath := flag.String("o", "", "path to a file to output the response body")

	flag.Parse()

	configFile, err := os.Open("hurl.json")
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

	err = readEnvironmentVariables(hurlConfigFile.EnvFilePath)
	if err != nil {
		return HurlConfig{}, err
	}

	return HurlConfig{
		*verbose,
		*bodyOutputPath,
	}, nil
}
