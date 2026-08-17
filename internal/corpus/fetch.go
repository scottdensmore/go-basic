package corpus

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	// PinnedCommit is the audited BASIC Computer Games corpus revision.
	PinnedCommit = "5301155192d91d74d337899cecc59dbda59c4c17"
	// CorpusCommitFile records the revision stored in a cache directory.
	CorpusCommitFile = ".corpus-commit"
)

var (
	commitPattern = regexp.MustCompile(`^[0-9a-f]{40}$`)
	gameDirectory = regexp.MustCompile(`^[0-9]{2}_`)
)

// FetchOptions configures an atomic pinned-corpus download.
type FetchOptions struct {
	Commit string
	Target string
	URL    string
	Client *http.Client
}

// Fetch downloads only original BASIC source files and reuses a complete cache.
func Fetch(ctx context.Context, options FetchOptions) error {
	if !commitPattern.MatchString(options.Commit) {
		return fmt.Errorf("commit must be a 40-character lowercase hexadecimal SHA")
	}
	target, err := validateTarget(options.Target)
	if err != nil {
		return err
	}
	if cached, err := cachedCommit(target); err != nil {
		return err
	} else if cached != "" {
		if cached != options.Commit {
			return fmt.Errorf("cache %s contains commit %s, want %s", target, cached, options.Commit)
		}
		return nil
	}
	if _, err := os.Stat(target); err == nil {
		return fmt.Errorf("cache target %s exists without %s", target, CorpusCommitFile)
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("inspect cache target: %w", err)
	}

	parent := filepath.Dir(target)
	if err := os.MkdirAll(parent, 0o755); err != nil {
		return fmt.Errorf("create cache parent: %w", err)
	}
	temporary, err := os.MkdirTemp(parent, ".basic-computer-games-")
	if err != nil {
		return fmt.Errorf("create temporary cache: %w", err)
	}
	defer func() {
		if temporary != "" {
			_ = os.RemoveAll(temporary)
		}
	}()

	url := options.URL
	if url == "" {
		url = "https://codeload.github.com/coding-horror/basic-computer-games/tar.gz/" + options.Commit
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("create corpus request: %w", err)
	}
	client := options.Client
	if client == nil {
		client = http.DefaultClient
	}
	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("download corpus: %w", err)
	}
	defer func() {
		_ = response.Body.Close()
	}()
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("download corpus: unexpected HTTP status %s", response.Status)
	}
	if err := extractOriginalBASIC(response.Body, temporary); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(temporary, CorpusCommitFile), []byte(options.Commit+"\n"), 0o644); err != nil {
		return fmt.Errorf("record corpus commit: %w", err)
	}
	if err := os.Rename(temporary, target); err != nil {
		return fmt.Errorf("publish corpus cache: %w", err)
	}
	temporary = ""
	return nil
}

func validateTarget(target string) (string, error) {
	if strings.TrimSpace(target) == "" {
		return "", errors.New("cache target is required")
	}
	cleaned, err := filepath.Abs(filepath.Clean(target))
	if err != nil {
		return "", fmt.Errorf("resolve cache target: %w", err)
	}
	if cleaned == filepath.VolumeName(cleaned)+string(filepath.Separator) {
		return "", errors.New("cache target cannot be a filesystem root")
	}
	return cleaned, nil
}

func cachedCommit(target string) (string, error) {
	contents, err := os.ReadFile(filepath.Join(target, CorpusCommitFile))
	if os.IsNotExist(err) {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("read cached corpus commit: %w", err)
	}
	return strings.TrimSpace(string(contents)), nil
}

// VerifyCommit confirms that root is a complete cache for expected.
func VerifyCommit(root, expected string) error {
	cached, err := cachedCommit(root)
	if err != nil {
		return err
	}
	if cached == "" {
		return fmt.Errorf("corpus cache %s does not record a commit", root)
	}
	if cached != expected {
		return fmt.Errorf("corpus cache %s contains commit %s, want %s", root, cached, expected)
	}
	return nil
}

func extractOriginalBASIC(reader io.Reader, target string) error {
	compressed, err := gzip.NewReader(reader)
	if err != nil {
		return fmt.Errorf("open corpus archive: %w", err)
	}
	defer func() {
		_ = compressed.Close()
	}()
	archive := tar.NewReader(compressed)
	extracted := 0
	for {
		header, err := archive.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return fmt.Errorf("read corpus archive: %w", err)
		}
		if header.Typeflag != tar.TypeReg {
			continue
		}
		path, selected := originalBASICPath(header.Name)
		if !selected {
			continue
		}
		destination := filepath.Join(target, filepath.FromSlash(path))
		if err := os.MkdirAll(filepath.Dir(destination), 0o755); err != nil {
			return fmt.Errorf("create corpus directory: %w", err)
		}
		file, err := os.OpenFile(destination, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
		if err != nil {
			return fmt.Errorf("create corpus file %s: %w", path, err)
		}
		_, copyErr := io.Copy(file, archive)
		closeErr := file.Close()
		if copyErr != nil {
			return fmt.Errorf("extract corpus file %s: %w", path, copyErr)
		}
		if closeErr != nil {
			return fmt.Errorf("close corpus file %s: %w", path, closeErr)
		}
		extracted++
	}
	if extracted == 0 {
		return errors.New("corpus archive contained no original BASIC files")
	}
	return nil
}

func originalBASICPath(name string) (string, bool) {
	cleaned := filepath.ToSlash(filepath.Clean(name))
	if cleaned == "." || strings.HasPrefix(cleaned, "../") || strings.HasPrefix(cleaned, "/") {
		return "", false
	}
	parts := strings.Split(cleaned, "/")
	if len(parts) < 3 {
		return "", false
	}
	parts = parts[1:]
	if len(parts) == 2 && parts[0] != "00_Alternate_Languages" &&
		gameDirectory.MatchString(parts[0]) && strings.EqualFold(filepath.Ext(parts[1]), ".bas") {
		return strings.Join(parts, "/"), true
	}
	if len(parts) == 3 && parts[0] == "00_Alternate_Languages" &&
		gameDirectory.MatchString(parts[1]) && strings.EqualFold(filepath.Ext(parts[2]), ".bas") {
		return strings.Join(parts, "/"), true
	}
	return "", false
}
