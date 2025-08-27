package cmd

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"time"
	"wget/internal/downloader"

	"github.com/spf13/cobra"
)

var (
	flagDepth      int
	flagOutDir     string
	flagTimeoutSec int
	flagParallel   int
)

var rootCmd = &cobra.Command{
	Use:   "gowget <url>",
	Short: "Mirror a website (HTML + assets) for offline use",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		raw := args[0]
		u, err := url.Parse(raw)
		if err != nil || u.Scheme == "" || u.Host == "" {
			return fmt.Errorf("invalid URL: %s", raw)
		}

		cfg := downloader.Config{
			StartURL: u,
			MaxDepth: flagDepth,
			OutDir:   flagOutDir,
			Timeout:  time.Duration(flagTimeoutSec) * time.Second,
			Parallel: flagParallel,
		}

		dl, err := downloader.New(cfg)
		if err != nil {
			return err
		}
		if err := dl.Run(); err != nil {
			return err
		}
		log.Printf("Done. Saved to %s\n", dl.RootDir())
		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		_, err := fmt.Fprintln(os.Stderr, err)
		if err != nil {
			log.Fatalf("failed to execute root command: %v", err)
		}
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().IntVarP(&flagDepth, "depth", "d", 1, "recursion depth for <a> links (0 = only the start page)")
	rootCmd.Flags().StringVarP(&flagOutDir, "outdir", "o", ".", "output directory (root for mirrored site)")
	rootCmd.Flags().IntVar(&flagTimeoutSec, "timeout", 20, "HTTP timeout in seconds")
	rootCmd.Flags().IntVarP(&flagParallel, "parallel", "c", 4, "number of parallel downloads")
}
