package githubactions

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"
)

var (
	valueDelimiter  = `ghadelimiter_[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}`
	valueDelimiterR = regexp.MustCompile(valueDelimiter)

	keyDelimiter  = `(\w+)<<` + valueDelimiter
	keyDelimiterR = regexp.MustCompile(keyDelimiter)
)

// ParseEnvFile reads a GitHub Actions environment file e.g.
//
//		KEY1=value1
//		# SKIPPED=comment
//	 KEY2<<ghadelimiter_CC992248-87BA-41AF-BF33-A52DCE9681A6
//		value2 # For some security reason(s), @actions/core writes to the environment this way.
//		ghadelimiter_CC992248-87BA-41AF-BF33-A52DCE9681A6
func ParseEnvFile(r io.Reader) (map[string]string, error) {
	var (
		values  = make(map[string]string)
		scanner = bufio.NewScanner(r)
	)
	for scanner.Scan() {
		line0 := scanner.Text()
		if strings.HasPrefix(line0, "# ") || strings.TrimSpace(line0) == "" {
			continue
		} else if matches := keyDelimiterR.FindStringSubmatch(line0); len(matches) == 2 {
			if scanner.Scan() {
				line1 := scanner.Text()
				if scanner.Scan() {
					line2 := scanner.Text()
					if valueDelimiterR.MatchString(line2) {
						values[matches[1]] = strings.SplitN(line1, " #", 2)[0]
						continue
					}
				}
			}

			return nil, fmt.Errorf("invalid multiline environment file entry")
		} else if matches := strings.SplitN(line0, "=", 2); len(matches) == 2 {
			values[matches[0]] = trimValue(matches[1])
		} else {
			return nil, fmt.Errorf("parse environment file line: %s", line0)
		}
	}

	return values, scanner.Err()
}

func trimValue(value string) string {
	if strings.HasPrefix(value, `"`) && strings.HasSuffix(value, `"`) {
		return strings.Trim(value, `"`)
	} else if strings.HasPrefix(value, `'`) && strings.HasSuffix(value, `'`) {
		return strings.Trim(value, `'`)
	}

	return value
}
