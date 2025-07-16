# GoTelnet - Simple Telnet Client

A simple telnet-like utility for establishing TCP connections and interacting with remote servers.

## Features

- Accepts the following command-line parameters:
    - `host` (string): Target hostname or IP address (required)
    - `port` (int): Target TCP port (required)
    - `--timeout` (duration): Optional connection timeout (e.g. `--timeout=5s`, `--timeout=1m`). Default is 10 seconds.

- Establishes a TCP connection to the specified `host:port` upon launch.
- Sends all user input from STDIN to the socket.
- Prints all data received from the socket to STDOUT.
- Gracefully shuts down when:
    - The user presses `Ctrl+D` (EOF)
    - The server closes the connection
- If the connection attempt fails (e.g., the server is unavailable), the program terminates after the specified 
timeout with an appropriate error message.

## Implementation Notes

- Input/output are handled in separate goroutines for concurrent communication.
- Use the `net` package (`net.Conn`) for TCP communication.

## Testing
- Comprehensive unit test coverage
- Tested locally by connecting to an SMTP server and manually sending commands.

## Example Usage

```bash

