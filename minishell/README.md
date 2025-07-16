# Minishell - Simple Unix Shell

A minimal UNIX shell implementation in Go.

## Features

### Built-in Commands
- `cd <path>` - change directory
- `pwd` - print working directory
- `echo <args>` - display arguments
- `kill <pid>` - send termination signal
- `ps` - list running processes

### Core Functionality
- External command execution via `exec`
- Pipeline support (`cmd1 | cmd2 | cmd3`)
- Signal handling:
    - Ctrl+D - shell termination
    - Ctrl+C - command interruption (without shell exit)
- Conditional execution (`&&`, `||`)
- Environment variable substitution (`$VAR`)
- I/O redirection (`>`, `<`, `>>`)

## Implementation Details
- Passes standard Go checks (`gofmt`, `golint`, `go vet`)
- Comprehensive unit test coverage

## Usage Examples

Basic commands:
```sh
