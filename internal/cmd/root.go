package cmd

import (
	"os"

	"github.com/hop-/gotchat/internal/config"
	"github.com/hop-/gotchat/pkg/log"
	"github.com/spf13/cobra"
)

var (
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
	// Flags for root command
	rootCmd.Flags().IntVarP(
		&generalServerPort,
		"port", "p",
		config.GetServerPort(),
		"port on which connection listener will be started",
	)
	rootCmd.Flags().StringVarP(
		&generalDataStorageFile,
		"storage", "s",
		config.GetDataStorageFilePath(),
		"file to store chat data and configurations",
	)

	// Add subcommands
	rootCmd.AddCommand(appCmd)
	rootCmd.AddCommand(clientCmd)
	rootCmd.AddCommand(versionCmd)
}

func Execute() {
	log.Configure().
		InMemory().
		StdOut().
		Level(log.DEBUG). // TODO: set log level from config
		Init()
	defer log.Close()

	err := createRootDirIfNotExists()
	if err != nil {
		log.Fatalf("Failed to create root directory: %v", err)
	}
	cobra.CheckErr(rootCmd.Execute())

}

func createRootDirIfNotExists() error {
	rootDir := config.GetRootDir()

	return os.MkdirAll(rootDir, 0755)
}
