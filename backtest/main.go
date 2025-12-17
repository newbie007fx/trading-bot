package main

import (
	"context"
	"fmt"
	"log"

	"github.com/newbie007fx/trading-bot/internal/domain"
	"github.com/newbie007fx/trading-bot/internal/execution"
	"github.com/newbie007fx/trading-bot/internal/indicator"
	"github.com/newbie007fx/trading-bot/internal/market"
	"github.com/newbie007fx/trading-bot/internal/model"
	"github.com/newbie007fx/trading-bot/internal/service"
)

type BacktestResult struct {
	FinalEquity float64
	Trades      int
	Wins        int
	Losses      int
}

func RunBacktest(
	candles []model.CandleData,
	state *domain.BotState,
) (*BacktestResult, error) {

	simulation := execution.NewSimulatedExecutor()
	ctx := context.Background()

	for i := 200; i < len(candles); i++ {

		window := candles[:i]

		closes := indicator.ExtractClosePrices(window)

		ema7Series, _ := indicator.EMASeries(closes, 7)
		ema50Series, _ := indicator.EMASeries(closes, 50)
		ema200Series, _ := indicator.EMASeries(closes, 200)

		rsi6Series, _ := indicator.RSISeries(closes, 6)
		rsi14Series, _ := indicator.RSISeries(closes, 14)

		ema7Last2 := indicator.LastN(ema7Series, 2)
		ema50Last2 := indicator.LastN(ema50Series, 2)
		ema200Last2 := indicator.LastN(ema200Series, 2)
		rsi6Last2 := indicator.LastN(rsi6Series, 2)
		rsi14Last2 := indicator.LastN(rsi14Series, 2)

		ema1Last2 := indicator.LastN(closes, 2)
		input := service.StrategyInput{
			Price:      closes[len(closes)-1],
			EMA7Prev:   ema7Last2[0],
			EMA1Prev:   ema1Last2[0],
			EMA50Prev:  ema50Last2[0],
			EMA200Prev: ema200Last2[0],
			EMA1Cur:    ema1Last2[1],
			EMA7Cur:    ema7Last2[1],
			EMA50Cur:   ema50Last2[1],
			EMA200Cur:  ema200Last2[1],
			RSI6Prev:   rsi6Last2[0],
			RSI14Prev:  rsi14Last2[0],
			RSI6Cur:    rsi6Last2[1],
			RSI14Cur:   rsi14Last2[1],
		}

		action := service.EvaluateStrategy(input, state)

		switch action {
		case domain.ActionBuy:
			if state.Position == "NONE" {
				_ = simulation.Buy(ctx, state, input.Price)
			}

		case domain.ActionSell:
			if state.Position == "LONG" {
				_ = simulation.Sell(ctx, state, input.Price)
			}
		}
		if action != domain.ActionCheck {
			state.LastAction = string(action)
		}
	}

	log.Printf(
		"[BACKTEST DONE] Equity=%.2f Trades=%d Win=%d Loss=%d",
		state.Equity,
		state.TotalTrades,
		state.WinTrades,
		state.LossTrades,
	)

	return &BacktestResult{
		FinalEquity: state.Equity,
		Trades:      state.TotalTrades,
		Wins:        state.WinTrades,
		Losses:      state.LossTrades,
	}, nil
}

func main() {
	ctx := context.Background()

	marketClient := market.NewBinanceAdapter(ctx, "", "", "ETHUSDT")
	var endTime *int64 = nil
	state := &domain.BotState{
		Position:    "NONE",
		CashBalance: 300,
		Equity:      300,
	}
	for range 3 {
		candles, err := marketClient.GetCandles(ctx,
			"1d",
			1000, endTime)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(candles[0].OpenTime)
		RunBacktest(candles, state)
		endTime = &candles[0].OpenTime
	}
}
