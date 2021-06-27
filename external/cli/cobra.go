package cli

import (
	"github.com/spf13/cobra"
)

var cmdList []*cobra.Command

func RegisterCommand(cmd *cobra.Command) {
	cmdList = append(cmdList, cmd)
}

type ConsoleService struct{}

func (ConsoleService) rootCmd() *cobra.Command {
	rootCmd := &cobra.Command{}
	rootCmd.Use = "telebot"
	rootCmd.Short = "A Simple Web"
	rootCmd.Long = `A simple web admin application for telebot`
	return rootCmd
}

func (cs *ConsoleService) Run() error {
	rootCmd := cs.rootCmd()
	for _, cmd := range cmdList {
		rootCmd.AddCommand(cmd)
	}

	return rootCmd.Execute()
}
