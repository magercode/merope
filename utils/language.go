package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"merope/models"
)

type LanguageManager struct {
	currentLang models.Language
	messages    models.Message
}

func NewLanguageManager(lang models.Language) (*LanguageManager, error) {
	lm := &LanguageManager{
		currentLang: lang,
	}

	err := lm.loadLanguage(lang)
	if err != nil {
		lm.loadLanguage(models.EN)
	}

	return lm, nil
}

func (lm *LanguageManager) loadLanguage(lang models.Language) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	langPath := filepath.Join(wd, "lang", string(lang)+".json")

	data, err := ioutil.ReadFile(langPath)
	if err != nil {
		return fmt.Errorf("failed to read language file: %v", err)
	}

	var messages models.Message
	err = json.Unmarshal(data, &messages)
	if err != nil {
		return fmt.Errorf("failed to parse language file: %v", err)
	}

	lm.messages = messages
	return nil
}

func (lm *LanguageManager) GetMessage(key string) string {
	switch strings.ToLower(key) {
	case "cpu_high":
		return lm.messages.CPUHigh
	case "memory_high":
		return lm.messages.MemoryHigh
	case "disk_high":
		return lm.messages.DiskHigh
	case "system_info":
		return lm.messages.SystemInfo
	case "alert_title":
		return lm.messages.AlertTitle
	case "alert_message":
		return lm.messages.AlertMessage
	case "time":
		return lm.messages.Time
	case "level":
		return lm.messages.Level
	case "ok":
		return lm.messages.OK
	case "warning":
		return lm.messages.Warning
	case "critical":
		return lm.messages.Critical
	case "service_started":
		return lm.messages.ServiceStarted
	default:
		return key
	}
}

func (lm *LanguageManager) FormatAlert(alert *models.Alert) string {
	return fmt.Sprintf("*%s*\n\n", alert.Title) +
		fmt.Sprintf("📊 *%s:* %s\n", lm.GetMessage("alert_message"), alert.Message) +
		fmt.Sprintf("⚠️ *%s:* %s\n", lm.GetMessage("level"), alert.Level) +
		fmt.Sprintf("⏰ *%s:* %s", lm.GetMessage("time"), alert.Time)
}