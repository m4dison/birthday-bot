package controller

import (
	"context"
	"log"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/m4dison/my-telegram-bot/internal/models"
	"github.com/m4dison/my-telegram-bot/internal/service"
)

type BotController struct {
	bot            *tgbotapi.BotAPI
	userService    *service.UserService
	pendingAddUser map[int64]struct{} // хранит идентификаторы чатов с ожидаемыми данными пользователей
}

func NewBotController(bot *tgbotapi.BotAPI, userService *service.UserService) *BotController {
	return &BotController{
		bot:            bot,
		userService:    userService,
		pendingAddUser: make(map[int64]struct{}),
	}
}

func (bc *BotController) Start(ctx context.Context) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bc.bot.GetUpdatesChan(u)

	for {
		select {
		case update := <-updates:
			if update.Message == nil { // игнорируем любые не Message обновления
				continue
			}

			go func(update tgbotapi.Update) {
				ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
				defer cancel()

				if update.Message.IsCommand() {
					bc.handleCommand(update.Message)
				} else {
					bc.handleNonCommand(ctx, update.Message)
				}
			}(update)
		case <-ctx.Done():
			log.Println("BotController: received shutdown signal")
			return
		}
	}
}

func (bc *BotController) handleCommand(message *tgbotapi.Message) {
	switch message.Command() {
	case "adduser":
		bc.pendingAddUser[message.Chat.ID] = struct{}{}
		msg := tgbotapi.NewMessage(message.Chat.ID, "Please send your information in the format: Name YYYY-MM-DD")
		bc.bot.Send(msg)
	default:
		msg := tgbotapi.NewMessage(message.Chat.ID, "I don't know that command")
		bc.bot.Send(msg)
	}
}

func (bc *BotController) handleNonCommand(ctx context.Context, message *tgbotapi.Message) {
	if _, waiting := bc.pendingAddUser[message.Chat.ID]; waiting {
		bc.processAddUser(ctx, message)
		delete(bc.pendingAddUser, message.Chat.ID)
	} else {
		msg := tgbotapi.NewMessage(message.Chat.ID, "I only understand commands")
		bc.bot.Send(msg)
	}
}

func (bc *BotController) processAddUser(ctx context.Context, message *tgbotapi.Message) {
	parts := strings.Split(message.Text, " ")
	if len(parts) != 2 {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Invalid format. Use: Name YYYY-MM-DD")
		bc.bot.Send(msg)
		return
	}

	name := parts[0]
	birthday, err := time.Parse("2006-01-02", parts[1])
	if err != nil {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Invalid date format. Use: YYYY-MM-DD")
		bc.bot.Send(msg)
		return
	}

	user := models.User{
		Name:     name,
		Birthday: birthday,
		ChatID:   message.Chat.ID, // Сохраняем chat_id
	}

	err = bc.userService.AddUser(ctx, user)
	if err != nil {
		log.Printf("Error adding user: %v", err) // Логирование ошибки
		msg := tgbotapi.NewMessage(message.Chat.ID, "Error adding user")
		bc.bot.Send(msg)
		return
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, "User added successfully!")
	bc.bot.Send(msg)
}
