package vm

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"gohanectl/hanectl/config"
	"gohanectl/hanectl/model"
	"os"
	"path"
	"strings"
	"text/template"
)

func getParsedScriptTemplate(dev *model.Device, templateFile string) (string, error) {
	if _, err := os.Stat(templateFile); err == nil {
		t, err := template.ParseFiles(templateFile)
		if err != nil {
			log.Error().Msgf("Failed to parse template: %s", templateFile)
		} else {
			var tpl bytes.Buffer
			if err := t.Execute(&tpl, dev); err != nil {
				log.Error().Msgf("Failed to execute template: %s, err: %v", templateFile, err)
			} else {
				var content = tpl.String()
				if strings.Index(content, "(function () {")  < 0 {
					content = fmt.Sprintf("(function () {%s})();", content)
				}
				return content, nil
			}
		}
	}
	return "", errors.New("not found")
}

func getScriptPath(cfg config.IConfiguration) string {
	return path.Join(
		cfg.GetStr(config.ConfigDirectory, config.DefConfigDirectory),
		cfg.GetStr(config.ScriptsDirectory, config.DefScriptsDirectory))
}

func getScriptTemplatesPath(cfg config.IConfiguration) string {
	return path.Join(
		getScriptPath(cfg),
		cfg.GetStr(config.ScriptsTemplatesDirectory, config.DefScriptsTemplatesDirectory))
}

func getNestedScriptPath(scriptName string, scriptDir string) (string, error) {
	if _, err := os.Stat(scriptDir); err != nil {
		log.Fatal().Msgf("Cant find directory %s", scriptDir)
		return "", errors.New("dir not found")
	}

	scriptPath := path.Join(scriptDir, scriptName)
	if !strings.HasSuffix(scriptPath, JsSuffix) {
		scriptPath = fmt.Sprintf("%s%s", scriptPath, JsSuffix)
	}
	return scriptPath, nil
}

func getScriptTemplatesNestedPath(cfg config.IConfiguration, scriptName string, configKey config.ConfigKey, defValue string) (string, error) {
	scriptDir := path.Join(getScriptTemplatesPath(cfg), cfg.GetStr(configKey, defValue))
	return getNestedScriptPath(scriptName, scriptDir)
}

func getScriptNestedPath(cfg config.IConfiguration, scriptName string, configKey config.ConfigKey, defValue string) (string, error) {
	scriptDir := path.Join(getScriptPath(cfg), cfg.GetStr(configKey, defValue))
	return getNestedScriptPath(scriptName, scriptDir)
}
