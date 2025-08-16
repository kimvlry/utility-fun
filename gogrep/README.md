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

## Usage 
### build
```bash 
  go install ./cmd/gogrep  
```
### examples
```bash
  echo "alpha\nBeta\nGAMMA" | gogrep -A 1 "alpha" 
```

```bash
  gogrep -c "beta" test.txt  
```

