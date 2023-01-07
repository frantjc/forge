package envconv

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// PathFromReader takes a Reader with newline-delimited directory paths e.g.
//
//	/usr/local/bin
//	/usr/bin
//
// and returns a corresponding PATH environment variable
//
//	/usr/local/bin:/usr/bin.
func PathFromReader(r io.Reader) (string, error) {
	var (
		lines   []string
		path    = ""
		scanner = bufio.NewScanner(r)
	)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	for _, line := range lines {
		if !shouldIgnore(line) {
			cleaned := filepath.Clean(line)
			if path == "" {
				path = cleaned
			} else if !strings.Contains(path, cleaned) {
				path += ":" + cleaned
			}
		}
	}

	return path, nil
}

// PathFromReader takes a path to a file with newline-delimited directory paths e.g.
//
//	/usr/local/bin
//	/usr/bin
//
// and returns a corresponding PATH environment variable
//
//	/usr/local/bin:/usr/bin.
func PathFromFile(name string) (string, error) {
	f, err := os.Open(name)
	if err != nil {
		return "", err
	}

	return PathFromReader(f)
}

func shouldIgnore(line string) bool {
	trimmedLine := strings.TrimSpace(line)
	return trimmedLine == "" || strings.HasPrefix(trimmedLine, "#")
}
