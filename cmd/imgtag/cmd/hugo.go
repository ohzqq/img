package cmd

import (
	"log"
	"os"
	"strings"

	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

var (
	hugoOutput string
	hugoTitle  string
)

// hugoCmd represents the hugo command
var hugoCmd = &cobra.Command{
	Use:     "hugo",
	Aliases: []string{"h"},
	Short:   "write hugo front matter",
	Run: func(cmd *cobra.Command, args []string) {
		out := lo.KebabCase(strings.ToLower(hugoTitle)) + ".md"
		tags := []string{}
		metas, err := metaSlice(args)
		if err != nil {
			log.Fatal(err)
		}
		for _, params := range metas {
			tags = append(tags, params.Subject...)
		}
		tags = lo.Uniq(tags)
		m := map[string]any{
			"title": hugoTitle,
			"tags":  tags,
			"params": map[string]any{
				"images": metas,
			},
		}
		w, err := os.Create(out)
		if err != nil {
			log.Fatal(err)
		}
		err = encodeMeta(w, m)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(hugoCmd)
	hugoCmd.Flags().StringVarP(&hugoTitle, "name", "n", "", "name of post")
}
