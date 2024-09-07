package main

import (
	"github.com/orenoid/telegram-account-bot/api"
	"github.com/orenoid/telegram-account-bot/conf"
	billdal "github.com/orenoid/telegram-account-bot/dal/bill"
	userdal "github.com/orenoid/telegram-account-bot/dal/user"
	billservice "github.com/orenoid/telegram-account-bot/service/bill"
	"github.com/orenoid/telegram-account-bot/service/user"
	"github.com/spf13/cobra"
)

var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "api - start the api server",
	Run: func(cmd *cobra.Command, args []string) {
		config, err := conf.GetConfigFromEnv()
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

		billService := billservice.NewService(billRepo, userRepo)
		userService := user.NewUserService(userRepo)

		controllersHub := api.NewControllersHub(userService, billService)
		e := api.GetEcho(controllersHub)
		e.Logger.Fatal(e.Start(":1323"))
	},
}

func main() {
	if err := apiCmd.Execute(); err != nil {
		panic(err)
	}
}
