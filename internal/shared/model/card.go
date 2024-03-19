package model

import "time"

type Card struct {
	ID          string    `json:"id"`
	Num         string    `json:"num"`
	CVC         string    `json:"cvc"`
	HolderName  string    `json:"holder_name"`
	New         bool      `json:"-"`
	UserID      string    `json:"user_id"`
	Status      Status    `json:"status"`
	ModifiedTms time.Time `json:"modified_tms"`
}

func NewCard(num string, cvc string, holderName string, userID string,
	status Status, modifiedTms time.Time) Card {
	return Card{
		Num:         num,
		CVC:         cvc,
		HolderName:  holderName,
		UserID:      userID,
		Status:      status,
		ModifiedTms: modifiedTms,
	}
}

func (c *Card) GetID() string {
	return c.ID
}

func (c *Card) IsNew() bool {
	return c.New
}

func (c *Card) GetModifiedTms() time.Time {
	return c.ModifiedTms
}

func (c *Card) GetStatus() Status {
	return c.Status
}

func (c *Card) SetStatus(status Status) {
	c.Status = status
}
