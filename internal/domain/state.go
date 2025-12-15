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

	CashBalance float64 `firestore:"cash_balance"`
	Equity      float64 `firestore:"equity"`
	TotalTrades int     `firestore:"total_trades"`
	WinTrades   int     `firestore:"win_trades"`
	LossTrades  int     `firestore:"loss_trades"`

	TargetPrice float64 `firestore:"loss_trades"`
	IsAdjusted  bool    `firestore:"is_adjusted"`

	Rule string `firestore:"rule"`
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
