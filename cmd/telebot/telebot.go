package main

import (
	"fmt"
	"github.com/orenoid/telegram-account-bot/conf"
	billdal "github.com/orenoid/telegram-account-bot/dal/bill"
	teledal "github.com/orenoid/telegram-account-bot/dal/telegram"
	userdal "github.com/orenoid/telegram-account-bot/dal/user"
	billservice "github.com/orenoid/telegram-account-bot/service/bill"
	teleservice "github.com/orenoid/telegram-account-bot/service/telegram"
	"github.com/orenoid/telegram-account-bot/service/user"
	"github.com/orenoid/telegram-account-bot/telebot"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	tele "gopkg.in/telebot.v3"
	"time"
)

var telebotCmd = &cobra.Command{
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
		userService := user.NewUserService(userRepo)

		telegramUserStateManager := telebot.NewInMemoryUserStateManager()

		hub := telebot.NewHandlerHub(billService, teleService, userService, telegramUserStateManager)
		bot, err := telebot.NewBot(settings, hub)
		if err != nil {
			panic(err)
		}
		err = bot.SetCommands([]tele.Command{
			{Text: "/help", Description: "查看使用帮助"},
			{Text: "/start", Description: "初始化"},
			{Text: "/day", Description: "今日收支"},
			{Text: "/month", Description: "本月收支"},
			{Text: "/cancel", Description: "取消当前操作"},
			{Text: "/set_keyboard", Description: "设置快捷键盘"},
			{Text: "/set_balance", Description: "设置余额"},
			{Text: "/balance", Description: "查询余额"},
			{Text: "/create_token", Description: "创建用于 OpenAPI 的 token"},
			{Text: "/disable_all_tokens", Description: "废弃所有 token"},
		})
		if err != nil {
			logrus.Warnf("failed to set commands, err: %+v", err)
		}

		fmt.Println("Running telebot with a LongPoller...")
		bot.Start()

	},
}

func main() {
	if err := telebotCmd.Execute(); err != nil {
		panic(err)
	}
}
