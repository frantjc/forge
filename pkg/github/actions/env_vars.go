package actions

const (
	EnvVarCI              = "CI"
	EnvVarWorkflow        = "GITHUB_WORKFLOW"
	EnvVarRunID           = "GITHUB_RUN_ID"
	EnvVarRunNumber       = "GITHUB_RUN_NUMBER"
	EnvVarRunAttempt      = "GITHUB_RUN_ATTEMPT"
	EnvVarJob             = "GITHUB_JOB"
	EnvVarAction          = "GITHUB_ACTION"
	EnvVarActionPath      = "GITHUB_ACTION_PATH"
	EnvVarActions         = "GITHUB_ACTIONS"
	EnvVarActor           = "GITHUB_ACTOR"
	EnvVarRepository      = "GITHUB_REPOSITORY"
	EnvVarEventName       = "GITHUB_EVENT_NAME"
	EnvVarEventPath       = "GITHUB_EVENT_PATH"
	EnvVarWorkspace       = "GITHUB_WORKSPACE"
	EnvVarSha             = "GITHUB_SHA"
	EnvVarRef             = "GITHUB_REF"
	EnvVarRefName         = "GITHUB_REF_NAME"
	EnvVarRefProtected    = "GITHUB_REF_PROTECTED"
	EnvVarRefType         = "GITHUB_REF_TYPE"
	EnvVarHeadRef         = "GITHUB_HEAD_REF"
	EnvVarBaseRef         = "GITHUB_BASE_REF"
	EnvVarServerURL       = "GITHUB_SERVER_URL"
	EnvVarAPIURL          = "GITHUB_API_URL"
	EnvVarGraphQLURL      = "GITHUB_GRAPHQL_URL"
	EnvVarRunnerName      = "RUNNER_NAME"
	EnvVarRunnerOS        = "RUNNER_OS"
	EnvVarRunnerArch      = "RUNNER_ARCH"
	EnvVarRunnerTemp      = "RUNNER_TEMP"
	EnvVarRunnerToolCache = "RUNNER_TOOL_CACHE"

	EnvVarEnv  = "GITHUB_ENV"
	EnvVarPath = "GITHUB_PATH"

	EnvVarToken = "GITHUB_TOKEN" //nolint:gosec

	EnvVarRepositoryOwner  = "GITHUB_REPOSITORY_OWNER"
	EnvVarRetentionDays    = "GITHUB_RETENTION_DAYS"
	EnvVarStepSummary      = "GITHUB_STEP_SUMMARY"
	EnvVarActionRepository = "GITHUB_ACTION_REPOSITORY"
)
