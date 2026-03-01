package services

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"merope/models"
	"merope/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramService struct {
	enabled  bool
	bot      *tgbotapi.BotAPI
	chatID   int64
	lang     *utils.LanguageManager
}

func NewTelegramService(lang *utils.LanguageManager) *TelegramService {
	enabled := os.Getenv("TELEGRAM_ENABLED") == "true"
	if !enabled {
		return &TelegramService{enabled: false}
	}

	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	chatIDStr := os.Getenv("TELEGRAM_CHAT_ID")

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		fmt.Printf("Err in telegram services: %v\n", err)
		return &TelegramService{enabled: false}
	}

	chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
	if err != nil {
		fmt.Printf("Invalid Telegram chat ID: %v\n", err)
		return &TelegramService{enabled: false}
	}

	svc := &TelegramService{
		enabled: true,
		bot:     bot,
		chatID:  chatID,
		lang:    lang,
	}

	go svc.listen()

	return svc
}

func (t *TelegramService) Send(alert *models.Alert) error {
	if !t.enabled {
		return nil
	}

	var emoji string
	switch alert.Level {
	case models.INFO:
		emoji = "ℹ️"
	case models.WARNING:
		emoji = "⚠️"
	case models.CRITICAL:
		emoji = "🚨"
	}

	escape := func(s string) string {
		return strings.NewReplacer("_", "\\_", "*", "\\*", "[", "\\[", "`", "\\`").Replace(s)
	}

	message := fmt.Sprintf("%s *%s*\n\n", emoji, escape(alert.Title))
	message += fmt.Sprintf("📊 *%s:* %s\n", t.lang.GetMessage("alert_message"), escape(alert.Message))
	message += fmt.Sprintf("⚡ *%s:* %s\n", t.lang.GetMessage("level"), escape(string(alert.Level)))
	if alert.Recommendation != "" {
		message += fmt.Sprintf("🤖 *AI:* %s\n", escape(alert.Recommendation))
	}
	message += fmt.Sprintf("⏰ *%s:* %s", t.lang.GetMessage("time"), escape(alert.Time))

	msg := tgbotapi.NewMessage(t.chatID, message)
	msg.ParseMode = "Markdown"

	_, err := t.bot.Send(msg)
	if err != nil {
		return fmt.Errorf("failed to send Telegram message: %v", err)
	}

	return nil
}

func (t *TelegramService) IsEnabled() bool {
	return t.enabled
}

func (t *TelegramService) listen() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := t.bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.From.ID != t.chatID {
			continue
		}

		if strings.HasPrefix(update.Message.Text, "$ ") {
			command := strings.TrimPrefix(update.Message.Text, "$ ")
			output := t.runCommand(command)

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("```\n%s\n```", output))
			msg.ParseMode = "Markdown"
			t.bot.Send(msg)
		}
	}
}

func (t *TelegramService) runCommand(command string) string {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", command)
	} else {
		cmd = exec.Command("sh", "-c", command)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Sprintf("Error: %v\n%s", err, string(output))
	}

	if len(output) == 0 {
		return "No output"
	}

	if len(output) > 4000 {
		return string(output[:4000]) + "\n... (truncated)"
	}

	return string(output)
}