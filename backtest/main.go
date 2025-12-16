package main

import (
	"context"
	"fmt"
	"log"

	"github.com/newbie007fx/trading-bot/internal/domain"
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
	initialCash float64,
) (*BacktestResult, error) {

	state := &domain.BotState{
		Position:    "NONE",
		CashBalance: initialCash,
		Equity:      initialCash,
	}

	for i := 200; i < len(candles); i++ {

		window := candles[:i]

		closes := indicator.ExtractClosePrices(window)

		ema7Series, _ := indicator.EMASeries(closes, 7)
		ema50Series, _ := indicator.EMASeries(closes, 50)
		ema200Series, _ := indicator.EMASeries(closes, 200)

		ema7Last2 := indicator.LastN(ema7Series, 2)
		ema50Last2 := indicator.LastN(ema50Series, 2)
		ema200Last2 := indicator.LastN(ema200Series, 2)

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
		}

		action := service.EvaluateStrategy(input, state)

		service.Simulate(state, action, input.Price)
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

	marketClient := market.NewBinanceAdapter("", "")
	var endTime *int64 = nil
	for range 3 {
		candles, err := marketClient.GetCandles(ctx,
			"ETHUSDT",
			"1d",
			1000, endTime)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(candles[0].OpenTime)
		RunBacktest(candles, 300)
		endTime = &candles[0].OpenTime
	}
}
