package domain

import "time"

type SystemStatus struct {
	Seeded       bool
	ServiceName  string
	RecordCounts map[string]int64
	LastSeededAt *time.Time
}

type SeedResult struct {
	Success        bool
	Message        string
	RecordsCreated int64
}
