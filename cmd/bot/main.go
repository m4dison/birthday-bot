package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/lib/pq"
	"github.com/m4dison/my-telegram-bot/internal/controller"
	"github.com/m4dison/my-telegram-bot/internal/repository/postgres"
	"github.com/m4dison/my-telegram-bot/internal/service"
)

func main() {
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
	dataSourceName := "user=botadmin password=admin dbname=botdatabase sslmode=disable"

	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	repo, err := postgres.NewUserStore(dataSourceName)
	if err != nil {
		log.Fatal(err)
	}
	userService := service.NewUserService(repo, &mu)
	botAdapter := service.NewTelegramBotAdapter(bot)
	botController := controller.NewBotController(bot, userService)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Создаем канал для сигналов
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Канал для завершения работы горутин
	done := make(chan struct{})

	go func() {
		botController.Start(ctx)
		done <- struct{}{}
	}()

	go func() {
		service.StartBirthdayNotifier(ctx, botAdapter, userService)
		done <- struct{}{}
	}()

	sig := <-sigChan
	log.Printf("Received signal: %s. Shutting down gracefully...", sig)
	cancel()

	// Ждем завершения всех горутин
	<-done
	<-done

	log.Println("Application stopped gracefully")
}
