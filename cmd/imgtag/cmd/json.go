package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

// jsonCmd represents the json command
var jsonCmd = &cobra.Command{
	Use:     "json",
	Short:   "write meta to json",
	Aliases: []string{"j"},
	Run: func(cmd *cobra.Command, args []string) {
		JSON = true
		var err error
		if batchOutput != "" {
			err = writeMetaBatch(args)
		} else {
			err = writeMeta(args)
		}
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(jsonCmd)
}
