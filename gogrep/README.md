# GoGrep - Text Filtering Utility

A Go implementation of UNIX grep functionality with support for common filtering options.

## Supported flags

### Matching Options
- Basic substring matching
- Regular expression support
- `-F` - fixed string matching (no regex)
- `-i` - case-insensitive matching
- `-v` - invert match (show non-matching lines)

### Output Control
- `-A N` - show N lines after match
- `-B N` - show N lines before match
- `-C N` - show N lines of context around match
- `-c` - count matching lines only
- `-n` - show line numbers

## Implementation Details

- Processes input from STDIN or files
- Properly combines multiple flags (example: `-C 2 -n -i`)
- Edge case handling:
    - File boundaries for context
    - Overlapping matches
    - Repeated matches
- Passes `gofmt` and `golint` checks
- Comprehensive unit test coverage

## Usage Examples

Basic matching:
```bash
