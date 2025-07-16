# GoCut - Column Extraction Tool

A Go implementation of UNIX cut functionality with support for field selection by column.

## Features

### Core Functionality
- `-f "fields"` - select fields/ranges (e.g., "1,3-5")
- `-d "delimiter"` - specify custom delimiter (default: tab)
- `-s` - suppress lines without delimiters

### Advanced Capabilities
- Handles complex field specifications:
    - Single fields (e.g., "3")
    - Ranges (e.g., "2-4")
    - Mixed combinations (e.g., "1,3-5,7")
- Gracefully handles out-of-bound columns
- Efficient large file processing

## Implementation Details

- Processes input from STDIN
- Optimized for performance with:
    - ,
- Robust error handling:
    - Invalid field specifications
    - Malformed input
- Passes all Go quality checks:
    - `gofmt`
    - `golint`
    - `go vet`
- Comprehensive unit test coverage

## Usage Examples

Basic column extraction:
```bash
