package forge

import (
	"context"

	"github.com/google/uuid"
)

func runOptsWithDefaults(opts ...RunOpt) *RunOpts {
	o := &RunOpts{
		WorkingDir: "/" + uuid.NewString(),
	}

	for _, opt := range opts {
		opt.Apply(o)
	}

	return o
}

type RunOpts struct {
	Streams             *Streams
	Mounts              []Mount
	InterceptDockerSock bool
	WorkingDir          string
}

func (o *RunOpts) Apply(opts *RunOpts) {
	if opts == nil {
		opts = &RunOpts{}
	}
	if o.Streams != nil {
		opts.Streams = o.Streams
	}
	opts.Mounts = overrideMounts(opts.Mounts, o.Mounts...)
	if o.InterceptDockerSock {
		opts.InterceptDockerSock = true
	}
	if o.WorkingDir != "" {
		opts.WorkingDir = o.WorkingDir
	}
}

func WithStreams(streams *Streams) RunOpt {
	return &RunOpts{
		Streams: streams,
	}
}

func WithStdStreams() RunOpt {
	return WithStreams(StdStreams())
}

type RunOpt interface {
	Apply(*RunOpts)
}

type Runnable interface {
	Run(context.Context, ContainerRuntime, ...RunOpt) error
}
