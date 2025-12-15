package service

import "github.com/newbie007fx/trading-bot/internal/domain"

type StrategyInput struct {
	Price      float64
	EMA50Prev  float64
	EMA50Cur   float64
	EMA200Prev float64
	EMA200Cur  float64
	EMA7Prev   float64
	EMA7Cur    float64
}

func EvaluateStrategy(in StrategyInput, state *domain.BotState) domain.Action {
	if state.Rule != string(domain.Rule1) {
		if in.EMA50Prev < in.EMA200Prev && in.EMA50Cur > in.EMA200Cur {
			state.Rule = string(domain.Rule1)
			if state.Rule == string(domain.Rule2) {
				target := state.EntryPrice + state.EntryPrice*15/100
				if target > state.TargetPrice {
					state.TargetPrice = target
				}
				return domain.ActionAdjustRule
			}
			state.TargetPrice = in.Price + 2/100*in.Price
			return domain.ActionBuy
		}
	}

	if state.Position == "NONE" {
		if (in.EMA7Prev < in.EMA200Prev && in.EMA7Cur > in.EMA200Cur) ||
			(in.EMA7Prev < in.EMA50Prev && in.EMA7Cur > in.EMA50Cur) {
			state.Rule = string(domain.Rule2)
			state.TargetPrice = in.Price + 2/100*in.Price
			return domain.ActionBuy
		}
	}

	if state.Position == "LONG" {
		if in.Price > state.TargetPrice {
			state.TargetPrice = in.Price + 1/100*in.Price
			state.IsAdjusted = true
		} else if state.IsAdjusted && in.Price <= (state.TargetPrice-2/100*state.TargetPrice) {
			state.Rule = ""
			state.TargetPrice = 0
			state.IsAdjusted = false
			return domain.ActionSell
		}
	}

	return domain.ActionCheck
}
