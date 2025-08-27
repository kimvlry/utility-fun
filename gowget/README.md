# GoWget - Web Mirroring Utility

A simplified `wget -m`-like tool for recursive website downloading with mirroring capabilities.
Creates a local copy of a website (or a subset) that can be browsed offline.

### Options
- Links recursion depth control (`--depth=N`)
- Timeout configuration (`--timeout=N`)
- Parallel downloads (`--parallel=N`) - 

TODO:
- robots.txt support

## Usage Examples
1.build
```bash
  go install ./cmd/gowget
```

2. examples
```
gowget https://github.com/ --depth=0 --parallel=10 --outdir=./mirror --timeout=10
```
