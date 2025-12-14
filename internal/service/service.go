package service

import (
	"context"
	"log"
	"telebot-trading/internal/repository"
	"time"
)

type BotService struct {
	repo repository.StateRepository
}

func NewBotService(repo repository.StateRepository) *BotService {
	return &BotService{repo: repo}
}

func (s *BotService) Run(ctx context.Context) error {
	state, err := s.repo.Load(ctx)
	if err != nil {
		return err
	}

	log.Printf("Current state: %+v", state)

	state.LastRun = time.Now().UTC()
	state.LastAction = "CHECK"

	return s.repo.Save(ctx, state)
}
