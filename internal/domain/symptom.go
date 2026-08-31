package domain

type Symptom struct {
	ID          int64
	Code        string
	Name        string
	Description string
	CategoryID  int64
}

type Category struct {
	ID   int64
	Name string
}
