package execution

import (
	"context"

	"github.com/newbie007fx/trading-bot/internal/domain"
	"github.com/newbie007fx/trading-bot/internal/market"
)

type LiveExecutor struct {
	binanceAdapter *market.BinanceAdapter
}

func NewLiveExecutor(binanceAdapter *market.BinanceAdapter) *LiveExecutor {
	return &LiveExecutor{binanceAdapter: binanceAdapter}
}

func (e *LiveExecutor) Buy(ctx context.Context, state *domain.BotState, price float64) error {
	err := e.binanceAdapter.Buy(ctx, price)
	if err != nil {
		return err
	}

	state.Position = "LONG"
	state.EntryPrice = price
	state.PositionSize = state.CashBalance / price
	state.CashBalance = 0

	ethValue := state.PositionSize * price
	state.Equity = state.CashBalance + ethValue

	return nil
}

func (e *LiveExecutor) Sell(ctx context.Context, state *domain.BotState, price float64) error {
	err := e.binanceAdapter.Sell(ctx)
	if err != nil {
		return err
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

	ethValue := state.PositionSize * price
	state.Equity = state.CashBalance + ethValue

	return nil
}

func (e *LiveExecutor) Name() string {
	return "LIVE"
}
