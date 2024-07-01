package service

import (
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

func (s *NotifyService) NotifyAll() {
	usersWithBirthday := s.userService.CheckBirthdays()

	for _, user := range usersWithBirthday {
		messageText := "С Днем рождения, " + user.Name + "!"
		msg := tgbotapi.NewMessage(user.ChatID, messageText)
		_, err := s.bot.Send(msg)
		if err != nil {
			log.Printf("Error sending birthday message to %s: %v", user.Name, err)
		}
	}
}

func StartBirthdayNotifier(bot BirthdayNotifier, userService *UserService) {
	notifyService := NewNotifyService(bot, userService)

	for {
		now := time.Now()
		// Устанавливаем время следующего запуска на 10:00
		next := time.Date(now.Year(), now.Month(), now.Day(), 10, 0, 0, 0, now.Location())
		if now.After(next) {
			// Если текущее время уже после 10:00, устанавливаем следующий день
			next = next.Add(24 * time.Hour)
		}

		// Вычисляем длительность до следующего запуска
		duration := next.Sub(now)
		log.Printf("Next birthday check in %v", duration)

		// Ждем до следующего запуска
		time.Sleep(duration)
		// timer := time.NewTimer(duration)
		// <-timer.C

		// Выполняем уведомление о днях рождения
		notifyService.NotifyAll()
	}
}
