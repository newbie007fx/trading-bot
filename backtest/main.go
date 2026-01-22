package main

import (
	"context"
	"fmt"
	"log"

	"github.com/newbie007fx/trading-bot/internal/domain"
	"github.com/newbie007fx/trading-bot/internal/execution"
	"github.com/newbie007fx/trading-bot/internal/market"
	"github.com/newbie007fx/trading-bot/internal/model"
	"github.com/newbie007fx/trading-bot/internal/repository"
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
	botService *service.BotService,
) (*BacktestResult, error) {

	ctx := context.Background()

	for i := 200; i < len(candles); i++ {
		window := candles[:i]
		err := botService.ProcessCandles(ctx, window, state)
		if err != nil {
			return nil, err
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

	marketClient := market.NewBinanceAdapter(ctx, "", "", "SOL")
	simulation := execution.NewSimulatedExecutor()
	repo := repository.NewMemoryStateRepo()

	state := &domain.BotState{
		Position:    "NONE",
		CashBalance: 300,
		Equity:      300,
	}

	if err := repo.Save(ctx, state); err != nil {
		log.Fatal(err)
	}

	botService := service.NewBotService(repo, marketClient, simulation)

	var endTime *int64 = nil
	for range 3 {
		candles, err := marketClient.GetCandles(ctx,
			"12h",
			1000, endTime)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(candles[0].OpenTime)
		RunBacktest(candles, state, botService)
		endTime = &candles[0].OpenTime
	}
}
