package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type AppConfig struct {
	BotToken string   `yaml:"bot_token"`
	BotProxy string   `yaml:"bot_proxy"`
	ChatID   int64    `yaml:"chat_id"`
	Message  string   `yaml:"message"`
	Time     string   `yaml:"time"`
	Timezone string   `yaml:"timezone"`
	Weekdays []string `yaml:"weekdays"`
}

func Default() AppConfig {
	return AppConfig{
		BotToken: "YOUR_TELEGRAM_BOT_TOKEN",
		BotProxy: "socks5://vpn:1080",
		Message:  "Коллеги, пора на дейлик!",
		Time:     "10:29",
		Timezone: "Europe/Moscow",
		Weekdays: []string{"monday", "tuesday", "wednesday", "thursday", "friday"},
	}
}

func LoadOrCreate(path string) (*AppConfig, bool, error) {
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		cfg := Default()
		data, err = yaml.Marshal(cfg)
		if err != nil {
			return nil, false, fmt.Errorf("marshal default config: %w", err)
		}
		if err = os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return nil, false, fmt.Errorf("create config directory: %w", err)
		}
		if err = os.WriteFile(path, data, 0600); err != nil {
			return nil, false, fmt.Errorf("create config file: %w", err)
		}
		return &cfg, true, nil
	}
	if err != nil {
		return nil, false, fmt.Errorf("read config: %w", err)
	}

	var cfg AppConfig
	if err = yaml.Unmarshal(data, &cfg); err != nil {
		return nil, false, fmt.Errorf("parse config: %w", err)
	}
	if err = cfg.Validate(); err != nil {
		return nil, false, err
	}
	return &cfg, false, nil
}

func (c *AppConfig) Validate() error {
	if strings.TrimSpace(c.BotToken) == "" || c.BotToken == "YOUR_TELEGRAM_BOT_TOKEN" {
		return errors.New("заполните bot_token в config.yaml")
	}
	if c.ChatID == 0 {
		return errors.New("заполните chat_id в config.yaml")
	}
	if strings.TrimSpace(c.Message) == "" {
		return errors.New("message в config.yaml не должен быть пустым")
	}
	if _, err := time.Parse("15:04", c.Time); err != nil {
		return fmt.Errorf("некорректное время %q: используйте формат HH:MM", c.Time)
	}
	if _, err := time.LoadLocation(c.Timezone); err != nil {
		return fmt.Errorf("не удалось загрузить часовой пояс %q: %w", c.Timezone, err)
	}
	if len(c.Weekdays) == 0 {
		return errors.New("укажите хотя бы один день в weekdays")
	}
	for _, day := range c.Weekdays {
		if _, ok := parseWeekday(day); !ok {
			return fmt.Errorf("неизвестный день недели %q в weekdays", day)
		}
	}
	return nil
}

func parseWeekday(day string) (time.Weekday, bool) {
	days := map[string]time.Weekday{
		"sunday": time.Sunday, "monday": time.Monday, "tuesday": time.Tuesday,
		"wednesday": time.Wednesday, "thursday": time.Thursday,
		"friday": time.Friday, "saturday": time.Saturday,
	}
	weekday, ok := days[strings.ToLower(strings.TrimSpace(day))]
	return weekday, ok
}

func (c *AppConfig) WeekdaySet() map[time.Weekday]bool {
	weekdays := make(map[time.Weekday]bool, len(c.Weekdays))
	for _, day := range c.Weekdays {
		weekday, ok := parseWeekday(day)
		if ok {
			weekdays[weekday] = true
		}
	}
	return weekdays
}
