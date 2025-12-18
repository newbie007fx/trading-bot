package service

import (
	"github.com/newbie007fx/trading-bot/internal/domain"
)

type StrategyInput struct {
	Price      float64
	EMA50Prev  float64
	EMA50Cur   float64
	EMA200Prev float64
	EMA200Cur  float64
	EMA7Prev   float64
	EMA7Cur    float64
	EMA1Prev   float64
	EMA1Cur    float64
	RSI6Prev   float64
	RSI6Cur    float64
	RSI14Prev  float64
	RSI14Cur   float64
}

func EvaluateStrategy(in StrategyInput, state *domain.BotState) domain.Action {
	if state.Rule != string(domain.Rule1) {
		if in.EMA50Prev < in.EMA200Prev && in.EMA50Cur > in.EMA200Cur {
			if state.Rule == string(domain.Rule2) {
				state.Rule = string(domain.Rule1)
				state.IsAdjusted = true
				return domain.ActionAdjustRule
			}
			state.Rule = string(domain.Rule1)
			state.IsAdjusted = true
			return domain.ActionBuy
		}
	}

	if state.Position == "NONE" {
		if (in.EMA7Prev < in.EMA200Prev && in.EMA7Cur > in.EMA200Cur) ||
			(in.EMA7Prev < in.EMA50Prev && in.EMA7Cur > in.EMA50Cur) {
			if in.RSI6Cur >= 60 || in.RSI14Cur >= 60 {
				state.Rule = string(domain.Rule2)
				return domain.ActionBuy
			}
		}

		if in.RSI6Cur >= 80 && in.RSI14Cur >= 70 {
			state.Rule = string(domain.Rule3)
			return domain.ActionBuy
		}
	}

	if state.Position == "LONG" {
		if state.Rule == string(domain.Rule1) && ((in.EMA1Prev > in.EMA200Prev && in.EMA1Cur < in.EMA200Cur) ||
			(in.EMA1Prev > in.EMA50Prev && in.EMA1Cur < in.EMA50Cur)) {
			return domain.ActionSell
		} else if (state.Rule == string(domain.Rule2) || state.Rule == string(domain.Rule3)) && ((in.EMA1Prev > in.EMA200Prev && in.EMA1Cur < in.EMA200Cur) ||
			(in.EMA1Prev > in.EMA50Prev && in.EMA1Cur < in.EMA50Cur) || (in.EMA1Prev > in.EMA7Prev && in.EMA1Cur < in.EMA7Cur)) {
			return domain.ActionSell
		}

		if in.RSI14Prev >= 60 && in.RSI6Prev >= 60 && in.RSI6Cur < 60 && in.RSI14Cur < 60 {
			if state.IsAdjusted {
				state.IsAdjusted = false
				return domain.ActionAdjustRule
			}

			return domain.ActionSell
		}

		if state.Rule != string(domain.Rule3) && in.RSI6Cur >= 80 && in.RSI14Cur >= 70 && state.LastAction != string(domain.ActionAdjustRule) {
			state.IsAdjusted = true
			return domain.ActionAdjustRule
		}
	}

	return domain.ActionCheck
}
