# GoWget - Web Mirroring Utility

A simplified `wget -m`-like tool for recursive website downloading with mirroring capabilities.
Creates a local copy of a website (or a subset) that can be browsed offline.

## Features

### Core Functionality
- Recursive downloading with configurable links depth
- Supports all web resources:
    - HTML pages
    - CSS stylesheets
    - JavaScript files
    - Images and media
- Domain-bound crawling (stays within original domain)
- URL normalization and deduplication

### Advanced Options
- Links recursion depth control (`--depth=N`)
- Parallel downloads (`--workers=N`)
- Timeout configuration (`--timeout=N`)
- `-robots`: Honor robots.txt directives (optional)
## Implementation 

### Architecture Components
1. 

### Details
- Proper error handling (network/filesystem)
- Link cycling detection
- Memory efficiency
- Code passes `golint`, `go vet`
- Comprehensive unit test coverage

## Usage Examples

Mirroring with set depth, timeout and num of parallel download workers:
```bash
