package cloudbuild

const (
	EnvVarProjectID              = "PROJECT_ID"
	EnvVarBuildID                = "BUILD_ID"
	EnvVarProjectNumber          = "PROJECT_NUMBER"
	EnvVarLocation               = "LOCATION"
	EnvVarTriggerName            = "TRIGGER_NAME"
	EnvVarCommitSha              = "COMMIT_SHA"
	EnvVarRevisionID             = "REVISION_ID"
	EnvVarShortSha               = "SHORT_SHA"
	EnvVarRepoName               = "REPO_NAME"
	EnvVarRepoFullName           = "REPO_FULL_NAME"
	EnvVarBranchName             = "BRANCH_NAME"
	EnvVarTagName                = "TAG_NAME"
	EnvVarRefName                = "REF_NAME"
	EnvVarTriggerBuildConfigPath = "TRIGGER_BUILD_CONFIG_PATH"
	EnvVarServiceAccountEmail    = "SERVICE_ACCOUNT_EMAIL"
	EnvVarServiceAccount         = "SERVICE_ACCOUNT"

	EnvVarGitHubHeadBranch  = "_HEAD_BRANCH"
	EnvVarGitHubBaseBranch  = "_BASE_BRANCH"
	EnvVarGitHubHeadRepoURL = "_HEAD_REPO_URL"
	EnvVarGitHubPRNumber    = "_PR_NUMBER"
)
