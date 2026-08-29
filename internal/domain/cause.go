package domain

type Cause struct {
	ID          int64
	Code        string
	Name        string
	Type        string
	Prior       float64
	Description string
	RemedyID    *int64
}
