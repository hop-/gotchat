package cmd

import "github.com/spf13/cobra"

var (
	clientCmd = &cobra.Command{
		Use:   "client",
		Short: "Run the client command",
		Run: func(cmd *cobra.Command, args []string) {
			executeClient()
		},
	}
)

func executeClient() {
	// TODO
}
