package downloader

import (
	"crypto/sha1"
	"encoding/hex"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// urlToLocalPath converts URL to a local filesystem path
func urlToLocalPath(rootDir string, u *url.URL) string {
	p := u.Path
	if p == "" || strings.HasSuffix(p, "/") {
		p = strings.TrimSuffix(p, "/") + "/index.html"
	}
	base := filepath.Base(p)
	if !strings.Contains(base, ".") {
		p = strings.TrimSuffix(p, "/") + "/index.html"
	}
	if u.RawQuery != "" {
		h := sha1.Sum([]byte(u.String()))
		h8 := hex.EncodeToString(h[:])[:8]
		base = filepath.Base(p)
		dir := filepath.Dir(p)
		dot := strings.LastIndex(base, ".")
		if dot < 0 {
			base = base + "-" + h8
		} else {
			base = base[:dot] + "-" + h8 + base[dot:]
		}
		p = filepath.Join(dir, base)
	}
	p = strings.TrimPrefix(p, "/")
	return filepath.Join(rootDir, p)
}

// rewriteHTML rewrites links and collects page links and assets
func rewriteHTML(base *url.URL, html []byte, mapFn func(*url.URL) string) ([]byte, []*url.URL, []*url.URL, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(html)))
	if err != nil {
		return nil, nil, nil, err
	}

	resolve := func(s string) *url.URL {
		if s == "" {
			return nil
		}
		u, err := url.Parse(s)
		if err != nil {
			return nil
		}
		u = base.ResolveReference(u)
		u.Fragment = ""
		return u
	}

	var links []*url.URL
	var assets []*url.URL

	doc.Find("a[href]").Each(func(_ int, s *goquery.Selection) {
		v, _ := s.Attr("href")
		if u := resolve(v); u != nil {
			links = append(links, u)
			if m := mapFn(u); m != "" {
				_ = s.SetAttr("href", m)
			}
		}
	})

	doc.Find("img[src],script[src]").Each(func(_ int, s *goquery.Selection) {
		v, _ := s.Attr("src")
		if u := resolve(v); u != nil {
			assets = append(assets, u)
			if m := mapFn(u); m != "" {
				_ = s.SetAttr("src", m)
			}
		}
	})

	doc.Find("link[href]").Each(func(_ int, s *goquery.Selection) {
		rel, _ := s.Attr("rel")
		v, _ := s.Attr("href")
		if u := resolve(v); u != nil {
			lower := strings.ToLower(rel)
			if strings.Contains(lower, "stylesheet") || strings.Contains(lower, "icon") || strings.Contains(lower, "preload") {
				assets = append(assets, u)
			} else {
				links = append(links, u)
			}
			if m := mapFn(u); m != "" {
				_ = s.SetAttr("href", m)
			}
		}
	})

	out, err := doc.Html()
	if err != nil {
		return nil, nil, nil, err
	}
	return []byte(out), links, assets, nil
}

// isSameHost checks if URL is in the root domain
func (d *Downloader) isSameHost(u *url.URL) bool {
	return strings.EqualFold(u.Host, d.rootHost)
}
