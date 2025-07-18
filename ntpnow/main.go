package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/beevik/ntp"
)

var (
	server  string
	timeout time.Duration
)

func init() {
	flag.StringVar(&server, "server", "pool.ntp.org", "NTP server to use (default: pool.ntp.org)")
	flag.DurationVar(&timeout, "timeout", 5*time.Second, "Timeout for NTP request (default: 5s)")
}

func NtpNow() time.Time {
	flag.Parse()

	opts := ntp.QueryOptions{Timeout: timeout}
	resp, err := ntp.QueryWithOptions(server, opts)

	if err != nil {
		var netErr net.Error
		if errors.As(err, &netErr) && netErr.Timeout() {
			log.Fatalf("Error: timeout after %s while connecting to %q", timeout, server)
		}

		log.Fatalf("Error querying %q: %v", server, err)
	}

	return resp.Time
}

func main() {
	t := NtpNow()

	fmt.Printf("UTC:   %s\n", t.Format(time.RFC3339Nano))
	fmt.Printf("Local: %s\n", t.Local().Format(time.RFC3339Nano))
}
