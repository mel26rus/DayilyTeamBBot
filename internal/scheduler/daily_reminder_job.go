package scheduler

import (
	"context"
	"log/slog"
	"time"

	"expense-bot/internal/bot"
)

type DailyReminderJob struct {
	Handler  *bot.Handler
	Time     string
	Location *time.Location
	Weekdays map[time.Weekday]bool
}

func NewDailyReminderJob(handler *bot.Handler, runTime string, location *time.Location, weekdays map[time.Weekday]bool) *DailyReminderJob {
	return &DailyReminderJob{Handler: handler, Time: runTime, Location: location, Weekdays: weekdays}
}

func (j *DailyReminderJob) Name() string {
	return "Daily Reminder Job"
}

/*
func (j *DailyReminderJob) NextRun(now time.Time) time.Time {
	localNow := now.In(j.Location)
	parsed, _ := time.Parse("15:04", j.Time)
	hour, minute, _ := parsed.Clock()
	for dayOffset := 0; dayOffset < 8; dayOffset++ {
		date := localNow.AddDate(0, 0, dayOffset)
		candidate := time.Date(date.Year(), date.Month(), date.Day(), hour, minute, 0, 0, j.Location)
		if j.Weekdays[candidate.Weekday()] && candidate.After(localNow) {
			slog.Debug("DailyReminderJob.NextRun", "name", j.Name(), "next", candidate)
			return candidate
		}
	}
	return time.Time{}
}
*/

func (j *DailyReminderJob) NextRun(now time.Time) time.Time {
	localNow := now.In(j.Location)
	parsed, _ := time.Parse("15:04", j.Time)
	hour, minute, _ := parsed.Clock()

	// Отступаем на секунду в прошлое, чтобы поймать границу и избежать микросекундных гонок
	// Но для защиты от двойного запуска внутри одной минуты лучше проверять с запасом:
	for dayOffset := 0; dayOffset < 8; dayOffset++ {
		date := localNow.AddDate(0, 0, dayOffset)
		candidate := time.Date(date.Year(), date.Month(), date.Day(), hour, minute, 0, 0, j.Location)

		// Если время кандидата совпадает с текущей минутой или уже прошло — пропускаем этот день
		if candidate.Sub(localNow) <= time.Second {
			continue
		}

		if j.Weekdays[candidate.Weekday()] {
			slog.Debug("DailyReminderJob.NextRun", "name", j.Name(), "next", candidate)
			return candidate
		}
	}
	return time.Time{}
}

func (j *DailyReminderJob) Run(ctx context.Context) error {
	slog.Debug("+DailyReminderJob.Run")
	defer slog.Debug("-DailyReminderJob.Run")
	return j.Handler.SendDailyReminder(ctx)
}
