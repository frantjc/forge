package cloudbuild

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/frantjc/forge/envconv"
	xslice "github.com/frantjc/x/slice"
	"github.com/go-git/go-git/v5"
	"golang.org/x/exp/maps"
)

var (
	configR                  = regexp.MustCompile(`([^\s]+)\s*=\s*([^\s]+)\n`)
	userDefinedSubstitutionR = regexp.MustCompile(`_[A-Z0-9_]+`)
)

type Substitutions struct {
	ProjectID              string
	BuildID                string
	ProjectNumber          int
	Location               string
	TriggerName            string
	CommitSha              string
	RevisionID             string
	ShortSha               string
	RepoName               string
	RepoFullName           string
	BranchName             string
	TagName                string
	RefName                string
	TriggerBuildConfigPath string
	ServiceAccountEmail    string
	ServiceAccount         string
	GitHubHeadBranch       string
	GitHubBaseBranch       string
	GitHubHeadRepoURL      string
	GitHubPRNumber         int
	UserDefined            map[string]string
}

func (s *Substitutions) Env() []string {
	projectNumber := ""
	if s.ProjectNumber > 0 {
		projectNumber = fmt.Sprint(s.ProjectNumber)
	}

	githubPRNumber := ""
	if s.GitHubPRNumber > 0 {
		githubPRNumber = fmt.Sprint(s.ProjectNumber)
	}

	substitutionsM := map[string]string{
		EnvVarProjectID:              s.ProjectID,
		EnvVarBuildID:                s.BuildID,
		EnvVarProjectNumber:          projectNumber,
		EnvVarLocation:               s.Location,
		EnvVarTriggerName:            s.TriggerName,
		EnvVarCommitSha:              xslice.Coalesce(s.CommitSha, s.RevisionID),
		EnvVarRevisionID:             xslice.Coalesce(s.RevisionID, s.CommitSha),
		EnvVarRepoName:               xslice.Coalesce(s.RepoName, strings.Split(s.RepoFullName, "/")[1]),
		EnvVarRepoFullName:           s.RepoFullName,
		EnvVarBranchName:             xslice.Coalesce(s.BranchName, s.RefName),
		EnvVarTagName:                xslice.Coalesce(s.TagName, s.RefName),
		EnvVarRefName:                xslice.Coalesce(s.RefName, s.BranchName, s.TagName),
		EnvVarTriggerBuildConfigPath: s.TriggerBuildConfigPath,
		EnvVarServiceAccountEmail:    s.ServiceAccountEmail,
		EnvVarServiceAccount:         s.ServiceAccountEmail,
		EnvVarGitHubHeadBranch:       s.GitHubHeadBranch,
		EnvVarGitHubBaseBranch:       s.GitHubBaseBranch,
		EnvVarGitHubHeadRepoURL:      s.GitHubHeadRepoURL,
		EnvVarGitHubPRNumber:         githubPRNumber,
	}

	i := len(substitutionsM[EnvVarCommitSha])
	if i > 7 {
		i = 7
	}

	j := len(substitutionsM[EnvVarRevisionID])
	if j > 7 {
		j = 7
	}

	substitutionsM[EnvVarShortSha] = xslice.Coalesce(s.ShortSha, s.CommitSha[:i], s.RevisionID[:j])

	maps.Copy(substitutionsM, s.UserDefined)

	return envconv.MapToArr(substitutionsM)
}

// NewSubstitutionsFromEnv returns a map of default substitutions
// whose values are sourced from the environment.
// See https://cloud.google.com/build/docs/configuring-builds/substitute-variable-values#using_default_substitutions.
func NewSubstitutionsFromEnv(userDefinedSubstitutions map[string]string) (*Substitutions, error) { //nolint: gocyclo
	substitutionsM := map[string]string{
		EnvVarProjectID:              os.Getenv(EnvVarProjectID),
		EnvVarBuildID:                os.Getenv(EnvVarBuildID),
		EnvVarProjectNumber:          os.Getenv(EnvVarProjectNumber),
		EnvVarLocation:               os.Getenv(EnvVarLocation),
		EnvVarTriggerName:            os.Getenv(EnvVarTriggerName),
		EnvVarCommitSha:              os.Getenv(EnvVarCommitSha),
		EnvVarRevisionID:             os.Getenv(EnvVarRevisionID),
		EnvVarShortSha:               os.Getenv(EnvVarShortSha),
		EnvVarRepoName:               os.Getenv(EnvVarRepoName),
		EnvVarRepoFullName:           os.Getenv(EnvVarRepoFullName),
		EnvVarBranchName:             os.Getenv(EnvVarBranchName),
		EnvVarTagName:                os.Getenv(EnvVarTagName),
		EnvVarRefName:                os.Getenv(EnvVarRefName),
		EnvVarTriggerBuildConfigPath: os.Getenv(EnvVarTriggerBuildConfigPath),
		EnvVarServiceAccountEmail:    os.Getenv(EnvVarServiceAccountEmail),
		EnvVarServiceAccount:         os.Getenv(EnvVarServiceAccount),
		EnvVarGitHubHeadBranch:       os.Getenv(EnvVarGitHubHeadBranch),
		EnvVarGitHubBaseBranch:       os.Getenv(EnvVarGitHubBaseBranch),
		EnvVarGitHubHeadRepoURL:      os.Getenv(EnvVarGitHubHeadRepoURL),
		EnvVarGitHubPRNumber:         os.Getenv(EnvVarGitHubPRNumber),
	}

	home := os.Getenv("HOME")
	if home == "" {
		if u, err := user.Current(); err == nil {
			home = u.HomeDir
		}
	}

	if home != "" {
		activeConfig := "default"
		if activeConfigFile, err := os.Open(filepath.Join(home, ".config/gcloud/active_config")); err == nil {
			if b, err := io.ReadAll(activeConfigFile); err == nil {
				activeConfig = strings.TrimSpace(string(b))
			}
		}

		if configFile, err := os.Open(filepath.Join(home, fmt.Sprintf(".config/gcloud/configurations/config_%s", activeConfig))); err == nil {
			if b, err := io.ReadAll(configFile); err == nil {
				matches := configR.FindAllStringSubmatch(string(b), -1)

				for _, line := range matches {
					if len(line) == 3 {
						switch line[1] {
						case "project":
							if substitutionsM[EnvVarProjectID] == "" {
								substitutionsM[EnvVarProjectID] = line[2]
							}
						case "account":
							if substitutionsM[EnvVarServiceAccount] == "" {
								substitutionsM[EnvVarServiceAccount] = line[2]
							}

							if substitutionsM[EnvVarServiceAccountEmail] == "" {
								substitutionsM[EnvVarServiceAccountEmail] = line[2]
							}
						}
					}
				}
			}
		}
	}

	for k, v := range userDefinedSubstitutions {
		if _, ok := substitutionsM[k]; !ok && !userDefinedSubstitutionR.MatchString(k) {
			return nil, fmt.Errorf("user-defined substitution must respect regular expression `%s`: %s", userDefinedSubstitutionR, k)
		}

		substitutionsM[k] = v
	}

	substitutions := &Substitutions{
		UserDefined: make(map[string]string),
	}

	for k, v := range substitutionsM {
		switch k {
		case EnvVarProjectID:
			substitutions.ProjectID = v
		case EnvVarBuildID:
			substitutions.BuildID = v
		case EnvVarProjectNumber:
			substitutions.ProjectNumber, _ = strconv.Atoi(v)
		case EnvVarLocation:
			substitutions.Location = v
		case EnvVarTriggerName:
			substitutions.TriggerName = v
		case EnvVarCommitSha:
			substitutions.CommitSha = v
		case EnvVarRevisionID:
			substitutions.RevisionID = v
		case EnvVarShortSha:
			substitutions.ShortSha = v
		case EnvVarRepoName:
			substitutions.RepoName = v
		case EnvVarRepoFullName:
			substitutions.RepoFullName = v
		case EnvVarBranchName:
			substitutions.BranchName = v
		case EnvVarTagName:
			substitutions.TagName = v
		case EnvVarRefName:
			substitutions.RefName = v
		case EnvVarTriggerBuildConfigPath:
			substitutions.TriggerBuildConfigPath = v
		case EnvVarServiceAccountEmail:
			substitutions.ServiceAccountEmail = v
		case EnvVarServiceAccount:
			substitutions.ServiceAccount = v
		case EnvVarGitHubHeadBranch:
			substitutions.GitHubHeadBranch = v
		case EnvVarGitHubBaseBranch:
			substitutions.GitHubBaseBranch = v
		case EnvVarGitHubHeadRepoURL:
			substitutions.GitHubHeadRepoURL = v
		case EnvVarGitHubPRNumber:
			substitutions.GitHubPRNumber, _ = strconv.Atoi(v)
		default:
			substitutions.UserDefined[k] = v
		}
	}

	return substitutions, nil
}

// NewSubstituionsFromPath returns a map of default
// substitutions whose values are sourced from the git
// repository at the given path and the environment.
// See https://cloud.google.com/build/docs/configuring-builds/substitute-variable-values#using_default_substitutions.
func NewSubstituionsFromPath(path string, userDefinedSubstitutions map[string]string) (*Substitutions, error) {
	r, err := git.PlainOpen(path)
	for ; err != nil && path != "/"; r, err = git.PlainOpen(path) {
		path = filepath.Dir(path)
	}
	if err != nil {
		return nil, err
	}

	substitutions, err := NewSubstitutionsFromEnv(userDefinedSubstitutions)
	if err != nil {
		return nil, err
	}

	if ref, err := r.Head(); err == nil {
		if sha := ref.Hash().String(); sha != "" {
			substitutions.CommitSha = xslice.Coalesce(substitutions.CommitSha, sha)
			substitutions.RevisionID = xslice.Coalesce(substitutions.RevisionID, sha)
			i := len(sha)
			if i > 7 {
				i = 7
			}
			substitutions.ShortSha = xslice.Coalesce(substitutions.ShortSha, sha[:i])
		}

		refName := ref.Name().Short()
		substitutions.RefName = xslice.Coalesce(substitutions.RefName, refName)
		if ref.Name().IsBranch() {
			substitutions.BranchName = xslice.Coalesce(substitutions.BranchName, refName)
		} else {
			substitutions.TagName = xslice.Coalesce(substitutions.TagName, refName)
		}

		if conf, err := r.Config(); err == nil {
			for _, remote := range conf.Remotes {
				for _, rurl := range remote.URLs {
					if prurl, err := url.Parse(rurl); err == nil {
						substitutions.GitHubHeadRepoURL = xslice.Coalesce(substitutions.GitHubHeadRepoURL, prurl.String())
						substitutions.RepoFullName = xslice.Coalesce(
							substitutions.RepoFullName,
							strings.TrimSuffix(
								strings.TrimPrefix(prurl.Path, "/"),
								".git",
							),
						)
						substitutions.RepoName = xslice.Coalesce(
							substitutions.RepoName,
							strings.Split(substitutions.RepoFullName, "/")[1],
						)
						break
					}
				}
			}
		}
	}

	return substitutions, nil
}
