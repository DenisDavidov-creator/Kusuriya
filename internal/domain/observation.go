package domain

import "time"

type Observation struct {
	ID         int64
	CaseID     int64
	Kind       string
	RefID      int64
	ObservedAt time.Time
	Note       string
}
