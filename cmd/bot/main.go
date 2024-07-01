package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/m4dison/my-telegram-bot/internal/controller"
	"github.com/m4dison/my-telegram-bot/internal/repository/memory"
	"github.com/m4dison/my-telegram-bot/internal/service"
)

const pidFilePath = "/Users/madisagi/go/src/github.com/m4dison/my-telegram-bot/bot.pid"

func main() {
	if checkIfAlreadyRunning() {
		log.Fatal("Bot is already running")
	}

	defer removePIDFile()

	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN is not set")
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	var mu sync.Mutex
	repo := memory.NewUserStore()
	userService := service.NewUserService(repo, &mu)
	botAdapter := service.NewTelegramBotAdapter(bot)
	botController := controller.NewBotController(bot, userService)

	go botController.Start()
	go service.StartBirthdayNotifier(botAdapter, userService)

	handleSignals()

	select {}
}

func checkIfAlreadyRunning() bool {
	if _, err := os.Stat(pidFilePath); err == nil {
		log.Printf("PID file %s exists, bot is already running", pidFilePath)
		return true
	}

	log.Printf("PID file %s does not exist, creating new PID file", pidFilePath)

	pid := os.Getpid()
	pidFile, err := os.Create(pidFilePath)
	if err != nil {
		log.Fatalf("Failed to create PID file: %v", err)
	}
	defer pidFile.Close()

	_, err = pidFile.WriteString(fmt.Sprintf("%d", pid))
	if err != nil {
		log.Fatalf("Failed to write PID to file: %v", err)
	}

	log.Printf("PID file %s created with PID %d", pidFilePath, pid)

	return false
}

func removePIDFile() {
	err := os.Remove(pidFilePath)
	if err != nil {
		log.Printf("Failed to remove PID file: %v", err)
	} else {
		log.Println("PID file removed successfully")
	}
}

func handleSignals() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-c
		log.Println("Received termination signal, exiting...")
		cleanup() // Добавить функцию очистки ресурсов
		os.Exit(0)
	}()
}

func cleanup() {
	removePIDFile()
	// Другие действия по очистке ресурсов
}
