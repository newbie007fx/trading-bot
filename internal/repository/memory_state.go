package repository

import (
	"context"

	"github.com/newbie007fx/trading-bot/internal/domain"
)

type MemoryStateRepo struct {
	state *domain.BotState
}

func NewMemoryStateRepo() *MemoryStateRepo {
	return &MemoryStateRepo{}
}

func (r *MemoryStateRepo) Load(ctx context.Context) (*domain.BotState, error) {
	if r.state == nil {
		return domain.NewInitialState(), nil
	}
	return r.state, nil
}

func (r *MemoryStateRepo) Save(ctx context.Context, state *domain.BotState) error {
	r.state = state
	return nil
}
