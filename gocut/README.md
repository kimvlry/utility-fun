# GoCut - Column Extraction Tool

A Go implementation of UNIX `cut` functionality with support for field selection by column.

## Supported flags

### Core Functionality
- `-f "fields"` — select fields/ranges (e.g., `1,3-5`)
- `-d "delimiter"` — specify custom delimiter (default: tab `\t`)
- `-s` — suppress lines without delimiters

### Field Specification
- Single fields (e.g., `3`)
- Ranges (e.g., `2-4`)
- Mixed combinations (e.g., `1,3-5,7`)
- Out-of-bound fields are ignored gracefully

## Usage
### build
```bash
  go install ./cmd/gocut
```

### examples
```
gocut -f 2,4 -d "," data.csv
```

```
echo "a:b:c" | gocut -f 2 -d ":"
```

```
cat data.csv | gocut -f 1,3 -d ","
```