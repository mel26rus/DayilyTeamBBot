package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"expense-bot/internal/bot"
	"expense-bot/internal/config"
	"expense-bot/internal/filesystem"
	"expense-bot/internal/scheduler"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const configFile = "config.yaml"

func main() {
	baseDir := filesystem.BaseDir(configFile)
	configPath := filepath.Join(baseDir, configFile)
	cfg, created, err := config.LoadOrCreate(configPath)
	if err != nil {
		slog.Error("Load config", "error", err)
		os.Exit(1)
	}
	if created {
		slog.Info("Создан конфиг. Заполните bot_token и chat_id, затем запустите бота снова", "path", configPath)
		return
	}
	if err = cfg.Validate(); err != nil {
		slog.Error("Config validation failed", "error", err)
		os.Exit(1)
	}

	botAPI, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		slog.Error("Bot initialization failed", "error", err)
		os.Exit(1)
	}

	location, err := time.LoadLocation(cfg.Timezone)
	if err != nil {
		slog.Error("Load timezone", "error", err)
		os.Exit(1)
	}
	handler := bot.NewHandler(botAPI, cfg.ChatID, cfg.Message)
	sch := scheduler.New()
	sch.Add(scheduler.NewDailyReminderJob(handler, cfg.Time, location, cfg.WeekdaySet()))

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	slog.Info("Бот запущен", "username", botAPI.Self.UserName, "timezone", cfg.Timezone, "time", cfg.Time)
	sch.Run(ctx)
	go handler.Start(ctx)
}
