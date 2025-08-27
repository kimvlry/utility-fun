package downloader

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func TestURLToLocalPath(t *testing.T) {
	root := t.TempDir()
	parse := func(s string) *url.URL {
		u, _ := url.Parse(s)
		return u
	}

	tests := map[string]string{
		"https://ex.com/":                       "index.html",
		"https://ex.com/about":                  "about/index.html",
		"https://ex.com/about/":                 "about/index.html",
		"https://ex.com/assets/app.js":          "assets/app.js",
		"https://ex.com/p?q=1":                  "p/index-",
		"https://ex.com/dir/file":               "dir/file/index.html",
		"https://ex.com/dir/file.html":          "dir/file.html",
		"https://ex.com/dir/sub/?a=b":           "dir/sub/index-",
		"https://ex.com/dir/sub/index.html?x=y": "dir/sub/index-",
		"https://ex.com/dir/sub/index.html":     "dir/sub/index.html",
		"https://ex.com/with.dot/":              "with.dot/index.html",
	}

	for in, wantPrefix := range tests {
		got := urlToLocalPath(root, parse(in))
		if !strings.HasPrefix(got, filepath.Join(root, wantPrefix)) {
			t.Fatalf("map %q => %q; want prefix %q", in, got, filepath.Join(root, wantPrefix))
		}
	}
}

func TestMirrorBasicAndRewrite(t *testing.T) {
	// Mock site:
	//   / -> index.html that links to /page and references /static/app.js and /img/logo.png
	//   /page -> simple html
	var concurrent int32
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&concurrent, 1)
		defer atomic.AddInt32(&concurrent, -1)

		switch r.URL.Path {
		case "/":
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			io := `<!doctype html>
<html>
<head>
  <link rel="stylesheet" href="/static/style.css">
  <script src="/static/app.js"></script>
</head>
<body>
  <img src="/img/logo.png">
  <a href="/page">Go page</a>
</body>
</html>`
			_, _ = w.Write([]byte(io))
		case "/page":
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			_, _ = w.Write([]byte("<html><body>PAGE</body></html>"))
		case "/static/app.js":
			time.Sleep(50 * time.Millisecond) // simulate latency
			w.Header().Set("Content-Type", "application/javascript")
			_, _ = w.Write([]byte("console.log('ok')"))
		case "/static/style.css":
			time.Sleep(50 * time.Millisecond)
			w.Header().Set("Content-Type", "text/css")
			_, _ = w.Write([]byte("body{background:#fff}"))
		case "/img/logo.png":
			w.Header().Set("Content-Type", "image/png")
			_, _ = w.Write([]byte{0x89, 0x50})
		default:
			http.NotFound(w, r)
		}
	}))
	defer ts.Close()

	u, _ := url.Parse(ts.URL + "/")
	out := t.TempDir()

	dl, err := New(Config{
		StartURL: u,
		MaxDepth: 1,
		OutDir:   out,
		Timeout:  5 * time.Second,
		Parallel: 4,
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := dl.Run(); err != nil {
		t.Fatal(err)
	}

	// Check files exist
	rootDir := filepath.Join(out, u.Host)
	want := []string{
		"index.html",
		"page/index.html",
		"static/app.js",
		"static/style.css",
		"img/logo.png",
	}
	for _, p := range want {
		if _, err := os.Stat(filepath.Join(rootDir, filepath.FromSlash(p))); err != nil {
			t.Fatalf("expected file %s; err=%v", p, err)
		}
	}

	// Check rewrite produced relative links
	idxBytes, err := os.ReadFile(filepath.Join(rootDir, "index.html"))
	if err != nil {
		t.Fatal(err)
	}
	idx := string(idxBytes)
	if !strings.Contains(idx, `href="static/style.css"`) {
		t.Fatalf("index rewrite missing stylesheet rel link: %s", idx)
	}
	if !strings.Contains(idx, `src="static/app.js"`) {
		t.Fatalf("index rewrite missing script rel link: %s", idx)
	}
	if !strings.Contains(idx, `src="img/logo.png"`) {
		t.Fatalf("index rewrite missing img rel link: %s", idx)
	}
	if !strings.Contains(idx, `href="page/index.html"`) {
		t.Fatalf("index rewrite missing page rel link: %s", idx)
	}
}

func TestConcurrencyLimit(t *testing.T) {
	var current, maxSeen int32
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c := atomic.AddInt32(&current, 1)
		for {
			time.Sleep(5 * time.Millisecond)
			break
		}
		if c > atomic.LoadInt32(&maxSeen) {
			atomic.StoreInt32(&maxSeen, c)
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(`<html><body></body></html>`))
		atomic.AddInt32(&current, -1)
	}))
	defer ts.Close()

	u, _ := url.Parse(ts.URL + "/")
	out := t.TempDir()

	dl, err := New(Config{
		StartURL: u,
		MaxDepth: 0,
		OutDir:   out,
		Parallel: 2,
		Timeout:  2 * time.Second,
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := dl.Run(); err != nil {
		t.Fatal(err)
	}

	if maxSeen > 2 {
		t.Fatalf("concurrency exceeded limit: maxSeen=%d", maxSeen)
	}
}
