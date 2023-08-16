package ore

import (
	"context"
	"fmt"

	"github.com/frantjc/forge"
)

type Task struct {
	Task   string            `json:"task,omitempty"`
	Inputs map[string]string `json:"inputs,omitempty"`
}

func (o *Task) Liquify(_ context.Context, _ forge.ContainerRuntime, _ *forge.Drains) error {
	return fmt.Errorf("unimplemented")
}
