package model

import "time"

type Binary struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Data        string    `json:"data"`
	New         bool      `json:"-"`
	UserID      string    `json:"user_id"`
	Status      Status    `json:"status"`
	ModifiedTms time.Time `json:"modified_tms"`
}

func NewBinary(name string, data string, userID string, status Status,
	modifiedTms time.Time) *Binary {
	return &Binary{
		Name:        name,
		Data:        data,
		UserID:      userID,
		Status:      status,
		ModifiedTms: modifiedTms,
	}
}

func (b *Binary) GetID() string {
	return b.ID
}

func (b *Binary) IsNew() bool {
	return b.New
}

func (b *Binary) GetModifiedTms() time.Time {
	return b.ModifiedTms
}

func (b *Binary) GetStatus() Status {
	return b.Status
}

func (b *Binary) SetStatus(status Status) {
	b.Status = status
}
