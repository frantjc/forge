package githubactions

const (
	// SecretActionsStepDebug--see https://docs.github.com/en/actions/monitoring-and-troubleshooting-workflows/enabling-debug-logging.
	SecretActionsStepDebug = "ACTIONS_STEP_DEBUG"
	// SecretActionsRunnerDebug--see https://docs.github.com/en/actions/monitoring-and-troubleshooting-workflows/enabling-debug-logging.
	SecretActionsRunnerDebug = "ACTIONS_RUNNER_DEBUG"
	// SecretRunnerDebug--see https://github.com/actions/toolkit/blob/master/packages/core/src/core.ts#L118-L123.
	SecretRunnerDebug = "RUNNER_DEBUG"
)
