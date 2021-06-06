package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"telebot-trading/utils"
)

const BASE_URL string = "https://api.telegram.org/bot"

type sendMessageReqBody struct {
	ChatID int64  `json:"chat_id"`
	Text   string `json:"text"`
}

func SendToTelegram(chatID int64, text string) error {
	reqBody := &sendMessageReqBody{
		ChatID: chatID,
		Text:   text,
	}

	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	resp, err := http.Post(
		buildSendMessageUrl(),
		"application/json",
		bytes.NewBuffer(reqBytes),
	)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("unexpected status" + resp.Status)
	}

	return err
}

func buildSendMessageUrl() string {
	token := utils.Env("TELEGRAM_TOKEN", "")
	return BASE_URL + token + "/sendMessage"
}
