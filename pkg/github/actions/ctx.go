package actions

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"os/user"
	"strconv"
	"strings"

	"github.com/frantjc/forge/pkg/envconv"
	"github.com/frantjc/go-js"
	"github.com/go-git/go-git/v5"
)

type globalContextKey struct{}

func WithGlobalContext(ctx context.Context, globalContext *GlobalContext) context.Context {
	return context.WithValue(ctx, globalContextKey{}, globalContext)
}

func GlobalContextFrom(ctx context.Context) (globalContext *GlobalContext, ok bool) {
	globalContext, ok = ctx.Value(globalContextKey{}).(*GlobalContext)
	return
}

func (c *GlobalContext) GetString(key string) string {
	keys := strings.Split(key, ".")
	if len(keys) > 0 {
		switch keys[0] {
		case "github":
			if len(keys) > 1 {
				return c.GetGitHubContext().GetString(strings.Join(keys[1:], "."))
			}
		case "env":
			if len(keys) > 1 {
				if v, ok := c.EnvContext[keys[1]]; ok {
					return v
				}
			}
		case "job":
			if len(keys) > 1 {
				return c.GetJobContext().GetString(strings.Join(keys[1:], "."))
			}
		case "steps":
			if len(keys) > 2 {
				if v, ok := c.StepsContext[keys[1]]; ok {
					return v.GetString(strings.Join(keys[2:], "."))
				}
			}
		case "runner":
			if len(keys) > 1 {
				return c.GetRunnerContext().GetString(strings.Join(keys[1:], "."))
			}
		case "inputs":
			if len(keys) > 1 {
				if v, ok := c.InputsContext[keys[1]]; ok {
					return v
				}
			}
		case "secrets":
			if len(keys) > 1 {
				if v, ok := c.SecretsContext[keys[1]]; ok {
					return v
				}
			}
		case "needs":
			if len(keys) > 2 {
				if v, ok := c.NeedsContext[keys[1]]; ok {
					return v.GetString(strings.Join(keys[2:], "."))
				}
			}
		}
	}

	return ""
}

func (c *GitHubContext) GetString(key string) string {
	keys := strings.Split(key, ".")
	if len(keys) > 0 {
		switch keys[0] {
		case "action":
			return c.GetAction()
		case "action_path":
			return c.GetActionPath()
		case "actor":
			return c.GetActor()
		case "base_ref":
			return c.GetBaseRef()
		case "event":
			return c.GetEvent()
		case "event_name":
			return c.GetEventName()
		case "event_path":
			return c.GetEventPath()
		case "head_ref":
			return c.GetHeadRef()
		case "job":
			return c.GetJob()
		case "ref":
			return c.GetRef()
		case "ref_name":
			return c.GetRefName()
		case "ref_protected":
			return fmt.Sprint(c.GetRefProtected())
		case "ref_type":
			return c.GetRefType()
		case "repository":
			return c.GetRepository()
		case "repository_owner":
			return c.GetRepositoryOwner()
		case "run_id":
			return c.GetRunId()
		case "run_number":
			return fmt.Sprint(c.GetRunNumber())
		case "run_attempt":
			return fmt.Sprint(c.GetRunAttempt())
		case "server_url":
			return c.GetServerUrl()
		case "sha":
			return c.GetSha()
		case "token":
			return c.GetToken()
		case "workflow":
			return c.GetWorkflow()
		case "workspace":
			return c.GetWorkspace()
		}
	}

	return ""
}

func (c *JobContext) GetString(key string) string {
	keys := strings.Split(key, ".")
	if len(keys) > 0 {
		switch keys[0] {
		case "container":
			if len(keys) > 1 {
				switch keys[1] {
				case "id":
					return c.Container.Id
				case "network":
					return c.Container.Network
				}
			}
		case "services":
			if len(keys) > 1 {
				if v, ok := c.Services[keys[1]]; ok {
					if len(keys) > 2 {
						switch keys[2] {
						case "id":
							return v.Id
						case "network":
							return v.Network
						case "ports":
							if len(keys) > 4 {
								if v, ok := v.Ports[keys[4]]; ok {
									return v
								}
							}
						}
					}
				}
			}
		case "status":
			return c.Status
		}
	}

	return ""
}

func (c *StepContext) GetString(key string) string {
	keys := strings.Split(key, ".")
	if len(keys) > 0 {
		switch keys[0] {
		case "outputs":
			if len(keys) > 1 {
				if v, ok := c.Outputs[keys[1]]; ok {
					return v
				}
			}
		case "outcome":
			return c.GetOutcome()
		case "conclusion":
			return c.GetConclusion()
		}
	}

	return ""
}

func (c *RunnerContext) GetString(key string) string {
	keys := strings.Split(key, ".")
	if len(keys) > 0 {
		switch keys[0] {
		case "name":
			return c.GetName()
		case "os":
			return c.GetOs()
		case "arch":
			return c.GetArch()
		case "temp":
			return c.GetTemp()
		case "tool_cache":
			return c.GetToolCache()
		}
	}

	return ""
}

func (c *NeedContext) GetString(key string) string {
	keys := strings.Split(key, ".")
	if len(keys) > 0 {
		if keys[0] == "outputs" && len(keys) > 1 {
			if v, ok := c.Outputs[keys[1]]; ok {
				return v
			}
		}
	}

	return ""
}

func (c *GlobalContext) AddEnv(env map[string]string) {
	if len(c.EnvContext) == 0 {
		c.EnvContext = env
		return
	}

	for k, v := range env {
		c.EnvContext[k] = v
	}
}

func NewGlobalContextFromEnv() *GlobalContext {
	u, _ := user.Current()

	actor := os.Getenv(EnvVarActor)
	if actor == "" {
		actor = u.Name
	}

	rawServerURL := os.Getenv(EnvVarServerURL)
	if rawServerURL == "" {
		rawServerURL = DefaultURL.String()
	}

	serverURL, err := url.Parse(rawServerURL)
	if err != nil {
		serverURL = DefaultURL
	}

	runnerOS := OS(os.Getenv(EnvVarRunnerOS))
	if runnerOS == "" {
		runnerOS = OSLinux
	}

	runnerArch := Arch(os.Getenv(EnvVarRunnerArch))
	if runnerArch == "" {
		runnerArch = ArchX86
	}

	refProtected, _ := strconv.ParseBool(EnvVarRefProtected)

	runNumber, err := strconv.Atoi(os.Getenv(EnvVarRunNumber))
	if err != nil {
		runNumber = 1
	}

	runAttempt, err := strconv.Atoi(os.Getenv(EnvVarRunAttempt))
	if err != nil {
		runAttempt = 1
	}

	runnerName := os.Getenv(EnvVarRunnerName)
	if runnerName == "" {
		runnerName = u.Name
	}

	runnerTemp := os.Getenv(EnvVarRunnerTemp)
	if runnerTemp == "" {
		runnerTemp = os.TempDir()
	}

	runnerToolCache := os.Getenv(EnvVarRunnerToolCache)
	if runnerToolCache == "" {
		runnerToolCache = os.TempDir()
	}

	return &GlobalContext{
		GitHubContext: &GitHubContext{
			Action:          os.Getenv(EnvVarAction),
			ActionPath:      os.Getenv(EnvVarActionPath),
			Actor:           actor,
			BaseRef:         os.Getenv(EnvVarBaseRef),
			EventName:       os.Getenv(EnvVarEventName),
			EventPath:       os.Getenv(EnvVarEventPath),
			HeadRef:         os.Getenv(EnvVarHeadRef),
			Job:             os.Getenv(EnvVarJob),
			Ref:             os.Getenv(EnvVarRef),
			RefName:         os.Getenv(EnvVarRefName),
			RefProtected:    refProtected,
			RefType:         RefType(os.Getenv(EnvVarRefType)).String(),
			Repository:      os.Getenv(EnvVarRepository),
			RepositoryOwner: os.Getenv(EnvVarRepositoryOwner),
			RunId:           os.Getenv(EnvVarRunID),
			RunNumber:       int64(runNumber),
			RunAttempt:      int64(runAttempt),
			ServerUrl:       serverURL.String(),
			Sha:             os.Getenv(EnvVarSha),
			Token:           os.Getenv(EnvVarToken),
			Workflow:        os.Getenv(EnvVarWorkflow),
			Workspace:       os.Getenv(EnvVarWorkspace),
		},
		EnvContext:   make(map[string]string),
		JobContext:   &JobContext{},
		StepsContext: make(map[string]*StepContext),
		RunnerContext: &RunnerContext{
			Name:      runnerName,
			Os:        runnerOS.String(),
			Arch:      runnerArch.String(),
			Temp:      runnerTemp,
			ToolCache: runnerToolCache,
		},
		InputsContext:  make(map[string]string),
		SecretsContext: make(map[string]string),
		NeedsContext:   make(map[string]*NeedContext),
	}
}

var (
	DefaultBranch = "main"
	DefaultRemote = "origin"
)

func NewGlobalContextFromPath(ctx context.Context, path string) (*GlobalContext, error) {
	var (
		c             = NewGlobalContextFromEnv()
		currentBranch = DefaultBranch
		currentRemote = DefaultRemote
	)

	r, err := git.PlainOpen(path)
	if err != nil {
		return nil, err
	}

	if ref, err := r.Head(); err == nil {
		c.GetGitHubContext().Sha = ref.Hash().String()
		if shaBranch := strings.Split(ref.String(), " "); len(shaBranch) > 1 {
			c.GitHubContext.RefName = strings.TrimPrefix(shaBranch[1], "refs/heads/")
			c.GitHubContext.Ref = shaBranch[1]
		} else {
			c.GitHubContext.RefName = ref.Hash().String()
			c.GitHubContext.Ref = ref.Hash().String()
		}

		if ref.Name().IsBranch() {
			currentBranch = ref.Name().Short()
			c.GitHubContext.RefType = RefTypeBranch.String()
		} else {
			c.GitHubContext.RefType = RefTypeTag.String()
		}
	}

	if conf, err := r.Config(); err == nil {
		c.GitHubContext.Actor = js.Coalesce(
			conf.User.Name,
			conf.Author.Name,
			conf.Committer.Name,
			conf.User.Email,
			conf.Author.Email,
			conf.Committer.Email,
			c.GetGitHubContext().GetActor(),
		)

		for _, remote := range conf.Remotes {
			for _, rurl := range remote.URLs {
				prurl, err := url.Parse(rurl)
				if err == nil {
					c.GetGitHubContext().Repository = strings.TrimSuffix(
						strings.TrimPrefix(prurl.Path, "/"),
						".git",
					)
					c.GetGitHubContext().RepositoryOwner = strings.Split(c.GetGitHubContext().Repository, "/")[0]
					break
				}
			}
		}
	}

	if branch, err := r.Branch(currentBranch); err == nil {
		currentRemote = branch.Remote
		c.GitHubContext.RefName = branch.Name
		c.GitHubContext.Ref = "refs/heads/" + branch.Name
		c.GitHubContext.RefType = RefTypeBranch.String()
	}

	if remote, err := r.Remote(currentRemote); err == nil {
		for _, u := range remote.Config().URLs {
			_, err := url.Parse(u)
			if err == nil {
				// TODO override default github urls
				break
			}
		}
	}

	return c, nil
}

func (c *GlobalContext) envMap() map[string]string {
	serverURL, _ := url.Parse(c.GetGitHubContext().GetServerUrl())
	apiURL, _ := APIURLFromBaseURL(serverURL)
	graphqlURL, _ := GraphQLURLFromBaseURL(serverURL)
	env := map[string]string{
		EnvVarCI:              fmt.Sprint(true),
		EnvVarWorkflow:        c.GetGitHubContext().GetWorkflow(),
		EnvVarRunID:           c.GetGitHubContext().GetRunId(),
		EnvVarRunNumber:       fmt.Sprint(c.GetGitHubContext().GetRunNumber()),
		EnvVarRunAttempt:      fmt.Sprint(c.GetGitHubContext().GetRunAttempt()),
		EnvVarJob:             c.GetGitHubContext().GetJob(),
		EnvVarAction:          c.GetGitHubContext().GetAction(),
		EnvVarActionPath:      c.GetGitHubContext().GetActionPath(),
		EnvVarActions:         fmt.Sprint(true),
		EnvVarActor:           c.GetGitHubContext().GetActor(),
		EnvVarRepository:      c.GetGitHubContext().GetRepository(),
		EnvVarEventName:       c.GetGitHubContext().GetEventName(),
		EnvVarEventPath:       c.GetGitHubContext().GetEventPath(),
		EnvVarWorkspace:       c.GetGitHubContext().GetWorkspace(),
		EnvVarSha:             c.GetGitHubContext().GetSha(),
		EnvVarRef:             c.GetGitHubContext().GetRef(),
		EnvVarRefName:         c.GetGitHubContext().GetRefName(),
		EnvVarRefProtected:    fmt.Sprint(c.GetGitHubContext().GetRefProtected()),
		EnvVarRefType:         c.GetGitHubContext().GetRefType(),
		EnvVarHeadRef:         c.GetGitHubContext().GetHeadRef(),
		EnvVarBaseRef:         c.GetGitHubContext().GetBaseRef(),
		EnvVarServerURL:       c.GetGitHubContext().GetServerUrl(),
		EnvVarAPIURL:          apiURL.String(),
		EnvVarGraphQLURL:      graphqlURL.String(),
		EnvVarRunnerName:      c.GetRunnerContext().GetName(),
		EnvVarRunnerOS:        c.GetRunnerContext().GetOs(),
		EnvVarRunnerArch:      c.GetRunnerContext().GetArch(),
		EnvVarRunnerTemp:      c.GetRunnerContext().GetTemp(),
		EnvVarRunnerToolCache: c.GetRunnerContext().GetToolCache(),
		EnvVarToken:           c.GetGitHubContext().GetToken(),
		EnvVarRepositoryOwner: c.GetGitHubContext().GetRepositoryOwner(),
	}

	if c.EnvContext != nil {
		for k, v := range c.EnvContext {
			env[k] = v
		}
	}

	if c.InputsContext != nil {
		for k, v := range c.GetInputsContext() {
			if v != "" {
				env[fmt.Sprintf("INPUT_%s", strings.ReplaceAll(strings.ToUpper(k), " ", "_"))] = v
			}
		}
	}

	return env
}

func (c *GlobalContext) Env() []string {
	return envconv.MapToArr(c.envMap())
}
