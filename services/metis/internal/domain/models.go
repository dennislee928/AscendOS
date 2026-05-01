package domain

import "time"

// Evidence is a source record used by Metis during analysis.
type Evidence struct {
	ID          string
	Source      string
	ExternalID  string
	CollectedAt time.Time
}

// Assessment stores a normalized risk inference built from evidence.
type Assessment struct {
	ID          string
	EvidenceID  string
	Category    string
	Confidence  float64
	Summary     string
	GeneratedAt time.Time
}
