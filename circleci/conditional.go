package circleci

import "strconv"

// Conditional is a `when` or `unless` in a CircleCI Orb step.
type Conditional struct {
	Condition any    `json:"condition,omitempty" yaml:",omitempty"`
	Steps     []Step `json:"steps,omitempty" yaml:",omitempty"`
}

// EvaluateConditional evaluates a Conditional, expanding any
// variables using the given ExpandFunc before evaluating.
func EvaluateConditional(expand ExpandFunc, c *Conditional) bool {
	switch v := c.Condition.(type) {
	case bool:
		return v
	case string:
		b, _ := strconv.ParseBool(expand(v))

		return b
	case map[string]any:
		if eq, ok := v["equal"]; ok {
			if equal, ok := eq.([]string); ok {
				lenEqual := len(equal)
				if lenEqual < 2 {
					return true
				}

				for i := 1; i < lenEqual; i++ {
					if equal[i] != equal[i-1] {
						return false
					}
				}

				return true
			}
		}
	}

	return false
}
