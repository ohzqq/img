package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

// eachCmd represents the each command
var eachCmd = &cobra.Command{
	Use:   "each",
	Short: "save each image's meta seperately",
	Run: func(cmd *cobra.Command, args []string) {
		err := writeMeta(args)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(eachCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// eachCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// eachCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
