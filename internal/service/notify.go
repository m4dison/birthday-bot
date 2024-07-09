package service

import (
	"context"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type BirthdayNotifier interface {
	Send(c tgbotapi.Chattable) (tgbotapi.Message, error)
}

type TelegramBotAdapter struct {
	bot *tgbotapi.BotAPI
}

func NewTelegramBotAdapter(bot *tgbotapi.BotAPI) *TelegramBotAdapter {
	return &TelegramBotAdapter{bot: bot}
}

func (t *TelegramBotAdapter) Send(c tgbotapi.Chattable) (tgbotapi.Message, error) {
	return t.bot.Send(c)
}

type NotifyService struct {
	bot         BirthdayNotifier
	userService *UserService
}

func NewNotifyService(bot BirthdayNotifier, userService *UserService) *NotifyService {
	return &NotifyService{
		bot:         bot,
		userService: userService,
	}
}

func (s *NotifyService) NotifyAll(ctx context.Context) {
	usersWithBirthday, err := s.userService.CheckBirthdays(ctx)
	if err != nil {
		log.Printf("Error checking birthdays: %v", err)
		return
	}

	for _, user := range usersWithBirthday {
		messageText := "С Днем рождения, " + user.Name + "!"
		msg := tgbotapi.NewMessage(user.ChatID, messageText)
		_, err := s.bot.Send(msg)
		if err != nil {
			log.Printf("Error sending birthday message to %s: %v", user.Name, err)
		}
	}
}

func StartBirthdayNotifier(ctx context.Context, bot BirthdayNotifier, userService *UserService) {
	notifyService := NewNotifyService(bot, userService)

	for {
		now := time.Now()

		next := time.Date(now.Year(), now.Month(), now.Day(), 10, 0, 0, 0, now.Location())
		if now.After(next) {
			next = next.Add(24 * time.Hour)
		}

		duration := next.Sub(now)
		log.Printf("Next birthday check in %v", duration)

		select {
		case <-time.After(duration):
			// Выполняем уведомления о днях рождения
			notifyService.NotifyAll(ctx)
		case <-ctx.Done():
			log.Println("BirthdayNotifier: received shutdown signal")
			return
		}
	}
}
