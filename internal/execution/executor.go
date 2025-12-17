package execution

import (
	"context"

	"github.com/newbie007fx/trading-bot/internal/domain"
)

type Executor interface {
	Buy(ctx context.Context, state *domain.BotState, price float64) error
	Sell(ctx context.Context, state *domain.BotState, price float64) error
	Name() string
}
