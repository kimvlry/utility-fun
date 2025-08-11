package cmd

import (
	"os"
	"wget/internal/downloader"

	"github.com/spf13/cobra"
)

var depth int64

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "wget",
	Short: "Create a local copy of a website (or a subset) that can be browsed offline.",
	Long: `A simplified wget -m-like tool for recursive website downloading with mirroring capabilities. 
Creates a local copy of a website (or a subset) that can be browsed offline.

Arguments: 
    URL - address of a site to download

Example:
    https://100go.co --depth=3`,

	Args: cobra.ExactArgs(1),

	RunE: func(cmd *cobra.Command, args []string) error {
		url := args[0]
		return downloader.DownloadPage(url, depth)
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
	rootCmd.Flags().Int64Var(&depth, "depth", 5, "Links recursion depth control. Default: 5")
}
