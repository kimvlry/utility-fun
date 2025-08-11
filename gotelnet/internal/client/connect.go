package client

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

// Connect establishes a TCP connection to host:port with the given timeout.
// After connecting, it concurrently copies stdin to the socket and socket to stdout.
// The function exits gracefully when Ctrl+D (EOF) is pressed or the server closes the connection.
func Connect(host string, port string, timeout time.Duration) error {
	addr := net.JoinHostPort(host, port)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	dialer := &net.Dialer{}
	// DialContext connects to the address on the named network using the provided context.
	conn, err := dialer.DialContext(ctx, "tcp", addr)

	if err != nil {
		return fmt.Errorf("could not connect to server: %w", err)
	}
	defer conn.Close()

	fmt.Printf("Connected to server: %s. Press Ctrl+D to exit.\n", addr)

	go func() {
		if err := readFromConn(ctx, conn); err != nil {
			fmt.Printf("read error: %v\n", err)
			cancel()
		}
	}()

	go func() {
		if err := writeToConn(ctx, conn); err != nil {
			fmt.Printf("write error: %v\n", err)
			cancel()
		}
	}()

	<-ctx.Done()
	fmt.Println("Connection closed.")
	return nil
}

// readFromConn reads data from connected server and writes to stdout.
// Returns error on failure or nil on normal close.
func readFromConn(ctx context.Context, conn net.Conn) error {
	buf := make([]byte, 4096)

	for {
		select {
		case <-ctx.Done():

		default:
			conn.SetReadDeadline(time.Now().Add(500 * time.Millisecond))

			n, err := conn.Read(buf)
			if err != nil {
				var ne net.Error
				if errors.As(err, &ne) && ne.Timeout() {
					continue
				}
				if err == io.EOF {
					return nil
				}
				return err
			}

			if n > 0 {
				os.Stdout.Write(buf[:n])
			}
		}
	}
}

// writeToConn reads lines from stdin and writes to the connection.
// Returns error on failure or nil on normal close (Ctrl+D).
func writeToConn(ctx context.Context, conn net.Conn) error {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			if !scanner.Scan() {
				return nil
			}
			line := scanner.Bytes()
			_, err := conn.Write(append(line, '\n'))
			if err != nil {
				return err
			}
		}
	}
}
