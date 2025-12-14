package domain

import "time"

type BotState struct {
	Position     string    `firestore:"position"`
	EntryPrice   float64   `firestore:"entry_price"`
	PositionSize float64   `firestore:"position_size"`
	LastAction   string    `firestore:"last_action"`
	LastRun      time.Time `firestore:"last_run"`
	CreatedAt    time.Time `firestore:"created_at"`
	UpdatedAt    time.Time `firestore:"updated_at"`
}

func NewInitialState() *BotState {
	now := time.Now().UTC()
	return &BotState{
		Position:   "NONE",
		LastAction: "INIT",
		LastRun:    now,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}
