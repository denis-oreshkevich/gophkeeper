package model

import "time"

type BinarySync struct {
	LastSyncTms time.Time `json:"last_sync_tms"`
	Binaries    []*Binary `json:"binaries"`
}

func NewBinarySync(lastSyncTms time.Time, binaries []*Binary) BinarySync {
	return BinarySync{
		LastSyncTms: lastSyncTms,
		Binaries:    binaries,
	}
}

type TextSync struct {
	LastSyncTms time.Time `json:"last_sync_tms"`
	Texts       []*Text   `json:"texts"`
}

func NewTextSync(lastSyncTms time.Time, text []*Text) TextSync {
	return TextSync{
		LastSyncTms: lastSyncTms,
		Texts:       text,
	}
}

type CredSync struct {
	LastSyncTms time.Time      `json:"last_sync_tms"`
	Credentials []*Credentials `json:"credentials"`
}

func NewCredSync(lastSyncTms time.Time, credentials []*Credentials) CredSync {
	return CredSync{
		LastSyncTms: lastSyncTms,
		Credentials: credentials,
	}
}

type CardSync struct {
	LastSyncTms time.Time `json:"last_sync_tms"`
	Cards       []*Card   `json:"cards"`
}

func NewCardSync(lastSyncTms time.Time, card []*Card) CardSync {
	return CardSync{
		LastSyncTms: lastSyncTms,
		Cards:       card,
	}
}
