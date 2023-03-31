package tarutil

type Opts struct {
	gzip bool
}

type Opt func(*Opts)

func WithGzip(o *Opts) {
	o.gzip = true
}
