package service

import (
	"context"
	"log"
	"time"

	"github.com/newbie007fx/trading-bot/internal/domain"
	"github.com/newbie007fx/trading-bot/internal/execution"
	"github.com/newbie007fx/trading-bot/internal/indicator"
	"github.com/newbie007fx/trading-bot/internal/market"
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

	action := EvaluateStrategy(input, state)

	log.Printf("[%s] running action %s: price %.2f rule %s target price %.2f, is adjusted %t\n", s.executor.Name(), action, input.Price, state.Rule, state.TargetPrice, state.IsAdjusted)
	log.Printf("ema1prev: %.2f, ema7prev: %.2f, ema50prev: %.2f, ema200prev: %.2f, ema1cur: %.2f, ema7cur: %.2f, ema50cur: %.2f, ema200cur: %.2f\n", input.EMA1Prev, input.EMA7Prev, input.EMA50Prev, input.EMA200Prev, input.EMA1Cur, input.EMA7Cur, input.EMA50Cur, input.EMA200Cur)

	switch action {
	case domain.ActionBuy:
		if state.Position == "NONE" {
			err = s.executor.Buy(ctx, state, input.Price)
			if err != nil {
				return err
			}
		}

	case domain.ActionSell:
		if state.Position == "LONG" {
			err = s.executor.Sell(ctx, state, input.Price)
			if err != nil {
				return err
			}
		}
	}

	state.LastRun = time.Now().UTC()
	if action != domain.ActionCheck {
		state.LastAction = string(action)
	}

	return s.repo.Save(ctx, state)
}
