package repository

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"gohanectl/hanectl/config"
	"gohanectl/hanectl/utils"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

func readConfiguration(cfg config.IConfiguration, configKey config.ConfigKey, defaultFile string, model interface{}) (interface{}, error) {
	fileName := cfg.GetStr(configKey, defaultFile)
	if fileName == "" {
		return nil, errors.New("config file not configured")
	}
	newFileName := path.Join(
		cfg.GetStr(config.ConfigDirectory, config.DefConfigDirectory),
		fileName)
	if utils.FileExists(newFileName) {
		fileName = newFileName
	} else if utils.FileNotExists(fileName) {
		return nil, errors.New(fmt.Sprintf("Could not find configuration file %s", fileName))
	}

	if cwd, err := os.Getwd(); err == nil {
		log.Info().Msgf("Reading configuration file: %s", path.Join(cwd, fileName))
	}

	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	return readConfigurationFromContent(fileName, content, model)
}

func readConfigurationFromContent(fileName string, content []byte, model interface{}) (interface{}, error) {

	if strings.HasSuffix(fileName, ".yml") {
		if err := yaml.Unmarshal(content, model); err != nil {
			log.Fatal().Msgf("Failed to read file:%s  %v", fileName, err)
			return nil, err
		}
		return model, nil
	} else if strings.HasSuffix(fileName, ".json") {
		decoder := json.NewDecoder(bytes.NewReader(content))
		if err := decoder.Decode(model); err != nil {
			log.Fatal().Msgf("Failed to read file: %s, %v", fileName, err)
			return nil, err
		}
		return model, nil
	}

	return nil, errors.New("failed to read file: unknown filetype")
}
