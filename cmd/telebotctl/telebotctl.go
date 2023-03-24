package main

import (
	"github.com/orenoid/telegram-account-bot/cmd/telebotctl/db"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "",
	Short: "",
}

func init() {
	Cmd.AddCommand(db.Cmd)
}

func main() {
	err := Cmd.Execute()
	if err != nil {
		panic(err)
	}
}
