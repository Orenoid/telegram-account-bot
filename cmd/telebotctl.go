package main

import (
	"fmt"
	"github.com/orenoid/telegram-account-bot/conf"
	billdal "github.com/orenoid/telegram-account-bot/dal/bill"
	teledal "github.com/orenoid/telegram-account-bot/dal/telegram"
	userdal "github.com/orenoid/telegram-account-bot/dal/user"
	billservice "github.com/orenoid/telegram-account-bot/service/bill"
	teleservice "github.com/orenoid/telegram-account-bot/service/telegram"
	"github.com/orenoid/telegram-account-bot/telebot"
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

		teleRepo, err := teledal.NewMysqlRepo(config.MysqlDSN)
		if err != nil {
			panic(err)
		}
		billRepo, err := billdal.NewMysqlRepo(config.MysqlDSN)
		if err != nil {
			panic(err)
		}
		userRepo, err := userdal.NewMysqlRepo(config.MysqlDSN)
		if err != nil {
			panic(err)
		}

		teleService := teleservice.NewService(teleRepo)
		billService := billservice.NewService(billRepo, userRepo)

		telegramUserStateManager := telebot.NewInMemoryUserStateManager()

		hub := telebot.NewHandlerHub(billService, teleService, telegramUserStateManager)
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
