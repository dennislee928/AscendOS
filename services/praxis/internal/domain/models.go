package domain

import "time"

// Playbook defines a reusable operation workflow in Praxis.
type Playbook struct {
	ID          string
	Name        string
	Description string
	Version     string
}

// Execution tracks a playbook run and outcome.
type Execution struct {
	ID         string
	PlaybookID string
	Status     string
	StartedAt  time.Time
	EndedAt    *time.Time
}
