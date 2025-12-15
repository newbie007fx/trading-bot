package function

import (
	"log"
	"net/http"

	"github.com/newbie007fx/trading-bot/internal/config"
	"github.com/newbie007fx/trading-bot/internal/market"
	"github.com/newbie007fx/trading-bot/internal/repository"
	"github.com/newbie007fx/trading-bot/internal/service"

	"cloud.google.com/go/firestore"
)

// helloHTTP is an HTTP Cloud Function with a request parameter.
func ExecuteBot(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	cfg := config.Load()

	client, err := firestore.NewClientWithDatabase(ctx, cfg.ProjectID, cfg.DatabaseID)
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

	marketClient := market.NewBinanceAdapter(ctx, cfg)
	bot := service.NewBotService(repo, marketClient)

	if err := bot.Run(ctx); err != nil {
		log.Println("Bot error:", err)
		http.Error(w, "bot error", 500)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}
