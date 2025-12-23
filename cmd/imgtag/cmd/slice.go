package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

// sliceCmd represents the slice command
var sliceCmd = &cobra.Command{
	Use:     "slice",
	Short:   "output meta as a slice",
	Aliases: []string{"s"},
	Run: func(cmd *cobra.Command, args []string) {
		metas, err := metaSlice(args)
		if err != nil {
			log.Fatal(err)
		}
		out := outputFile + getExt()
		err = saveMeta(out, metas)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(sliceCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// sliceCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// sliceCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
