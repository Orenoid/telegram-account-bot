package db

import (
	"fmt"
	"github.com/orenoid/telegram-account-bot/conf"
	"github.com/orenoid/telegram-account-bot/models"

	"github.com/spf13/cobra"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Cmd = &cobra.Command{
	Use:   "db",
	Short: "操作数据库",
}

func init() {
	Cmd.AddCommand(migrateCmd)
}

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Automatically migrate database schema",
	Run: func(cmd *cobra.Command, args []string) {
		config, err := conf.GetConfigFromEnv()
		if err != nil {
			panic(err)
		}
		db, err := gorm.Open(mysql.Open(config.MysqlDSN), &gorm.Config{DisableAutomaticPing: false})

		err = db.AutoMigrate(&models.User{}, &models.Bill{}, &models.TelegramUser{})
		if err != nil {
			panic(err)
		}

		fmt.Println("Database schema migrated successfully")
	},
}
