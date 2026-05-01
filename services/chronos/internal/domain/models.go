package domain

import "time"

// Schedule is a recurring job definition managed by Chronos.
type Schedule struct {
	ID        string
	Name      string
	CronExpr  string
	Enabled   bool
	NextRunAt time.Time
}

// Run captures one execution attempt for a schedule.
type Run struct {
	ID         string
	ScheduleID string
	Status     string
	StartedAt  time.Time
	FinishedAt *time.Time
}
