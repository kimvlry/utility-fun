package cmd

import (
	"gotelnet/internal/client"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var timeout time.Duration

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gotelnet <host> <port>",
	Short: `A simple telnet utility for TCP connections`,
	Long: `gotelnet connects to a TCP server and lets you interact with it.

Arguments:
  host  - Target hostname or IP address to connect to (e.g. example.com, 192.168.0.10)
  port  - TCP port number on the target host (e.g. 23 for Telnet, 25 for SMTP)

Example:
  gotelnet example.com 23 --timeout=5s
  gotelnet 192.168.0.10 80`,

	Args: cobra.ExactArgs(2),

	RunE: func(cmd *cobra.Command, args []string) error {
		host := args[0]
		port := args[1]
		return client.Connect(host, port, timeout)
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
	rootCmd.Flags().DurationVar(&timeout, "timeout", time.Second*10, "connection timeout (e.g. 5s, 1m)")
}
