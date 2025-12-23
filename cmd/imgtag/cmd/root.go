package cmd

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/evanoberholster/imagemeta/xmp"
	"github.com/goccy/go-yaml"
	"github.com/ohzqq/imgtag"
	"github.com/spf13/cobra"
)

var (
	YAML        bool
	JSON        bool
	EXT         string
	outputFile  string
	batchOutput string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "imgtag",
	Short: "get some image meta",
	Long:  ``,
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

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolP("creator", "c", false, "get creator")
	rootCmd.PersistentFlags().BoolP("title", "t", false, "get title")
	rootCmd.PersistentFlags().BoolP("description", "d", false, "get description")
	rootCmd.PersistentFlags().BoolP("subject", "s", false, "get subject")
	rootCmd.PersistentFlags().StringVarP(&batchOutput, "batch", "b", "", "output the metadata as a batch/slice/array")
	rootCmd.PersistentFlags().BoolVarP(&YAML, "yaml", "y", true, "marshal to yaml")
	rootCmd.PersistentFlags().BoolVarP(&JSON, "json", "j", false, "marshal to json")
	rootCmd.PersistentFlags().StringP("write", "w", ".yaml", "output meta to file")
	rootCmd.PersistentFlags().StringVarP(&EXT, "ext", "e", "", "extension for meta files")
	rootCmd.PersistentFlags().StringVarP(&outputFile, "output", "o", "images.yaml", "file output name")
}

func writeMeta(args []string) error {
	metas, err := decodeManyMuchMeta(args)
	if err != nil {
		return err
	}
	w := os.Stdout
	for i, meta := range metas {
		dir, name := filepath.Split(args[i])
		name = strings.TrimSuffix(name, filepath.Ext(args[i]))
		ext := ".yaml"
		if JSON {
			ext = ".json"
		}
		if EXT != "" {
			ext = EXT
		}
		w, err = os.Create(filepath.Join(dir, name) + ext)
		if err != nil {
			return err
		}
		defer w.Close()
		err = encodeMeta(w, meta.DublinCore())
		if err != nil {
			return err
		}
	}
	return nil
}

func writeMetaBatch(args []string) error {
	all, err := metaSlice(args)
	if err != nil {
		return err
	}
	ext := filepath.Ext(batchOutput)
	switch ext {
	case "":
		batchOutput = batchOutput + getExt()
	case ".yaml", ".yml":
		YAML = true
	case ".json":
		JSON = true
	}
	out, err := os.Create(batchOutput)
	if err != nil {
		return err
	}
	return encodeMeta(out, all)
}

func encodeMeta(w io.Writer, meta any) error {
	if JSON {
		err := encodeJSON(w, meta)
		if err != nil {
			return err
		}
	} else if YAML {
		err := encodeYAML(w, meta)
		if err != nil {
			return err
		}
		w.Write([]byte("---"))
	}
	return nil
}

func encodeJSON(w io.Writer, meta any) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(meta)
}

func encodeYAML(w io.Writer, meta any) error {
	enc := yaml.NewEncoder(w, yaml.Indent(2), yaml.AutoInt(), yaml.OmitEmpty())
	defer enc.Close()
	return enc.Encode(meta)
}

func decodeManyMuchMeta(args []string) ([]*imgtag.Img, error) {
	imgs := make([]*imgtag.Img, len(args))
	for i, arg := range args {
		im, err := decodeMeta(arg)
		if err != nil {
			return nil, err
		}
		imgs[i] = im
	}
	return imgs, nil
}

func decodeMeta(name string) (*imgtag.Img, error) {
	i, err := imgtag.NewImg(name)
	if err != nil {
		return nil, err
	}
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	err = i.DecodeMeta(f)
	if err != nil {
		return nil, err
	}
	return i, nil
}

func getExt() string {
	ext := ".yaml"
	if JSON {
		ext = ".json"
	}
	if EXT != "" {
		ext = EXT
	}
	return ext
}

func metaSlice(args []string) ([]xmp.DublinCore, error) {
	metas, err := decodeManyMuchMeta(args)
	if err != nil {
		return nil, err
	}
	all := make([]xmp.DublinCore, len(args))
	for i, meta := range metas {
		all[i] = meta.DublinCore()
		if len(all[i].Title) == 0 {
			all[i].Title = []string{
				strings.TrimSuffix(filepath.Base(all[i].Identifier), filepath.Ext(all[i].Identifier)),
			}
		}
	}
	return all, nil
}

func saveMeta(name string, m any) error {
	w, err := os.Create(name)
	if err != nil {
		return err
	}
	switch {
	case YAML:
		enc := yaml.NewEncoder(w, yaml.Indent(2), yaml.AutoInt(), yaml.OmitEmpty())
		err = enc.Encode(m)
		if err != nil {
			return err
		}
		w.Write([]byte("---"))
	case JSON:
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		err = enc.Encode(m)
		if err != nil {
			return err
		}
	}
	return nil
}
