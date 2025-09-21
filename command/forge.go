package command

import (
	"fmt"
	"io"
	"log/slog"
	"net/url"
	"os"
	"strconv"
	"strings"

	"dagger.io/dagger"
	"github.com/frantjc/forge/concourse"
	"github.com/frantjc/forge/githubactions"
	"github.com/frantjc/forge/internal/client"
	"github.com/frantjc/forge/internal/envconv"
	"github.com/frantjc/forge/internal/logutil"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	stagePre = iota
	stageMain
	stagePost
)

type genericBool[T any] struct {
	Value *T
	IfSet T
}

var _ pflag.Value = new(genericBool[any])

// Set implements pflag.Value.
func (b *genericBool[T]) Set(s string) error {
	v, err := strconv.ParseBool(s)
	if v {
		*b.Value = b.IfSet
	}
	return err
}

// String implements pflag.Value.
func (b *genericBool[T]) String() string {
	return fmt.Sprint(b.Value)
}

// Type implements pflag.Value.
func (b *genericBool[T]) Type() string {
	return "bool"
}

// Type implements pflag.boolFlag.
func (b *genericBool[T]) IsBoolFlag() bool {
	return true
}

// NewForge returns the command which acts as
// the entrypoint for `forge use`.
func NewForge(version string) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "forge",
		SilenceErrors: true,
		SilenceUsage:  true,
		Version:       version,
	}

	cmd.AddCommand(
		NewUse(),
		NewResource(concourse.MethodCheck),
		NewResource(concourse.MethodGet),
		NewResource(concourse.MethodPut),
	)

	cmd.Flags().BoolP("help", "h", false, "Help for "+cmd.Name())
	cmd.Flags().Bool("version", false, "Version for "+cmd.Name())

	return cmd
}

const errWhenNoExecs = "no command has been set"

// NewUse returns the command which acts as
// the entrypoint for `forge use`.
func NewUse() *cobra.Command {
	var (
		with       = map[string]string{}
		env        = []string{}
		token      string
		repo       string
		stage      = stageMain
		export     bool
		slogConfig = &logutil.SlogConfig{}
		cmd        = &cobra.Command{
			Use:           "use action [-w go-version=1.24] [--env KEY[=VALUE] [--pre | --post] [-r https://github.com/frantjc/forge] [-t $GH_TOKEN] [-e] [-dqv]",
			Aliases:       []string{"u", "uses"},
			SilenceErrors: true,
			SilenceUsage:  true,
			Args:          cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				var (
					log   = slog.New(slog.NewTextHandler(cmd.OutOrStdout(), &slog.HandlerOptions{Level: slogConfig}))
					ctx   = logutil.SloggerInto(cmd.Context(), log)
					debug = log.Enabled(ctx, slog.LevelDebug)
				)

				uses, err := githubactions.Parse(args[0])
				if err != nil {
					return err
				}

				if !cmd.Flag("token").Changed {
					token = os.Getenv(githubactions.EnvVarToken)
				}

				opts := []dagger.ClientOpt{
					dagger.WithLogOutput(io.Discard),
				}

				if debug {
					opts = []dagger.ClientOpt{
						dagger.WithLogOutput(cmd.ErrOrStderr()),
						dagger.WithVerbosity(int(slogConfig.Level())),
					}
				}

				dag, err := client.Connect(ctx, opts...)
				if err != nil {
					return err
				}
				defer dag.Close()

				var (
					workspace = dag.Host().Directory(".")
					ref       = workspace.AsGit().Head()
				)

				if cmd.Flag("repo").Changed {
					src, ver, found := strings.Cut(repo, "@")
					if found {
						ref = dag.Git(src).Ref(ver)
					} else {
						r, err := url.Parse(repo)
						if err != nil {
							return err
						}

						switch r.Scheme {
						case "file", "":
							ref = dag.Host().Directory(r.Path).AsGit().Head()
						default:
							ref = dag.Git(src).Head()
						}
					}
				}

				preAction := dag.
					Forge().
					Use(uses.String(), client.ForgeUseOpts{Workspace: workspace}).
					WithToken(dag.SetSecret("token", token)).
					WithRef(ref)

				if debug {
					preAction = preAction.WithDebug()
				}

				for _, e := range env {
					if k, v, found := strings.Cut(e, "="); !found {
						preAction.WithEnv(k, os.Getenv(k))
					} else {
						preAction.WithEnv(k, v)
					}
				}

				for k, v := range with {
					preAction = preAction.WithInput(k, v)
				}

				finalize := func() error {
					return nil
				}

				if export {
					finalize = func() error {
						// This is the same as action.Workspace() and postAction.Workspace().
						if _, err := preAction.Workspace().Changes(workspace).Export(ctx, "."); err != nil {
							return err
						}

						return nil
					}
				}

				action := preAction.Pre()

				logs, err := action.Container().CombinedOutput(ctx)
				if err != nil {
					if err.Error() != errWhenNoExecs {
						return err
					}
				}

				if _, err := fmt.Fprint(cmd.OutOrStdout(), logs); err != nil {
					return err
				}

				if stage < stageMain {
					return finalize()
				}

				postAction := action.Main()

				logs, err = postAction.Container().CombinedOutput(ctx)
				if err != nil {
					if err.Error() != errWhenNoExecs {
						return err
					}
				}

				if _, err := fmt.Fprint(cmd.OutOrStdout(), logs); err != nil {
					return err
				}

				if stage < stagePost {
					return finalize()
				}

				logs, err = postAction.Post().CombinedOutput(ctx)
				if err != nil {
					if err.Error() != errWhenNoExecs {
						return err
					}
				}

				if _, err := fmt.Fprint(cmd.OutOrStdout(), logs); err != nil {
					return err
				}

				return finalize()
			},
		}
	)

	cmd.Flags().StringToStringVarP(&with, "with", "w", nil, "With params")
	cmd.Flags().StringArrayVar(&env, "env", nil, "Env")
	cmd.Flags().StringVarP(&token, "token", "t", "", "GitHub token")
	cmd.Flags().StringVarP(&repo, "repo", "r", "", "Git repository to gather context from")

	cmd.Flags().BoolVarP(&export, "export", "e", false, "Apply changes that the action made to your workspace")

	slogConfig.AddFlags(cmd.Flags())

	cmd.Flags().AddFlag(&pflag.Flag{
		Name: "pre",
		Value: &genericBool[int]{
			Value: &stage,
			IfSet: stagePre,
		},
		NoOptDefVal: "true",
		Usage:       "Run only the pre-action step",
	})
	cmd.Flags().AddFlag(&pflag.Flag{
		Name: "post",
		Value: &genericBool[int]{
			Value: &stage,
			IfSet: stagePost,
		},
		NoOptDefVal: "true",
		Usage:       "Run the post-action step",
	})

	cmd.MarkFlagsMutuallyExclusive("pre", "post")

	return cmd
}

// NewResource returns the command which acts as
// the entrypoint for `forge check`, `forge get`, and `forge put`.
func NewResource(method string) *cobra.Command {
	var (
		version    map[string]string
		param      map[string]string
		pipeline   string
		export     bool
		slogConfig = &logutil.SlogConfig{}
		cmd        = &cobra.Command{
			Use:           fmt.Sprintf("%s resource [-p pipeline.yml] [-e] [-dqv]", method),
			SilenceErrors: true,
			SilenceUsage:  true,
			Args:          cobra.MinimumNArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				var (
					log   = slog.New(slog.NewTextHandler(cmd.OutOrStdout(), &slog.HandlerOptions{Level: slogConfig}))
					ctx   = logutil.SloggerInto(cmd.Context(), log)
					debug = log.Enabled(ctx, slog.LevelDebug)
				)

				opts := []dagger.ClientOpt{
					dagger.WithLogOutput(io.Discard),
				}

				if debug {
					opts = []dagger.ClientOpt{
						dagger.WithLogOutput(cmd.ErrOrStderr()),
						dagger.WithVerbosity(int(slogConfig.Level())),
					}
				}

				f, err := os.Open(pipeline)
				if err != nil {
					return err
				}
				defer f.Close()

				b, err := io.ReadAll(f)
				if err != nil {
					return err
				}

				dag, err := client.Connect(ctx, opts...)
				if err != nil {
					return err
				}
				defer dag.Close()

				workdir := dag.Host().Directory(".")

				resource := dag.Forge().Resource(args[0], client.ForgeResourceOpts{
					Pipeline: dag.File(pipeline, string(b)),
					Workdir:  workdir,
				})

				switch method {
				case concourse.MethodCheck:
					logs, err := resource.Check(client.ForgeResourceCheckOpts{
						Version: envconv.MapToArr(version),
					}).
						Container().
						CombinedOutput(ctx)
					if err != nil {
						return err
					}

					if _, err := fmt.Fprint(cmd.OutOrStdout(), logs); err != nil {
						return err
					}
				case concourse.MethodGet:
					logs, err := resource.Get(client.ForgeResourceGetOpts{
						Version: envconv.MapToArr(version),
						Param:   envconv.MapToArr(param),
					}).
						Container().
						CombinedOutput(ctx)
					if err != nil {
						return err
					}

					if _, err := fmt.Fprint(cmd.OutOrStdout(), logs); err != nil {
						return err
					}
				case concourse.MethodPut:
					logs, err := resource.Put(client.ForgeResourcePutOpts{
						Param: envconv.MapToArr(param),
					}).
						Container().
						CombinedOutput(ctx)
					if err != nil {
						return err
					}

					if _, err := fmt.Fprint(cmd.OutOrStdout(), logs); err != nil {
						return err
					}
				default:
					return fmt.Errorf("unknown resource method %s", method)
				}

				if export {
					if _, err := resource.Workdir().Changes(workdir).Export(ctx, "."); err != nil {
						return err
					}
				}

				return nil
			},
		}
	)

	cmd.Flags().BoolVarP(&export, "export", "e", false, "Apply changes that the resource made to your working directory")

	cmd.Flags().StringVarP(&pipeline, "pipeline", "p", ".forge.yml", "Concourse pipeline file to get resources from")

	switch method {
	case concourse.MethodCheck:
		cmd.Flags().StringToStringVarP(&version, "version", "V", nil, "Concourse resource version")
	case concourse.MethodGet:
		cmd.Flags().StringToStringVarP(&version, "version", "V", nil, "Concourse resource version")
		cmd.Flags().StringToStringVarP(&param, "param", "P", nil, "Concourse resource params")
	case concourse.MethodPut:
		cmd.Flags().StringToStringVarP(&param, "param", "P", nil, "Concourse resource params")
	}

	slogConfig.AddFlags(cmd.Flags())

	return cmd
}
