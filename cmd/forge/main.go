package main

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"dagger.io/dagger"
	"github.com/frantjc/forge/githubactions"
	client "github.com/frantjc/forge/internal/client"
	"github.com/frantjc/forge/internal/envconv"
	xos "github.com/frantjc/x/os"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"golang.org/x/exp/constraints"
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

type incrementalCount[T constraints.Integer] struct {
	Value     *T
	Increment T
}

var _ pflag.Value = new(incrementalCount[int])

// Set implements pflag.Value.
func (c *incrementalCount[T]) Set(s string) error {
	v, err := strconv.ParseInt(s, 0, 0)
	*c.Value += (T(v) * c.Increment)
	return err
}

// String implements pflag.Value.
func (c *incrementalCount[T]) String() string {
	return strconv.Itoa(int(*c.Value))
}

// Type implements pflag.Value.
func (c *incrementalCount[T]) Type() string {
	return "count"
}

const errWhenNoExecs = "empty result reference"

// NewForge returns the command which acts as
// the entrypoint for `forge`.
func NewForge() *cobra.Command {
	var (
		with      = map[string]string{}
		debug, _  = strconv.ParseBool(os.Getenv("DEBUG"))
		token     string
		repo      string
		stage     = stageMain
		verbosity = 0
		cmd       = &cobra.Command{
			Use:           "forge",
			SilenceErrors: true,
			SilenceUsage:  true,
			Args:          cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				ctx := cmd.Context()

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

				if verbosity > 4 {
					opts = append(opts,
						dagger.WithLogOutput(cmd.ErrOrStderr()),
						dagger.WithVerbosity(verbosity),
					)
					debug = true
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
					return nil
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
					return nil
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

				return nil
			},
		}
	)

	cmd.Flags().StringToStringVarP(&with, "with", "w", nil, "With params")
	cmd.Flags().StringVarP(&token, "token", "t", "", "GitHub token")
	cmd.Flags().StringVarP(&repo, "repo", "r", "", "Git repository to gather context from")

	cmd.Flags().BoolVarP(&debug, "debug", "d", debug, "Debug logging")
	cmd.Flags().AddFlag(&pflag.Flag{
		Name:      "quiet",
		Shorthand: "q",
		Value: &genericBool[bool]{
			Value: &debug,
			IfSet: false,
		},
		NoOptDefVal: "true",
		Usage:       "Run only the pre-action step",
	})
	cmd.Flags().AddFlag(&pflag.Flag{
		Name:      "verbose",
		Shorthand: "v",
		Value: &incrementalCount[int]{
			Value:     &verbosity,
			Increment: 1,
		},
		NoOptDefVal: "+1",
		Usage:       "More vebose logging",
	})

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

func main() {
	var (
		cmd       = NewForge()
		ctx, stop = signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	)

	err := cmd.ExecuteContext(ctx)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
	}

	stop()
	xos.ExitFromError(err)
}
