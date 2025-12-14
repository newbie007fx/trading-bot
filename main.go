package main

import (
	"fmt"
	"log"
	"net/http"
	"telebot-trading/internal/config"
	"telebot-trading/internal/indicator"
	"telebot-trading/internal/market"
	"telebot-trading/internal/repository"
	"telebot-trading/internal/service"

	"cloud.google.com/go/firestore"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

func init() {
	functions.HTTP("HelloHTTP", helloHTTP)
}

// helloHTTP is an HTTP Cloud Function with a request parameter.
func helloHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	cfg := config.Load()

	client, err := firestore.NewClient(ctx, cfg.ProjectID)
	if err != nil {
		http.Error(w, "firestore error", 500)
		log.Println(err)
		return
	}
	defer client.Close()

	repo := repository.NewFirestoreStateRepo(
		client,
		cfg.Collection,
		cfg.DocumentID,
	)

	bot := service.NewBotService(repo)

	if err := bot.Run(ctx); err != nil {
		log.Println("Bot error:", err)
		http.Error(w, "bot error", 500)
		return
	}

	marketClient := market.NewBinanceAdapter(ctx, cfg)
	candles, err := marketClient.GetCandles(ctx,
		"ETHUSDT",
		"1d",
		250)
	if err != nil {
		http.Error(w, "error getting data", 500)
		return
	}

	closes := indicator.ExtractClosePrices(candles)

	ema50Series, _ := indicator.EMASeries(closes, 50)
	ema200Series, _ := indicator.EMASeries(closes, 200)

	ema50Last2 := indicator.LastN(ema50Series, 2)
	ema200Last2 := indicator.LastN(ema200Series, 2)

	if ema50Last2 == nil || ema200Last2 == nil {
		http.Error(w, "not enough EMA data", 500)
		return
	}

	result := fmt.Sprintf(
		"EMA50 prev=%.2f curr=%.2f | EMA200 prev=%.2f curr=%.2f",
		ema50Last2[0], ema50Last2[1],
		ema200Last2[0], ema200Last2[1],
	)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(result))
}
