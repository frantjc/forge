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
	"github.com/google/uuid"
)

type globalContextKey struct{}

func WithGlobalContext(ctx context.Context, globalContext *GlobalContext) context.Context {
	return context.WithValue(ctx, globalContextKey{}, globalContext)
}

func GlobalContextFrom(ctx context.Context) (globalContext *GlobalContext, ok bool) {
	globalContext, ok = ctx.Value(globalContextKey{}).(*GlobalContext)
	return
}

type GlobalContext struct {
	GitHubContext  *GitHubContext
	EnvContext     map[string]string
	JobContext     *JobContext
	StepsContext   map[string]*StepsContext
	RunnerContext  *RunnerContext
	InputsContext  map[string]string
	SecretsContext map[string]string
	NeedsContext   map[string]*NeedsContext
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

// GitHubContext represents the GitHub Context
// https://docs.github.com/en/actions/learn-github-actions/Context#github-context
type GitHubContext struct {
	Action          string
	ActionPath      string
	Actor           string
	BaseRef         string
	Event           string
	EventName       string
	EventPath       string
	HeadRef         string
	Job             string
	Ref             string
	RefName         string
	RefProtected    bool
	RefType         RefType
	Repository      string
	RepositoryOwner string
	RunID           string
	RunNumber       int
	RunAttempt      int
	ServerURL       *url.URL
	Sha             string
	Token           string
	Workflow        string
	Workspace       string
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
			return c.RefType.String()
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
			return c.ServerURL.String()
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

type JobContext struct {
	Container *struct {
		ID      string
		Network string
	}
	Services map[string]struct {
		ID      string
		Network string
		// not sure if this is the correct representation
		Ports map[string]string
	}
	Status string
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

type StepsContext struct {
	Outputs    map[string]string
	Conclusion string
	Outcome    string
}

func (c *StepsContext) GetString(key string) string {
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

type RunnerContext struct {
	Name      string
	OS        OS
	Arch      Arch
	Temp      string
	ToolCache string
}

func (c *RunnerContext) GetString(key string) string {
	keys := strings.Split(key, ".")
	if len(keys) > 0 {
		switch keys[0] {
		case "name":
			return c.Name
		case "os":
			return c.OS.String()
		case "arch":
			return c.Arch.String()
		case "temp":
			return c.Temp
		case "tool_cache":
			return c.ToolCache
		}
	}

	return ""
}

type NeedsContext struct {
	Outputs map[string]string
}

func (c *NeedsContext) GetString(key string) string {
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

	serverURL, err := url.Parse(os.Getenv(EnvVarServerURL))
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

	runID := os.Getenv(EnvVarRunID)
	if runID == "" {
		runID = uuid.NewString()
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
			RefType:         RefType(os.Getenv(EnvVarRefType)),
			Repository:      os.Getenv(EnvVarRepository),
			RepositoryOwner: os.Getenv(EnvVarRepositoryOwner),
			RunID:           runID,
			RunNumber:       runNumber,
			RunAttempt:      runAttempt,
			ServerURL:       serverURL,
			Sha:             os.Getenv(EnvVarSha),
			Token:           os.Getenv(EnvVarToken),
			Workflow:        os.Getenv(EnvVarWorkflow),
			Workspace:       os.Getenv(EnvVarWorkspace),
		},
		EnvContext:   make(map[string]string),
		JobContext:   &JobContext{},
		StepsContext: make(map[string]*StepsContext),
		RunnerContext: &RunnerContext{
			Name:      runnerName,
			OS:        runnerOS,
			Arch:      runnerArch,
			Temp:      runnerTemp,
			ToolCache: runnerToolCache,
		},
		InputsContext:  make(map[string]string),
		SecretsContext: make(map[string]string),
		NeedsContext:   make(map[string]*NeedsContext),
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
		c.GitHubContext.Actor = js.Coalesce(
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
	apiURL, _ := APIURLFromBaseURL(c.GitHubContext.ServerURL)
	graphqlURL, _ := GraphQLURLFromBaseURL(c.GitHubContext.ServerURL)
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
		EnvVarRefType:         c.GitHubContext.RefType.String(),
		EnvVarHeadRef:         c.GitHubContext.HeadRef,
		EnvVarBaseRef:         c.GitHubContext.BaseRef,
		EnvVarServerURL:       c.GitHubContext.ServerURL.String(),
		EnvVarAPIURL:          apiURL.String(),
		EnvVarGraphQLURL:      graphqlURL.String(),
		EnvVarRunnerName:      c.RunnerContext.Name,
		EnvVarRunnerOS:        c.RunnerContext.OS.String(),
		EnvVarRunnerArch:      c.RunnerContext.Arch.String(),
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
