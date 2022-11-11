package githubactions

import "errors"

var (
	ErrNotAnAction         = errors.New("action.yaml/action.yml not found")
	ErrNotAWorkflowCommand = errors.New("not a workflow command")
)

func IsErrNotAnAction(err error) bool {
	return errors.Is(err, ErrNotAnAction)
}

func IsErrNotAWorkflowCommand(err error) bool {
	return errors.Is(err, ErrNotAWorkflowCommand)
}
