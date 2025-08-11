# GoTelnet - Simple Telnet Client

A simple telnet utility for establishing TCP connections and interacting with remote servers,
implemented in go

## Usage
- Accepts the following command-line parameters:
  - `host` (string): Target hostname or IP address (required)
  - `port` (int): Target TCP port (required)
  - `--timeout` (duration): Optional connection timeout (e.g. `--timeout=5s`, `--timeout=1m`). Default is 10 seconds.
  ```text
    gotelnet <host> <port> [--timeout=10s]
  ```

  ```bash 
   gotelnet git:(main) ✗ gotelnet tcpbin.com 4242 --timeout=3s
  ```

## Details
- Establishes a TCP connection to the specified `host:port` upon launch.
- Sends all user input from STDIN to the socket.
- Prints all data received from the socket to STDOUT.
- Gracefully shuts down when:
    - The user presses `Ctrl+D` (EOF)
    - The server closes the connection
- If the connection attempt fails (e.g., the server is unavailable), the program terminates after the specified
timeout with an appropriate error message.


- Input/output are handled in separate goroutines for concurrent communication.
- `net` package (`net.Conn`) used for TCP communication.

## Testing
- Tested locally by connecting to an SMTP server and manually sending commands.

