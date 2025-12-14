package repository

import (
	"context"
	"telebot-trading/internal/domain"

	"cloud.google.com/go/firestore"
)

type StateRepository interface {
	Load(ctx context.Context) (*domain.BotState, error)
	Save(ctx context.Context, state *domain.BotState) error
}

type FirestoreStateRepo struct {
	client     *firestore.Client
	collection string
	document   string
}

func NewFirestoreStateRepo(
	client *firestore.Client,
	collection string,
	document string,
) *FirestoreStateRepo {
	return &FirestoreStateRepo{
		client:     client,
		collection: collection,
		document:   document,
	}
}

func (r *FirestoreStateRepo) Load(ctx context.Context) (*domain.BotState, error) {
	ref := r.client.Collection(r.collection).Doc(r.document)

	snap, err := ref.Get(ctx)
	if err != nil {
		state := domain.NewInitialState()
		_, err := ref.Set(ctx, state)
		return state, err
	}

	var state domain.BotState
	if err := snap.DataTo(&state); err != nil {
		return nil, err
	}
	return &state, nil
}

func (r *FirestoreStateRepo) Save(ctx context.Context, state *domain.BotState) error {
	state.UpdatedAt = state.LastRun
	_, err := r.client.
		Collection(r.collection).
		Doc(r.document).
		Set(ctx, state)

	return err
}
