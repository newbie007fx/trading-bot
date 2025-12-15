package service

import (
	"github.com/newbie007fx/trading-bot/internal/domain"
)

func Simulate(
	state *domain.BotState,
	action domain.Action,
	price float64,
) {
	switch action {

	case domain.ActionBuy:
		if state.CashBalance <= 0 {
			return
		}

		state.Position = "LONG"
		state.EntryPrice = price
		state.PositionSize = state.CashBalance / price
		state.CashBalance = 0

	case domain.ActionSell:
		if state.Position != "LONG" {
			return
		}

		value := state.PositionSize * price
		pnl := value - (state.PositionSize * state.EntryPrice)

		if pnl > 0 {
			state.WinTrades++
		} else {
			state.LossTrades++
		}

		state.CashBalance = value
		state.Position = "NONE"
		state.PositionSize = 0
		state.TotalTrades++
	}

	// update equity
	ethValue := state.PositionSize * price
	state.Equity = state.CashBalance + ethValue
}
