package actions

import (
	"io"
	"strings"
)

type workflowCommandWriter struct {
	callback func(*WorkflowCommand) []byte
	w        io.Writer
}

// TODO bufio.ScanLines.
func (w *workflowCommandWriter) Write(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}

	lines := strings.Split(string(p), "\n")
	for i, line := range lines {
		a := []byte{}
		if i < len(lines)-1 {
			a = []byte{'\n'}
		}

		if len(line) == 0 || strings.HasPrefix(line, "##[add-matcher]") {
		} else if c, err := ParseWorkflowCommandString(line); err == nil {
			if b := w.callback(c); len(b) != 0 {
				if _, err = w.w.Write(append(b, a...)); err != nil {
					return len(p), err
				}
			}
		} else {
			if _, err := w.w.Write(append([]byte(line), a...)); err != nil {
				return len(p), err
			}
		}
	}

	return len(p), nil
}

func NewWorkflowCommandWriter(callback func(*WorkflowCommand) []byte, w io.Writer) io.Writer {
	return &workflowCommandWriter{callback, w}
}
