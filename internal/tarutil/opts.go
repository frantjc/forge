package tarutil

type Opts struct {
	readerGzipped bool
}

type Opt func(*Opts)

func IsGzipped(o *Opts) {
	o.readerGzipped = true
}
