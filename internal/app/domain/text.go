package domain

import "time"

type Text struct {
	ID          string    `json:"id"`
	Txt         string    `json:"txt"`
	New         bool      `json:"-"`
	UserID      string    `json:"user_id"`
	Status      Status    `json:"status"`
	ModifiedTms time.Time `json:"modified_tms"`
}

func NewText(txt string, userID string, status Status, modifiedTms time.Time) *Text {
	return &Text{
		Txt:         txt,
		UserID:      userID,
		Status:      status,
		ModifiedTms: modifiedTms,
	}
}

func (t *Text) GetID() string {
	return t.ID
}

func (t *Text) IsNew() bool {
	return t.New
}

func (t *Text) GetModifiedTms() time.Time {
	return t.ModifiedTms
}

func (t *Text) GetStatus() Status {
	return t.Status
}

func (t *Text) SetStatus(status Status) {
	t.Status = status
}
