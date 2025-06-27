package cmd

import (
	"github.com/spf13/cobra"
)

var (
	port int = 7665

	rootCmd = &cobra.Command{
		Use:   "gotchat",
		Short: "A simple chat application",
		Long:  `A simple terminal based chat application built with Go.`,
		Run: func(cmd *cobra.Command, args []string) {
			executeApp()
		},
	}
)

// autorun: This function is called by the main package to initialize the command line interface.
func init() {
	// Flags for root and "app"
	rootCmd.Flags().IntVarP(&port, "port", "p", 7665, "port on which connection listener will be started")

	// Add subcommands
	rootCmd.AddCommand(appCmd)
	rootCmd.AddCommand(clientCmd)
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
