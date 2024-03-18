package model

import "time"

type Client struct {
	ID      string    `json:"id"`
	UserID  string    `json:"user_id"`
	SyncTms time.Time `json:"sync_tms"`
}

func NewClient(userID string, syncTms time.Time) Client {
	return Client{
		UserID:  userID,
		SyncTms: syncTms,
	}
}
