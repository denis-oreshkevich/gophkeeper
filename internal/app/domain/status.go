package domain

// Status field to represent state of record.
type Status string

const (
	StatusActive Status = "ACTIVE"

	StatusDeleted Status = "DELETED"
)
