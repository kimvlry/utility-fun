package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"grep/internal/grepper"
)

var cfg grepper.Config

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gogrep [flags] pattern [file]",
	Short: "a minimal unix grep-like text filter for files or stdin",
	Long: `A minimal unix grep-like tool for filtering text streams.

It reads input from a file or stdin and prints lines that match a given pattern
(substring or regular expression). Behavior is intentionally close to UNIX grep.

If [file] is omitted or set to "-" the input is read from stdin.
    
Examples:
  gogrep -F -i "todo" main.go

  echo -e "alpha\nBeta\nGAMMA" | gogrep -A 1 "alpha"

  gogrep -n -C 2 'error|warn' /var/log/syslog

  gogrep -c "pattern" file.txt

  cat data.txt | gogrep -n "foo" -
`,

	Args: cobra.RangeArgs(1, 2),

	RunE: runGrep,
}

func runGrep(cmd *cobra.Command, args []string) error {
	pattern := args[0]

	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	source, err := openInputSource(args)
	if err != nil {
		return err
	}

	_, err = grepper.GrepLines(source, cmd.OutOrStdout(), pattern, cfg)
	if err != nil {
		return fmt.Errorf("grep failed: %w", err)
	}
	return nil
}

// openInputSource returns an io.Reader for the input source (file or stdin).
func openInputSource(args []string) (io.Reader, error) {
	if len(args) == 2 && args[1] != "-" {
		f, err := os.Open(args[1])
		if err != nil {
			return nil, fmt.Errorf("failed to open file %s: %w", args[1], err)
		}
		return &fileReader{f}, nil
	}
	return os.Stdin, nil
}

// fileReader wraps *os.File to close automatically after reading is done.
// Implements io.Reader.
type fileReader struct {
	f *os.File
}

func (fr *fileReader) Read(p []byte) (int, error) {
	n, err := fr.f.Read(p)
	if err == io.EOF {
		closeErr := fr.f.Close()
		if closeErr != nil {
			_, err := fmt.Fprintf(os.Stderr, "failed to close file: %v\n", closeErr)
			if err != nil {
				return 0, err
			}
		}
	}
	return n, err
}

// init registers all command-line flags for the root command.
func init() {
	rootCmd.Flags().IntVarP(&cfg.After, "after", "A", 0, "print N lines of trailing context after each match")
	rootCmd.Flags().IntVarP(&cfg.Before, "before", "B", 0, "print N lines of leading context before each match")
	rootCmd.Flags().IntVarP(&cfg.Context, "context", "C", 0, "print N lines of context around each match (equivalent to -A N -B N)")

	rootCmd.Flags().BoolVarP(&cfg.CountOnly, "count", "c", false, "print only a count of matching lines")
	rootCmd.Flags().BoolVarP(&cfg.IgnoreCase, "ignore-case", "i", false, "ignore case distinctions")
	rootCmd.Flags().BoolVarP(&cfg.Invert, "invert-match", "v", false, "select non-matching lines")
	rootCmd.Flags().BoolVarP(&cfg.Fixed, "fixed-strings", "F", false, "interpret pattern as a fixed substring (not a regular expression)")
	rootCmd.Flags().BoolVarP(&cfg.WithLineNo, "line-number", "n", false, "print line number with each output line")
}
