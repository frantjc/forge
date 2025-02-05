package forge

import (
	xslice "github.com/frantjc/x/slice"
)

type Mount struct {
	Source      string `json:"source,omitempty"`
	Destination string `json:"destination,omitempty"`
}

func overrideMounts(oldMounts []Mount, newMounts ...Mount) []Mount {
	return append(xslice.Filter(oldMounts, func(m Mount, _ int) bool {
		return !xslice.Some(newMounts, func(n Mount, _ int) bool {
			return m.Destination == n.Destination
		})
	}), newMounts...)
}
