package tarutil

type Opts struct {
	gzipped bool
}

type Opt func(*Opts)

func IsGzipped(o *Opts) {
	o.gzipped = true
}
