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

	cmd.AddCommand(NewUse())

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
			Use:           "use",
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
					dagger.WithEnvironmentVariable(githubactions.EnvVarToken, token),
				}

				if debug {
					opts = append(opts,
						dagger.WithLogOutput(cmd.ErrOrStderr()),
						dagger.WithVerbosity(int(slogConfig.Level())),
					)
				}

				dag, err := client.Connect(ctx, opts...)
				if err != nil {
					return err
				}
				defer dag.Close()

				workspace := dag.Host().Directory(".")
				repository := workspace

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

	return cmd
}
