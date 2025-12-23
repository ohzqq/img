package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

// yamlCmd represents the yaml command
var yamlCmd = &cobra.Command{
	Use:     "yaml",
	Aliases: []string{"y"},
	Short:   "meta to yaml",
	Run: func(cmd *cobra.Command, args []string) {
		YAML = true
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
	rootCmd.AddCommand(yamlCmd)
}
