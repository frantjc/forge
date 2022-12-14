package githubactions

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/frantjc/forge"
	"github.com/frantjc/forge/envconv"
	"github.com/frantjc/go-fn"
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
				return c.GitHubContext.GetString(strings.Join(keys[1:], "."))
			}
		case "env":
			if len(keys) > 1 {
				if v, ok := c.EnvContext[keys[1]]; ok {
					return v
				}
			}
		case "job":
			if len(keys) > 1 {
				return c.JobContext.GetString(strings.Join(keys[1:], "."))
			}
		case "steps":
			if len(keys) > 2 {
				if v, ok := c.StepsContext[keys[1]]; ok {
					return v.GetString(strings.Join(keys[2:], "."))
				}
			}
		case "runner":
			if len(keys) > 1 {
				return c.RunnerContext.GetString(strings.Join(keys[1:], "."))
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
			return c.Action
		case "action_path":
			return c.ActionPath
		case "actor":
			return c.Actor
		case "base_ref":
			return c.BaseRef
		case "event":
			return c.Event
		case "event_name":
			return c.EventName
		case "event_path":
			return c.EventPath
		case "head_ref":
			return c.HeadRef
		case "job":
			return c.Job
		case "ref":
			return c.Ref
		case "ref_name":
			return c.RefName
		case "ref_protected":
			return fmt.Sprint(c.RefProtected)
		case "ref_type":
			return c.RefType
		case "repository":
			return c.Repository
		case "repository_owner":
			return c.RepositoryOwner
		case "run_id":
			return c.RunID
		case "run_number":
			return fmt.Sprint(c.RunNumber)
		case "run_attempt":
			return fmt.Sprint(c.RunAttempt)
		case "server_url":
			return c.ServerURL
		case "sha":
			return c.Sha
		case "token":
			return c.Token
		case "workflow":
			return c.Workflow
		case "workspace":
			return c.Workspace
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
					return c.Container.ID
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
							return v.ID
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
			return c.Outcome
		case "conclusion":
			return c.Conclusion
		}
	}

	return ""
}

func (c *RunnerContext) GetString(key string) string {
	keys := strings.Split(key, ".")
	if len(keys) > 0 {
		switch keys[0] {
		case "name":
			return c.Name
		case "os":
			return c.OS
		case "arch":
			return c.Arch
		case "temp":
			return c.Temp
		case "tool_cache":
			return c.ToolCache
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

	runnerOS := os.Getenv(EnvVarRunnerOS)
	if runnerOS == "" {
		runnerOS = OSLinux
	}

	runnerArch := os.Getenv(EnvVarRunnerArch)
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
			RefType:         os.Getenv(EnvVarRefType),
			Repository:      os.Getenv(EnvVarRepository),
			RepositoryOwner: os.Getenv(EnvVarRepositoryOwner),
			RunID:           os.Getenv(EnvVarRunID),
			RunNumber:       int64(runNumber),
			RunAttempt:      int64(runAttempt),
			ServerURL:       serverURL.String(),
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
			OS:        runnerOS,
			Arch:      runnerArch,
			Temp:      runnerTemp,
			ToolCache: runnerToolCache,
		},
		InputsContext:  make(map[string]string),
		SecretsContext: make(map[string]string),
		NeedsContext:   make(map[string]*NeedContext),
	}
}

const (
	DefaultBranch = "main"
	DefaultRemote = "origin"
)

func NewGlobalContextFromPath(ctx context.Context, path string) (*GlobalContext, error) {
	var (
		_             = forge.LoggerFrom(ctx)
		c             = NewGlobalContextFromEnv()
		currentBranch = DefaultBranch
		currentRemote = DefaultRemote
	)

	r, err := git.PlainOpen(path)
	for ; err != nil && path != "/"; r, err = git.PlainOpen(path) {
		path = filepath.Dir(path)
	}
	if err != nil {
		return nil, err
	}

	if ref, err := r.Head(); err == nil {
		c.GitHubContext.Sha = ref.Hash().String()
		if shaBranch := strings.Split(ref.String(), " "); len(shaBranch) > 1 {
			c.GitHubContext.RefName = strings.TrimPrefix(shaBranch[1], "refs/heads/")
			c.GitHubContext.Ref = shaBranch[1]
		} else {
			c.GitHubContext.RefName = ref.Hash().String()
			c.GitHubContext.Ref = ref.Hash().String()
		}

		if ref.Name().IsBranch() {
			currentBranch = ref.Name().Short()
			c.GitHubContext.RefType = RefTypeBranch
		} else {
			c.GitHubContext.RefType = RefTypeTag
		}
	}

	if conf, err := r.Config(); err == nil {
		c.GitHubContext.Actor = fn.Coalesce(
			conf.User.Name,
			conf.Author.Name,
			conf.Committer.Name,
			conf.User.Email,
			conf.Author.Email,
			conf.Committer.Email,
			c.GitHubContext.Actor,
		)

		for _, remote := range conf.Remotes {
			for _, rurl := range remote.URLs {
				prurl, err := url.Parse(rurl)
				if err == nil {
					c.GitHubContext.Repository = strings.TrimSuffix(
						strings.TrimPrefix(prurl.Path, "/"),
						".git",
					)
					c.GitHubContext.RepositoryOwner = strings.Split(c.GitHubContext.Repository, "/")[0]
					break
				}
			}
		}
	}

	if branch, err := r.Branch(currentBranch); err == nil {
		currentRemote = branch.Remote
		c.GitHubContext.RefName = branch.Name
		c.GitHubContext.Ref = "refs/heads/" + branch.Name
		c.GitHubContext.RefType = RefTypeBranch
	}

	if remote, err := r.Remote(currentRemote); err == nil {
		for _, u := range remote.Config().URLs {
			if _, err := url.Parse(u); err == nil {
				// TODO override default github urls
				break
			}
		}
	}

	return c, nil
}

func (c *GlobalContext) envMap() map[string]string {
	serverURL, _ := url.Parse(c.GitHubContext.ServerURL)
	apiURL, _ := APIURLFromBaseURL(serverURL)
	graphqlURL, _ := GraphQLURLFromBaseURL(serverURL)
	env := map[string]string{
		EnvVarCI:              fmt.Sprint(true),
		EnvVarWorkflow:        c.GitHubContext.Workflow,
		EnvVarRunID:           c.GitHubContext.RunID,
		EnvVarRunNumber:       fmt.Sprint(c.GitHubContext.RunNumber),
		EnvVarRunAttempt:      fmt.Sprint(c.GitHubContext.RunAttempt),
		EnvVarJob:             c.GitHubContext.Job,
		EnvVarAction:          c.GitHubContext.Action,
		EnvVarActionPath:      c.GitHubContext.ActionPath,
		EnvVarActions:         fmt.Sprint(true),
		EnvVarActor:           c.GitHubContext.Actor,
		EnvVarRepository:      c.GitHubContext.Repository,
		EnvVarEventName:       c.GitHubContext.EventName,
		EnvVarEventPath:       c.GitHubContext.EventPath,
		EnvVarWorkspace:       c.GitHubContext.Workspace,
		EnvVarSha:             c.GitHubContext.Sha,
		EnvVarRef:             c.GitHubContext.Ref,
		EnvVarRefName:         c.GitHubContext.RefName,
		EnvVarRefProtected:    fmt.Sprint(c.GitHubContext.RefProtected),
		EnvVarRefType:         c.GitHubContext.RefType,
		EnvVarHeadRef:         c.GitHubContext.HeadRef,
		EnvVarBaseRef:         c.GitHubContext.BaseRef,
		EnvVarServerURL:       c.GitHubContext.ServerURL,
		EnvVarAPIURL:          apiURL.String(),
		EnvVarGraphQLURL:      graphqlURL.String(),
		EnvVarRunnerName:      c.RunnerContext.Name,
		EnvVarRunnerOS:        c.RunnerContext.OS,
		EnvVarRunnerArch:      c.RunnerContext.Arch,
		EnvVarRunnerTemp:      c.RunnerContext.Temp,
		EnvVarRunnerToolCache: c.RunnerContext.ToolCache,
		EnvVarToken:           c.GitHubContext.Token,
		EnvVarRepositoryOwner: c.GitHubContext.RepositoryOwner,
	}

	if c.EnvContext != nil {
		for k, v := range c.EnvContext {
			env[k] = v
		}
	}

	if c.InputsContext != nil {
		for k, v := range c.InputsContext {
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

type GlobalContext struct {
	GitHubContext  *GitHubContext          `json:"git_hub_context,omitempty"`
	EnvContext     map[string]string       `json:"env_context,omitempty"`
	JobContext     *JobContext             `json:"job_context,omitempty"`
	StepsContext   map[string]*StepContext `json:"steps_context,omitempty"`
	RunnerContext  *RunnerContext          `json:"runner_context,omitempty"`
	InputsContext  map[string]string       `json:"inputs_context,omitempty"`
	SecretsContext map[string]string       `json:"secrets_context,omitempty"`
	NeedsContext   map[string]*NeedContext `json:"needs_context,omitempty"`
}

type GitHubContext struct {
	Action          string `json:"action,omitempty"`
	ActionPath      string `json:"action_path,omitempty"`
	Actor           string `json:"actor,omitempty"`
	BaseRef         string `json:"base_ref,omitempty"`
	Event           string `json:"event,omitempty"`
	EventName       string `json:"event_name,omitempty"`
	EventPath       string `json:"event_path,omitempty"`
	HeadRef         string `json:"head_ref,omitempty"`
	Job             string `json:"job,omitempty"`
	Ref             string `json:"ref,omitempty"`
	RefName         string `json:"ref_name,omitempty"`
	RefProtected    bool   `json:"ref_protected,omitempty"`
	RefType         string `json:"ref_type,omitempty"`
	Repository      string `json:"repository,omitempty"`
	RepositoryOwner string `json:"repository_owner,omitempty"`
	RunID           string `json:"run_id,omitempty"`
	RunNumber       int64  `json:"run_number,omitempty"`
	RunAttempt      int64  `json:"run_attempt,omitempty"`
	ServerURL       string `json:"server_url,omitempty"`
	Sha             string `json:"sha,omitempty"`
	Token           string `json:"token,omitempty"`
	Workflow        string `json:"workflow,omitempty"`
	Workspace       string `json:"workspace,omitempty"`
}

type JobContext struct {
	Container *JobContextContainer          `json:"container,omitempty"`
	Services  map[string]*JobContextService `json:"services,omitempty"`
	Status    string                        `json:"status,omitempty"`
}

type StepContext struct {
	Outputs    map[string]string `json:"outputs,omitempty"`
	Conclusion string            `json:"conclusion,omitempty"`
	Outcome    string            `json:"outcome,omitempty"`
}

type RunnerContext struct {
	Name      string `json:"name,omitempty"`
	OS        string `json:"os,omitempty"`
	Arch      string `json:"arch,omitempty"`
	Temp      string `json:"temp,omitempty"`
	ToolCache string `json:"tool_cache,omitempty"`
}

type NeedContext struct {
	Outputs map[string]string `json:"outputs,omitempty"`
}

type JobContextContainer struct {
	ID      string `json:"id,omitempty"`
	Network string `json:"network,omitempty"`
}

type JobContextService struct {
	ID      string            `json:"id,omitempty"`
	Network string            `json:"network,omitempty"`
	Ports   map[string]string `json:"ports,omitempty"`
}
