package downloader

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

// getCtx performs HTTP GET with context
func (d *Downloader) getCtx(ctx context.Context, u *url.URL) ([]byte, string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, "", err
	}
	resp, err := d.client.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatalf("failed to close response body: %v", err)
		}
	}(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, "", fmt.Errorf("HTTP %d", resp.StatusCode)
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}
	return b, resp.Header.Get("Content-Type"), nil
}

// fetchAssetCtx downloads an asset
func (d *Downloader) fetchAssetCtx(ctx context.Context, u *url.URL) error {
	b, _, err := d.getCtx(ctx, u)
	if err != nil {
		return err
	}
	p := urlToLocalPath(d.rootDir, u)
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		return err
	}
	return os.WriteFile(p, b, 0o644)
}
