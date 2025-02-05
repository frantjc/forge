package forge

const (
	DefaultNode10ImageReference = "docker.io/library/node:10"
	DefaultNode12ImageReference = "docker.io/library/node:12"
	DefaultNode16ImageReference = "docker.io/library/node:16"
	DefaultNode20ImageReference = "docker.io/library/node:20"
	DefaultNodeImageReference   = DefaultNode16ImageReference
)

var (
	Node10ImageReference = DefaultNode10ImageReference
	Node12ImageReference = DefaultNode12ImageReference
	Node16ImageReference = DefaultNode16ImageReference
	Node20ImageReference = DefaultNode20ImageReference
	NodeImageReference   = DefaultNodeImageReference
)
