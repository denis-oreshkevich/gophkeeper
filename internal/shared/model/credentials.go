package model

import "time"

type Credentials struct {
	ID          string    `json:"id"`
	Login       string    `json:"login"`
	Password    string    `json:"password"`
	New         bool      `json:"-"`
	UserID      string    `json:"user_id"`
	Status      Status    `json:"status"`
	ModifiedTms time.Time `json:"modified_tms"`
}

func NewCredentials(login string, password string, status Status,
	userID string, modifiedTms time.Time) Credentials {
	return Credentials{
		Login:       login,
		Password:    password,
		UserID:      userID,
		Status:      status,
		ModifiedTms: modifiedTms,
	}
}

func (c *Credentials) GetID() string {
	return c.ID
}

func (c *Credentials) IsNew() bool {
	return c.New
}

func (c *Credentials) GetModifiedTms() time.Time {
	return c.ModifiedTms
}

func (c *Credentials) GetStatus() Status {
	return c.Status
}

func (c *Credentials) SetStatus(status Status) {
	c.Status = status
}
