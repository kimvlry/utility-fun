# GoSort - Simple Sort Utility

A simplified UNIX-style sort implementation in Go with support for common sorting options.

## Features

### Basic Sorting Options
- `-k N` - sort by column N (tab-delimited by default)
- `-n` - numeric sort
- `-r` - reverse sort order
- `-u` - output only unique lines

### ...more options
- `-M` - sort by month names (Jan, Feb,...Dec)
- `-b` - ignore trailing blanks
- `-c` - check if input is sorted (if not - notify)
- `-h` - human-readable numeric sort (K - kilobyte, M - megabyte suffixes)

## Implementation Details

- Processes input from files or STDIN
- Properly combines multiple flags (e.g., `-nr` for reverse numeric sort)
- Passes `go vet` and `golint` checks
- Comprehensive unit test coverage

## Usage Examples
```bash
