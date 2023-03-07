package ore

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/frantjc/forge"
	"github.com/frantjc/forge/circleci"
	"github.com/frantjc/forge/envconv"
	"github.com/frantjc/forge/forgecircleci"
	"github.com/frantjc/forge/internal/bin"
	"github.com/frantjc/forge/internal/containerutil"
	errorcode "github.com/frantjc/go-error-code"
	"github.com/frantjc/go-fn"
)

type Orb struct {
	Orb        string            `json:"orb,omitempty"`
	Command    string            `json:"command,omitempty"`
	Parameters map[string]string `json:"parameters,omitempty"`
}

func (o *Orb) Liquify(ctx context.Context, containerRuntime forge.ContainerRuntime, drains *forge.Drains) error {
	_ = forge.LoggerFrom(ctx)

	orb, err := circleci.Parse(o.Orb)
	if err != nil {
		return err
	}

	source, err := circleci.GetOrbSource(ctx, orb)
	if err != nil {
		return err
	}

	job := &circleci.Job{}

	command, ok := source.Commands[o.Command]
	if ok {
		job.Command = command
		for name := range source.Executors {
			job.Executor = map[string]any{
				"name": name,
			}
			break
		}
	} else {
		if djob, ok := source.Jobs[o.Command]; ok {
			job = &djob
		} else {
			return fmt.Errorf("not a job or command: %s", o.Command)
		}
	}

	executorName, ok := job.Executor["name"].(string)
	if !ok {
		return fmt.Errorf("get executor name")
	}

	executor, ok := source.Executors[executorName]
	if !ok {
		return fmt.Errorf("find executor: %s", executorName)
	}

	parameters := o.Parameters
	if parameters == nil {
		parameters = map[string]string{}
	}

	for name, parameter := range executor.Parameters {
		if value, ok := parameters[name]; !ok {
			parameters[name] = fmt.Sprint(parameter.Default)
		} else if len(parameter.Enum) > 0 && !fn.Includes(parameter.Enum, value) {
			return fmt.Errorf("invalid parameter: %s, must be one of: %s", value, strings.Join(parameter.Enum, ", "))
		}
	}

	expander := circleci.ExpandFunc(func(s string) string {
		return parameters[strings.TrimPrefix(s, "parameters.")]
	})

	if len(executor.Docker) == 0 {
		return fmt.Errorf("exector has no image: %s", executorName)
	}

	image, err := containerRuntime.PullImage(ctx, expander.ExpandString(executor.Docker[0].Image))
	if err != nil {
		return err
	}

	container, err := containerutil.CreateSleepingContainer(ctx, containerRuntime, image, &forge.ContainerConfig{
		Env: []string{
			"BASH_ENV=" + filepath.Join(forgecircleci.DefaultRootPath, "bash.env"),
			"HOME=" + forgecircleci.DefaultHome,
		},
		WorkingDir: forgecircleci.DefaultHome,
	})
	if err != nil {
		return err
	}

	if err = container.CopyTo(ctx, forgecircleci.DefaultRootPath, bin.NewTarArchiveWithEmptyFiles("bash.env")); err != nil {
		return err
	}

	for name, parameter := range command.Parameters {
		if value, ok := parameters[name]; !ok {
			parameters[name] = fmt.Sprint(parameter.Default)
		} else if len(parameter.Enum) > 0 && !fn.Includes(parameter.Enum, value) {
			return fmt.Errorf("parameter: %s, must be one of: %s", value, strings.Join(parameter.Enum, ", "))
		}
	}

	return o.runSteps(ctx, containerRuntime, container, drains.ToStreams(nil), command.Steps, parameters)
}

func (o *Orb) runSteps(ctx context.Context, containerRuntime forge.ContainerRuntime, container forge.Container, streams *forge.Streams, steps []circleci.Step, parameters map[string]string) error {
	var (
		_        = forge.LoggerFrom(ctx)
		expander = circleci.ExpandFunc(func(s string) string {
			return parameters[strings.TrimPrefix(s, "parameters.")]
		})
	)

	for _, step := range steps {
		switch {
		case step.Run != nil:
			if exitCode, err := container.Exec(ctx, &forge.ContainerConfig{
				Entrypoint: []string{"/bin/bash", "-c", step.Run.Command},
				Env: fn.Map(envconv.MapToArr(step.Run.Environment), func(s string, _ int) string {
					return expander.ExpandString(s)
				}),
				WorkingDir: forgecircleci.DefaultRootPath,
			}, streams); err != nil {
				return err
			} else if exitCode > 0 {
				return errorcode.New(ErrContainerExitedWithNonzeroExitCode, errorcode.WithExitCode(exitCode))
			}
		case step.Unless != nil:
			if !circleci.EvaluateConditional(expander, step.Unless) {
				if err := o.runSteps(ctx, containerRuntime, container, streams, step.Unless.Steps, parameters); err != nil {
					return err
				}
			}
		case step.When != nil:
			if circleci.EvaluateConditional(expander, step.When) {
				if err := o.runSteps(ctx, containerRuntime, container, streams, step.When.Steps, parameters); err != nil {
					return err
				}
			}
		default:
			switch len(step.Dynamic) {
			case 0:
				return fmt.Errorf("empty step")
			case 1:
				for command, commandParameters := range step.Dynamic {
					if command == "checkout" {
						return fmt.Errorf("checkout unimplemented")
					}

					subParameters := map[string]string{}
					for name, parameter := range commandParameters {
						subParameters[name] = expander.ExpandString(fmt.Sprint(parameter))
					}

					orb := &Orb{
						Orb:        o.Orb,
						Command:    command,
						Parameters: subParameters,
					}

					return orb.Liquify(ctx, containerRuntime, streams.Drains)
				}
			default:
				return fmt.Errorf("ambiguous step")
			}
		}
	}

	return nil
}
