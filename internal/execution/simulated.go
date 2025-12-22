package execution

import (
	"context"
	"errors"

	"github.com/newbie007fx/trading-bot/internal/domain"
)

type SimulatedExecutor struct {
}

func NewSimulatedExecutor() *SimulatedExecutor {
	return &SimulatedExecutor{}
}

func (e *SimulatedExecutor) Buy(ctx context.Context, state *domain.BotState, price float64) error {
	if state.CashBalance <= 0 {
		return errors.New("balance is empty")
	}

	state.Position = "LONG"
	state.EntryPrice = price
	state.PositionSize = state.CashBalance / price
	state.CashBalance = 0

	assetValue := state.PositionSize * price
	state.Equity = state.CashBalance + assetValue
	return nil
}

func (e *SimulatedExecutor) Sell(ctx context.Context, state *domain.BotState, price float64) error {
	if state.Position != "LONG" {
		return errors.New("no coin holded")
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

	assetValue := state.PositionSize * price
	state.Equity = state.CashBalance + assetValue
	return nil
}

func (e *SimulatedExecutor) Name() string {
	return "SIMULATION"
}
