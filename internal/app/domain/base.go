package domain

import "time"

type Base interface {
	GetID() string
	IsNew() bool
	GetModifiedTms() time.Time
	GetStatus() Status
	SetStatus(status Status)
}
