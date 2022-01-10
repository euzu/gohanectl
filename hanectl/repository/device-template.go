package repository

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"gohanectl/hanectl/config"
	"gohanectl/hanectl/utils"
	"os"
	"path"
	"text/template"
)

const YmlSuffix = ".yml"
const JsonSuffix = ".json"

func getDevicesTemplatePath(cfg config.IConfiguration) string {
	return path.Join(
		cfg.GetStr(config.ConfigDirectory, config.DefConfigDirectory),
		cfg.GetStr(config.ScriptsDirectory, config.DefScriptsDirectory),
		cfg.GetStr(config.ScriptsTemplatesDirectory, config.DefScriptsTemplatesDirectory))
}

func getTemplateNestedPath(cfg config.IConfiguration, teplateName string, configKey config.ConfigKey, defValue string) (string, error) {
	templateDir := cfg.GetStr(configKey, defValue)
	if _, err := os.Stat(templateDir); err != nil {
		templateDir = path.Join(getDevicesTemplatePath(cfg), templateDir)
		if _, err := os.Stat(templateDir); err != nil {
			log.Fatal().Msgf("Cant find device template directory %s", templateDir)
			return "", errors.New("dir not found")
		}
	}

	templatePath := path.Join(templateDir, teplateName)
	templateYml := fmt.Sprintf("%s%s", templatePath, YmlSuffix)
	if _, err := os.Stat(templateDir); err == nil {
		return templateYml, nil
	} else {
		templateJson := fmt.Sprintf("%s%s", templatePath, JsonSuffix)
		if _, err := os.Stat(templateDir); err == nil {
			return templateJson, nil

		}
	}
	return "", errors.New("cant find template file")
}

func getDevicesConfigPath(cfg config.IConfiguration, templateName string) (string, error) {
	return getTemplateNestedPath(cfg, templateName, config.ScriptsTemplatesDevicesDirectory,
		config.DefScriptsTemplatesDevicesDirectory)
}

func getDeviceTemplateContent(templateFile string, model interface{}) ([]byte, error) {
	if utils.IsBlank(templateFile) {
		return nil, errors.New("template file not configured")
	}

	if _, err := os.Stat(templateFile); err == nil {
		t, err := template.ParseFiles(templateFile)
		if err != nil {
			log.Error().Msgf("Failed to parse template: %s", templateFile)
		} else {
			var tpl bytes.Buffer
			if err := t.Execute(&tpl, model); err != nil {
				log.Error().Msgf("Failed to execute template: %s, err: %v", templateFile, err)
			} else {
				return tpl.Bytes(), nil
			}
		}
	}
	return nil, errors.New("not found or failed")
}
