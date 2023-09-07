package forgeazure

import (
	"fmt"
	"strings"
)

type Execution int

const (
	ExecutionNode Execution = iota
	ExecutionNode16
	ExecutionNode10
	ExecutionPowershell
	ExecutionPowershell3
)

func ParseExecution(execution string) (Execution, error) {
	switch strings.ToLower(execution) {
	case "node":
		return ExecutionNode, nil
	case "node16":
		return ExecutionNode16, nil
	case "node10":
		return ExecutionNode10, nil
	case "powershell":
		return ExecutionPowershell, nil
	case "powershell3":
		return ExecutionPowershell3, nil
	}

	return -1, fmt.Errorf("not an exeuction: %s", execution)
}
