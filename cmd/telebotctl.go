package main

import (
	"fmt"
	"github.com/orenoid/account-bot/conf"
	"github.com/orenoid/account-bot/telebot"
	"github.com/spf13/cobra"
	tele "gopkg.in/telebot.v3"
	"time"
)

var cmd = &cobra.Command{
	Use:   "telebotctl",
	Short: "telebotctl - start the telegram bot",
	Run: func(cmd *cobra.Command, args []string) {
		config, err := conf.GetConfigFromEnv()
		if err != nil {
			panic(err)
		}

		settings := tele.Settings{
			Token:  config.TelebotToken,
			Poller: &tele.LongPoller{Timeout: 10 * time.Second},
		}

		hub := telebot.NewHandlerHub()
		bot, err := telebot.NewBot(settings, hub)
		if err != nil {
			panic(err)
		}

		fmt.Println("Running telebot with a LongPoller...")
		bot.Start()

	},
}

func main() {
	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}
