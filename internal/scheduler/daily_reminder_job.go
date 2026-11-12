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

func (j *DailyReminderJob) NextRun(now time.Time) time.Time {
	localNow := now.In(j.Location)
	parsed, _ := time.Parse("15:04", j.Time)
	hour, minute, _ := parsed.Clock()

	next := time.Date(
		localNow.Year(),
		localNow.Month(),
		localNow.Day(),
		hour, minute, 0, 0,
		j.Location,
	)

	if !next.After(localNow) {
		next = next.AddDate(0, 0, 1)
	}

	for i := 0; i < 7; i++ {
		if j.Weekdays[next.Weekday()] {
			slog.Info("DailyReminderJob.NextRun", "name", j.Name(), "next", next)
			return next
		}
		next = next.AddDate(0, 0, 1)
	}
	slog.Info("DailyReminderJob.NextRun", "name", j.Name(), "next time", time.Time{})
	return time.Time{}
}

func (j *DailyReminderJob) Run(ctx context.Context) error {
	slog.Debug("+DailyReminderJob.Run")
	defer slog.Debug("-DailyReminderJob.Run")
	return j.Handler.SendDailyReminder(ctx)
}
