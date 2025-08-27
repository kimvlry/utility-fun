package downloader

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sync"
	"time"

	"log/slog"
)

// Config holds downloader settings
type Config struct {
	StartURL *url.URL
	MaxDepth int
	OutDir   string
	Timeout  time.Duration
	Parallel int
}

// Downloader is the main engine
type Downloader struct {
	cfg      Config
	client   *http.Client
	rootHost string
	rootDir  string

	mu        sync.Mutex
	seenPage  map[string]struct{}
	seenAsset map[string]struct{}
}

// New creates a new Downloader
func New(cfg Config) (*Downloader, error) {
	if cfg.StartURL == nil {
		return nil, errors.New("nil StartURL")
	}
	if cfg.OutDir == "" {
		cfg.OutDir = "."
	}
	if cfg.Timeout <= 0 {
		cfg.Timeout = 20 * time.Second
	}
	if cfg.Parallel <= 0 {
		cfg.Parallel = 4
	}

	d := &Downloader{
		cfg:       cfg,
		client:    &http.Client{Timeout: cfg.Timeout},
		rootHost:  cfg.StartURL.Host,
		rootDir:   filepath.Join(cfg.OutDir, cfg.StartURL.Host),
		seenPage:  make(map[string]struct{}),
		seenAsset: make(map[string]struct{}),
	}

	if err := os.MkdirAll(d.rootDir, 0o755); err != nil {
		return nil, err
	}
	return d, nil
}

// RootDir returns local root directory
func (d *Downloader) RootDir() string { return d.rootDir }

// job represents a download task
type job struct {
	u     *url.URL
	depth int
	kind  string // "page" or "asset"
}

// Run starts downloading with context
func (d *Downloader) Run() error {
	ctx, cancel := context.WithTimeout(context.Background(), d.cfg.Timeout)
	defer cancel()

	jobs := make(chan job, 256)
	var wg sync.WaitGroup
	var active sync.WaitGroup

	enqueue := func(j job) {
		active.Add(1)
		select {
		case jobs <- j:
		case <-ctx.Done():
			slog.Warn("enqueue cancelled", "url", j.u.String())
			active.Done()
		}
	}

	for i := 0; i < d.cfg.Parallel; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case j, ok := <-jobs:
					if !ok {
						return
					}
					switch j.kind {
					case "page":
						d.processPageCtx(ctx, j, enqueue)
					case "asset":
						if err := d.fetchAssetCtx(ctx, j.u); err != nil {
							slog.Warn("asset fetch failed", "url", j.u.String(), "err", err)
						}
					}
					active.Done()
				case <-ctx.Done():
					return
				}
			}
		}()
	}

	enqueue(job{u: d.cfg.StartURL, depth: 0, kind: "page"})

	go func() {
		active.Wait()
		close(jobs)
	}()

	wg.Wait()
	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		return fmt.Errorf("download timed out after %s", d.cfg.Timeout)
	}
	return nil
}
