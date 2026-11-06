package scheduler

import (
	"context"
	"log/slog"
	"time"
)

type Scheduler struct {
	jobs []Job
}

func New() *Scheduler {
	return &Scheduler{}
}

func (s *Scheduler) Add(jobs ...Job) {
	s.jobs = append(s.jobs, jobs...)
}

func (s *Scheduler) Run(ctx context.Context) {
	for {
		now := time.Now()
		var next time.Time
		var nextJob Job
		for _, job := range s.jobs {
			runAt := job.NextRun(now)
			if next.IsZero() || runAt.Before(next) {
				next = runAt
				nextJob = job
			}
		}
		if nextJob == nil {
			<-ctx.Done()
			return
		}

		timer := time.NewTimer(time.Until(next))
		select {
		case <-ctx.Done():
			timer.Stop()
			return
		case <-timer.C:
			if err := nextJob.Run(ctx); err != nil {
				slog.Error("Scheduler job failed", "name", nextJob.Name(), "error", err)
			}
		}
	}
}
