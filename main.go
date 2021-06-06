package main

import (
	_ "telebot-trading/app/console"
	"telebot-trading/bootstrap"
	_ "telebot-trading/database"
	_ "telebot-trading/external"
)

func main() {
	bootstraper := bootstrap.Bootstraper{}
	bootstraper.RegistServices()
	defer bootstraper.ShutdownServices()
	bootstraper.Run()
}
