package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"merope/models"
	"merope/services"
	"merope/utils"

	"github.com/joho/godotenv"
)

type NotificationService interface {
	Send(alert *models.Alert) error
	IsEnabled() bool
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	langCode := os.Getenv("NOTIFICATION_LANGUAGE")
	if langCode == "" {
		langCode = "en"
	}
	
	langManager, err := utils.NewLanguageManager(models.Language(langCode))
	if err != nil {
		log.Printf("Failed to load language: %v, using default", err)
	}

	geminiService := services.NewGeminiService()
	if geminiService.IsEnabled() {
		defer geminiService.Close()
		log.Println("Gemini AI integration enabled")
	}

	services := []NotificationService{
		services.NewEmailService(langManager),
		services.NewTelegramService(langManager),
	}

	hasEnabledService := false
	for _, s := range services {
		if s.IsEnabled() {
			hasEnabledService = true
			break
		}
	}

	if !hasEnabledService {
		log.Println("Warning: No notification service is enabled. Please check your .env configuration.")
	}

	sendStartupNotification(services, langManager)

	checkInterval := getCheckInterval()

	log.Printf("Merope monitoring started. Checking every %d seconds", checkInterval)
	
	for {
		alert, err := utils.CheckSystem()
		if err != nil {
			log.Printf("Error checking system: %v", err)
		}

		if alert != nil {
			log.Printf("Alert triggered: %s - %s", alert.Title, alert.Message)
			
			if geminiService.IsEnabled() {
				analysis, err := geminiService.AnalyzeAlert(alert)
				if err != nil {
					log.Printf("Gemini analysis failed: %v", err)
				} else {
					alert.Recommendation = analysis
				}
			}

			sendAlertToServices(services, alert)
		}

		time.Sleep(time.Duration(checkInterval) * time.Second)
	}
}

func getCheckInterval() int {
	interval := 60 
	intervalStr := os.Getenv("CHECK_INTERVAL")
	if intervalStr != "" {
		fmt.Sscanf(intervalStr, "%d", &interval)
	}
	if interval < 5 {
		interval = 5 
	}
	return interval
}

func sendStartupNotification(services []NotificationService, lang *utils.LanguageManager) {
	alert := &models.Alert{
		Title:   lang.GetMessage("service_started"),
		Message: fmt.Sprintf("System: %s\nCPU Cores: %d\nMemory: %.2f GB\nDisk: %.2f GB", 
			utils.GetOsType(),
			utils.GetTotalCore(),
			float64(utils.GetTotalMemory())/1024/1024/1024,
			float64(utils.GetTotalDisk())/1024/1024/1024),
		Level: models.INFO,
		Time:  time.Now().Format(time.RFC3339),
	}

	sendAlertToServices(services, alert)
}

func sendAlertToServices(services []NotificationService, alert *models.Alert) {
	for _, service := range services {
		if service.IsEnabled() {
			go func(s NotificationService, a *models.Alert) {
				err := s.Send(a)
				if err != nil {
					log.Printf("Failed to send notification: %v", err)
				}
			}(service, alert)
		}
	}
}