package service

import (
	"context"
	"log/slog"
	"time"

	"github.com/newbie007fx/trading-bot/internal/domain"
	"github.com/newbie007fx/trading-bot/internal/execution"
	"github.com/newbie007fx/trading-bot/internal/indicator"
	"github.com/newbie007fx/trading-bot/internal/market"
	"github.com/newbie007fx/trading-bot/internal/model"
	"github.com/newbie007fx/trading-bot/internal/repository"
)

type BotService struct {
	repo           repository.StateRepository
	binanceAdapter *market.BinanceAdapter
	executor       execution.Executor
}

func NewBotService(repo repository.StateRepository, binanceAdapter *market.BinanceAdapter, executor execution.Executor) *BotService {
	return &BotService{
		repo:           repo,
		binanceAdapter: binanceAdapter,
		executor:       executor,
	}
}

// ProcessCandles evaluates the strategy based on the provided candles and executes actions
func (s *BotService) ProcessCandles(ctx context.Context, candles []model.CandleData, state *domain.BotState) error {
	lastCandle := candles[len(candles)-1]
	percentDecreaseFromHight := ((lastCandle.Hight - lastCandle.Close) / (lastCandle.Hight - lastCandle.Open)) * 100

	closes := indicator.ExtractClosePrices(candles)

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
	input := StrategyInput{
		Price:                    closes[len(closes)-1],
		EMA7Prev:                 ema7Last2[0],
		EMA1Prev:                 ema1Last2[0],
		EMA50Prev:                ema50Last2[0],
		EMA200Prev:               ema200Last2[0],
		EMA1Cur:                  ema1Last2[1],
		EMA7Cur:                  ema7Last2[1],
		EMA50Cur:                 ema50Last2[1],
		EMA200Cur:                ema200Last2[1],
		RSI6Prev:                 rsi6Last2[0],
		RSI14Prev:                rsi14Last2[0],
		RSI6Cur:                  rsi6Last2[1],
		RSI14Cur:                 rsi14Last2[1],
		PercentDecreaseFromHight: percentDecreaseFromHight,
	}

	action := EvaluateStrategy(input, state)

	slog.Info("running action",
		"executor", s.executor.Name(),
		"action", action,
		"price", input.Price,
		"rule", state.Rule,
		"target_price", state.TargetPrice,
		"is_adjusted", state.IsAdjusted)
	slog.Info("indicator values",
		"ema1prev", input.EMA1Prev,
		"ema7prev", input.EMA7Prev,
		"ema50prev", input.EMA50Prev,
		"ema200prev", input.EMA200Prev,
		"ema1cur", input.EMA1Cur,
		"ema7cur", input.EMA7Cur,
		"ema50cur", input.EMA50Cur,
		"ema200cur", input.EMA200Cur)

	switch action {
	case domain.ActionBuy:
		if state.Position == "NONE" {
			err := s.executor.Buy(ctx, state, input.Price)
			if err != nil {
				return err
			}
		}

	case domain.ActionSell:
		if state.Position == "LONG" {
			err := s.executor.Sell(ctx, state, input.Price)
			if err != nil {
				return err
			}
		}
	}

	if action != domain.ActionCheck {
		state.LastAction = string(action)
	}

	return nil
}

func (s *BotService) Run(ctx context.Context) error {
	state, err := s.repo.Load(ctx)
	if err != nil {
		return err
	}

	candles, err := s.binanceAdapter.GetCandles(ctx,
		"12h",
		1000, nil)
	if err != nil {
		return err
	}

	err = s.ProcessCandles(ctx, candles, state)
	if err != nil {
		return err
	}

	state.LastRun = time.Now().UTC()

	return s.repo.Save(ctx, state)
}
