package bot

import (
	"context"
	"log/slog"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Handler struct {
	bot     *tgbotapi.BotAPI
	chatID  int64
	message string
}

func NewHandler(botAPI *tgbotapi.BotAPI, chatID int64, message string) *Handler {
	return &Handler{bot: botAPI, chatID: chatID, message: message}
}

func (h *Handler) SendDailyReminder(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	msg := tgbotapi.NewMessage(h.chatID, h.message)
	if _, err := h.bot.Send(msg); err != nil {
		slog.Error("SendDailyReminder", "error", err)
		return err
	}
	return nil
}

func (h *Handler) Start(ctx context.Context) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 30

	updates := h.bot.GetUpdatesChan(u)

	for {
		select {

		case <-ctx.Done():
			slog.Info("Stopping telegram handler...")
			h.bot.StopReceivingUpdates()
			return

		case update, ok := <-updates:
			if !ok {
				return
			}
			if update.Message != nil {
				h.handleMessage(update.Message)
			}
		}
	}
}

func (h *Handler) handleMessage(msg *tgbotapi.Message) {
	const fn_name = "Handler.handleMessage"

	if msg.Chat.ID < 0 {
		slog.Info("Received message", "text", msg.Text, "chat_id", msg.Chat.ID)
	}
}
