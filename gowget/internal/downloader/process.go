package downloader

import (
	"context"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"log/slog"
)

// processPageCtx downloads a page, rewrites links, and enqueues next jobs
func (d *Downloader) processPageCtx(ctx context.Context, j job, enqueue func(job)) {
	select {
	case <-ctx.Done():
		return
	default:
	}

	abs := j.u.String()
	d.mu.Lock()
	if _, ok := d.seenPage[abs]; ok {
		d.mu.Unlock()
		return
	}
	d.seenPage[abs] = struct{}{}
	d.mu.Unlock()

	slog.Info("page", "url", abs, "depth", j.depth)

	body, ctype, err := d.getCtx(ctx, j.u)
	if err != nil {
		slog.Warn("get page failed", "url", abs, "err", err)
		return
	}

	pageLocal := urlToLocalPath(d.rootDir, j.u)
	if err := os.MkdirAll(filepath.Dir(pageLocal), 0o755); err != nil {
		slog.Warn("mkdir failed", "path", pageLocal, "err", err)
		return
	}

	if !strings.HasPrefix(ctype, "text/html") {
		_ = os.WriteFile(pageLocal, body, 0o644)
		return
	}

	rewritten, links, assets, err := rewriteHTML(j.u, body, func(ref *url.URL) string {
		if !d.isSameHost(ref) {
			return ref.String()
		}
		local := urlToLocalPath(d.rootDir, ref)
		rel, err := filepath.Rel(filepath.Dir(pageLocal), local)
		if err != nil {
			return ref.String()
		}
		return filepath.ToSlash(rel)
	})
	if err != nil {
		rewritten = body
	}
	_ = os.WriteFile(pageLocal, rewritten, 0o644)

	// enqueue new pages
	if j.depth < d.cfg.MaxDepth {
		for _, l := range links {
			if d.isSameHost(l) {
				d.mu.Lock()
				if _, ok := d.seenPage[l.String()]; !ok {
					enqueue(job{u: l, depth: j.depth + 1, kind: "page"})
				}
				d.mu.Unlock()
			}
		}
	}

	// enqueue assets
	for _, a := range assets {
		if !d.isSameHost(a) {
			continue
		}
		d.mu.Lock()
		if _, ok := d.seenAsset[a.String()]; ok {
			d.mu.Unlock()
			continue
		}
		d.seenAsset[a.String()] = struct{}{}
		d.mu.Unlock()
		enqueue(job{u: a, kind: "asset"})
	}
}
