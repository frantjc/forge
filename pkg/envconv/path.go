package envconv

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
	"strings"
)

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

	for i, line := range lines {
		if !shouldIgnore(line) {
			cleaned := filepath.Clean(line)
			if i == 0 {
				path = cleaned
			} else if !strings.Contains(path, cleaned) {
				path += ":" + cleaned
			}
		}
	}

	return path, nil
}

func PathFromFile(name string) (string, error) {
	f, err := os.Open(name)
	if err != nil {
		return "", err
	}

	return PathFromReader(f)
}

func shouldIgnore(line string) bool {
	trimmedLine := strings.TrimSpace(line)
	return len(trimmedLine) == 0 || strings.HasPrefix(trimmedLine, "#")
}
