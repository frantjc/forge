package command

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/url"
	"os"
	"slices"
	"strconv"
	"strings"

	"dagger.io/dagger"
	"github.com/frantjc/forge/concourse"
	"github.com/frantjc/forge/githubactions"
	client "github.com/frantjc/forge/internal/client"
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
func NewForge() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "forge",
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	cmd.AddCommand(
		NewUse(),
		NewCloudBuild(),
		NewResource(concourse.MethodCheck),
		NewResource(concourse.MethodGet),
		NewResource(concourse.MethodPut),
	)

	return cmd
}

const errWhenNoExecs = "empty result reference"

// NewUse returns the command which acts as
// the entrypoint for `forge use`.
func NewUse() *cobra.Command {
	var (
		with       = map[string]string{}
		token      string
		repo       string
		stage      = stageMain
		export     bool
		slogConfig = &logutil.SlogConfig{}
		cmd        = &cobra.Command{
			Use:           "use action [-w go-version=1.24] [--pre | --post] [-r https://github.com/frantjc/forge] [-t $GH_TOKEN] [-E] [-dqv]",
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
					workspace  = dag.Host().Directory(".")
					repository = workspace
				)

				if cmd.Flag("repo").Changed {
					src, ref, found := strings.Cut(repo, "@")
					if found {
						repository = dag.Git(src).Ref(ref).Tree()
					} else {
						r, err := url.Parse(repo)
						if err != nil {
							return err
						}

						switch r.Scheme {
						case "file", "":
							repository = dag.Host().Directory(r.Path)
						default:
							repository = dag.Git(src).Head().Tree()
						}
					}
				}

				preAction := dag.Forge().Use(uses.String(), client.ForgeUseOpts{
					Workspace: workspace,
					Repo:      repository,
					With:      envconv.MapToArr(with),
					Token:     dag.SetSecret("token", token),
					Debug:     debug,
					Env:       os.Environ(),
				})

				finalize := func() error {
					return nil
				}

				if export {
					finalize = func() error {
						// This is the same as action.Workspace() and postAction.Workspace().
						if _, err := preAction.Workspace().Export(ctx, "."); err != nil {
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

// NewCloudBuild returns the command which acts as
// the entrypoint for `forge cloudbuild`.
func NewCloudBuild() *cobra.Command {
	var (
		script                   string
		entrypoint               string
		userDefinedSubstitutions = map[string]string{}
		automapSubstituations    bool
		dynamicSubstituations    bool
		gcloudConfig             string
		export                   bool
		slogConfig               = &logutil.SlogConfig{}
		cmd                      = &cobra.Command{
			Use:           "cloudbuild cloudbuilder [-S script.sh | -E entrypoint.sh | arg...] [-s user_defined=substitution] [-AD] [-E] [-dqv] [-c ~/.gcloud/config]",
			Aliases:       []string{"cb"},
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

				dag, err := client.Connect(ctx, opts...)
				if err != nil {
					return err
				}
				defer dag.Close()

				workdir := dag.Host().Directory(".")
				gc := dag.Directory()

				if _, err := os.Stat(gcloudConfig); errors.Is(err, os.ErrNotExist) {
					if cmd.Flag("gcloud-config").Changed {
						return err
					}
				} else if err != nil {
					return err
				} else {
					gc = dag.Host().Directory(gcloudConfig)
				}


				
				cloudbuild := dag.Forge().CloudBuild(args[0], client.ForgeCloudBuildOpts{
					Workdir:      workdir,
					Entrypoint:   slices.DeleteFunc(strings.Split(entrypoint, " "), func(s string) bool {
						return s == ""
					}),
					Args:         args[1:],
					Env:          os.Environ(),
					GcloudConfig: gc,
					Script: dag.File("script", script, client.FileOpts{
						Permissions: 700,
					}),
					// TODO(frantjc): Get additional substitutions from gcloud-config.
					Substitutions:        envconv.MapToArr(userDefinedSubstitutions),
					DynamicSubstitutions: dynamicSubstituations,
					AutomapSubstitutions: automapSubstituations,
				})

				logs, err := cloudbuild.CombinedOutput(ctx)
				if err != nil {
					return err
				}

				if _, err := fmt.Fprint(cmd.OutOrStdout(), logs); err != nil {
					return err
				}

				if export {
					if _, err := cloudbuild.Workdir().Export(ctx, "."); err != nil {
						return err
					}
				}

				return nil
			},
		}
	)

	cmd.Flags().StringVarP(&entrypoint, "entrypoint", "E", "", "Entrypoint to execute")
	cmd.Flags().StringVarP(&script, "script", "S", "", "Script to run")

	cmd.Flags().StringToStringVarP(&userDefinedSubstitutions, "substitution", "s", nil, "Substitutions")
	cmd.Flags().BoolVarP(&automapSubstituations, "automap-substitutions", "A", false, "Automap substitutions")
	cmd.Flags().BoolVarP(&dynamicSubstituations, "dynamic-substitutions", "D", false, "Dynamic substitutions")

	cmd.Flags().StringVarP(&gcloudConfig, "gcloud-config", "c", "~/.gcloud/config", "GCloud config directory")

	cmd.Flags().BoolVarP(&export, "export", "e", false, "Apply changes that the action made to your working directory")

	slogConfig.AddFlags(cmd.Flags())

	cmd.MarkFlagsMutuallyExclusive("entrypoint", "script")

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
			Use:           fmt.Sprintf("%s resource [-p pipeline.yml] [-E] [-dqv]", method),
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

				return nil
			},
		}
	)

	cmd.Flags().BoolVarP(&export, "export", "e", false, "Apply changes that the action made to your working directory")

	cmd.Flags().StringVarP(&pipeline, "pipeline", "p", ".forge.yml", "Pipeline")

	switch method {
	case concourse.MethodCheck:
		cmd.Flags().StringToStringVarP(&version, "version", "V", nil, "Version")
	case concourse.MethodGet:
		cmd.Flags().StringToStringVarP(&version, "version", "V", nil, "Version")
		cmd.Flags().StringToStringVarP(&param, "param", "P", nil, "Params")
	case concourse.MethodPut:
		cmd.Flags().StringToStringVarP(&param, "param", "P", nil, "Params")
	}

	slogConfig.AddFlags(cmd.Flags())

	return cmd
}
