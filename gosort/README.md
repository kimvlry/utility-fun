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


## Usage
### 1. build
```bash
    go install ./cmd/gosort
```

### 2. examples
```bash
    echo -e "3 apple\n2\n1 banana\n" | gosort 
```

```bash
    gosort -k2       
    3 apple
    2
    1 banana
    
    CtrlD
```

```bash
    gosort cmd/example/example.txt -k2 -nr
```