// Package corpus tests pinned external-corpus acquisition and execution.
package corpus

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func TestFetchCachesOnlyOriginalBASICFiles(t *testing.T) {
	archive := corpusArchive(t, map[string]string{
		"basic-computer-games-pinned/01_Game/game.bas":                              "10 END\n",
		"basic-computer-games-pinned/01_Game/python/game.py":                        "ignored",
		"basic-computer-games-pinned/00_Alternate_Languages/01_Game/game.bas":       "20 END\n",
		"basic-computer-games-pinned/00_Alternate_Languages/01_Game/csharp/game.cs": "ignored",
		"basic-computer-games-pinned/README.md":                                     "ignored",
	})
	var requests atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, _ *http.Request) {
		requests.Add(1)
		response.Header().Set("Content-Type", "application/gzip")
		_, _ = response.Write(archive)
	}))
	t.Cleanup(server.Close)

	target := filepath.Join(t.TempDir(), "cache", "pinned")
	if err := Fetch(context.Background(), FetchOptions{
		Commit: "0123456789abcdef0123456789abcdef01234567",
		Target: target,
		URL:    server.URL,
	}); err != nil {
		t.Fatal(err)
	}
	if err := Fetch(context.Background(), FetchOptions{
		Commit: "0123456789abcdef0123456789abcdef01234567",
		Target: target,
		URL:    server.URL,
	}); err != nil {
		t.Fatal(err)
	}

	for _, path := range []string{
		"01_Game/game.bas",
		"00_Alternate_Languages/01_Game/game.bas",
		CorpusCommitFile,
	} {
		if _, err := os.Stat(filepath.Join(target, path)); err != nil {
			t.Fatalf("selected path %s: %v", path, err)
		}
	}
	for _, path := range []string{"01_Game/python/game.py", "README.md"} {
		if _, err := os.Stat(filepath.Join(target, path)); !os.IsNotExist(err) {
			t.Fatalf("unrelated path %s was cached", path)
		}
	}
	if got, want := requests.Load(), int32(1); got != want {
		t.Fatalf("requests: got %d, want %d", got, want)
	}
}

func TestDiscoverIncludesMainAndOnlyByteDifferentAlternates(t *testing.T) {
	root := t.TempDir()
	writeCorpusFile(t, root, "01_One/one.bas", "10 END\n")
	writeCorpusFile(t, root, "00_Alternate_Languages/01_One/one.bas", "10 END\n")
	writeCorpusFile(t, root, "02_Two/two.bas", "20 END\n")
	writeCorpusFile(t, root, "00_Alternate_Languages/02_Two/two.bas", "25 END\n")
	writeCorpusFile(t, root, "03_Three/first.bas", "30 END\n")
	writeCorpusFile(t, root, "03_Three/second.bas", "40 END\n")
	writeCorpusFile(t, root, "03_Three/python/translated.bas", "ignored\n")

	variants, err := Discover(root)
	if err != nil {
		t.Fatal(err)
	}
	paths := make([]string, 0, len(variants))
	for _, variant := range variants {
		paths = append(paths, variant.Path)
	}
	want := []string{
		"01_One/one.bas",
		"02_Two/two.bas",
		"03_Three/first.bas",
		"03_Three/second.bas",
		"00_Alternate_Languages/02_Two/two.bas",
	}
	if got := strings.Join(paths, "\n"); got != strings.Join(want, "\n") {
		t.Fatalf("paths:\ngot:\n%s\nwant:\n%s", got, strings.Join(want, "\n"))
	}
}

func TestSmokeReportsEachVariantAndActionableFailures(t *testing.T) {
	root := t.TempDir()
	writeCorpusFile(t, root, "01_Complete/complete.bas", "10 PRINT \"OK\"\n20 END\n")
	writeCorpusFile(t, root, "02_Input/input.bas", "10 INPUT A\n")
	writeCorpusFile(t, root, "03_Bounded/bounded.bas", "10 GOTO 10\n")
	writeCorpusFile(t, root, "04_Broken/broken.bas", "10 GOTO 99\n")
	variants, err := Discover(root)
	if err != nil {
		t.Fatal(err)
	}

	results := Smoke(context.Background(), root, variants, SmokeOptions{
		StatementLimit: 5,
		Timeout:        time.Second,
	})
	if got, want := len(results), 4; got != want {
		t.Fatalf("results: got %d, want %d", got, want)
	}
	for _, result := range results {
		switch result.Path {
		case "01_Complete/complete.bas":
			if result.Status != StatusComplete || result.LastLine != 20 || result.Err != nil {
				t.Fatalf("complete result: %#v", result)
			}
		case "02_Input/input.bas":
			if result.Status != StatusInput || result.LastLine != 10 || result.Err != nil {
				t.Fatalf("input result: %#v", result)
			}
		case "03_Bounded/bounded.bas":
			if result.Status != StatusBounded || result.LastLine != 10 || result.Err != nil {
				t.Fatalf("bounded result: %#v", result)
			}
		case "04_Broken/broken.bas":
			if result.Status != StatusFailed || result.LastLine != 10 ||
				result.Err == nil || !strings.Contains(result.Err.Error(), "undefined BASIC line 99") {
				t.Fatalf("broken result: %#v", result)
			}
		default:
			t.Fatalf("unexpected result: %#v", result)
		}
	}
}

func corpusArchive(t *testing.T, files map[string]string) []byte {
	t.Helper()
	var compressed bytes.Buffer
	gzipWriter := gzip.NewWriter(&compressed)
	tarWriter := tar.NewWriter(gzipWriter)
	for path, contents := range files {
		if err := tarWriter.WriteHeader(&tar.Header{
			Name: path,
			Mode: 0o644,
			Size: int64(len(contents)),
		}); err != nil {
			t.Fatal(err)
		}
		if _, err := io.WriteString(tarWriter, contents); err != nil {
			t.Fatal(err)
		}
	}
	if err := tarWriter.Close(); err != nil {
		t.Fatal(err)
	}
	if err := gzipWriter.Close(); err != nil {
		t.Fatal(err)
	}
	return compressed.Bytes()
}

func writeCorpusFile(t *testing.T, root, path, contents string) {
	t.Helper()
	fullPath := filepath.Join(root, filepath.FromSlash(path))
	if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(fullPath, []byte(contents), 0o600); err != nil {
		t.Fatal(err)
	}
}
