package cmd

import (
	"bufio"
	"cut/internal/cutter"
	"cut/internal/parser"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var (
	fieldSpec string
	delimiter string
	separated bool
)

func openInputSource(args []string) (io.Reader, error) {
	if len(args) > 0 && args[0] != "-" {
		f, err := os.Open(args[0])
		if err != nil {
			return nil, fmt.Errorf("failed to open file %s: %w", args[0], err)
		}
		return f, nil
	}
	return os.Stdin, nil
}

var rootCmd = &cobra.Command{
	Use:   "gocut",
	Short: "",
	RunE: func(cmd *cobra.Command, args []string) error {
		fields, err := parser.ParseFieldSpec(fieldSpec)
		if err != nil {
			return fmt.Errorf("invalid format of -f value: %w", err)
		}

		extractor := cutter.NewExtractor(delimiter, fields, separated)

		source, err := openInputSource(args)
		if err != nil {
			return err
		}

		scanner := bufio.NewScanner(source)
		for scanner.Scan() {
			line := scanner.Text()
			if out := extractor.Extract(line); out != "" {
				fmt.Println(out)
			}
		}
		if err := scanner.Err(); err != nil {
			return fmt.Errorf("error reading: %w", err)
		}
		return nil
	},
}

func init() {
	rootCmd.Flags().StringVarP(&fieldSpec, "fields", "f", "", "fields number (example, \"1,3-5\") (required)")
	err := rootCmd.MarkFlagRequired("fields")
	if err != nil {
		log.Fatalf("could not init cobra flags: %v", err)
	}

	rootCmd.Flags().StringVarP(&delimiter, "delimiter", "d", "\t", "field delimiter (default: tab)")
	rootCmd.Flags().BoolVarP(&separated, "separated", "s", false, "only show lines that include delimiter")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
