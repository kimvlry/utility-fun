package cmd

import (
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	"gosort/internal/sorter"
	"log"
	"os"
	"strings"
)

var cfg sorter.Config

func init() {
	RootCmd.PersistentFlags().BoolP("help", "", false, "disable default help")

	RootCmd.Flags().IntVarP(&cfg.Key, "key", "k", 0, "column number to sort by (starting from 1)")
	RootCmd.Flags().BoolVarP(&cfg.Numeric, "numeric", "n", false, "sort numerically")
	RootCmd.Flags().BoolVarP(&cfg.Reverse, "reverse", "r", false, "reverse sort order")
	RootCmd.Flags().BoolVarP(&cfg.Unique, "unique", "u", false, "output unique lines only")
	RootCmd.Flags().BoolVarP(&cfg.Month, "month", "M", false, "sort by month name")
	RootCmd.Flags().BoolVarP(&cfg.IgnoreTrailing, "ignore-trailing-blanks", "b", false, "ignore trailing blanks")
	RootCmd.Flags().BoolVarP(&cfg.CheckIfSorted, "check-if-sorted", "c", false, "checkIfSorted if input is sorted")
	RootCmd.Flags().BoolVarP(&cfg.HumanNum, "human-readable-numeric", "h", false, "human readable numeric sort")
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		_, err := fmt.Fprintln(os.Stderr, err)
		if err != nil {
			log.Fatalf("failed to execute: %v", err)
		}
		os.Exit(1)
	}
}

var RootCmd = &cobra.Command{
	Use:   "sort",
	Short: "Sort lines of text",
	Long: `Sort lines of text from a file or standard input.
If no filename is provided, input is read from stdin.

Examples:
  gosort file.txt
  gosort -n -k2 data.tsv
  cat file.txt | gosort -r

Args: 
    filename - name of file to read lines from`,
	Args: cobra.RangeArgs(0, 1),

	PreRun: func(cmd *cobra.Command, args []string) {
		for i, arg := range os.Args[1:] {
			if strings.HasPrefix(arg, "-k") && len(arg) > 2 {
				os.Args = append(os.Args[:i+1], append([]string{"-k", arg[2:]}, os.Args[i+2:]...)...)
			}
		}
	},

	RunE: func(cmd *cobra.Command, args []string) error {
		if err := cfg.Validate(); err != nil {
			return err
		}

		var scanner *bufio.Scanner
		if len(args) > 0 {
			file := args[0]
			f, err := os.Open(file)
			if err != nil {
				return err
			}
			defer func(f *os.File) {
				err := f.Close()
				if err != nil {
					log.Fatalf("failed to close file: %v", err)
				}
			}(f)
			scanner = bufio.NewScanner(f)
		} else {
			scanner = bufio.NewScanner(os.Stdin)
		}

		var lines []string
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			return err
		}

		sortedLines, err := sorter.SortLines(lines, cfg)
		if err != nil {
			return err
		}

		out := cmd.OutOrStdout()
		for _, line := range sortedLines {
			_, err := fmt.Fprintln(out, line)
			if err != nil {
				return err
			}
		}

		return nil
	},
}
