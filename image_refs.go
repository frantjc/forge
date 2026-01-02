package forge


const (
	DefaultNode10ImageReference = "docker.io/library/node:10"
	// DefaultNode12ImageReference is the default container image reference that
	// is used by GitHub Actions that run using node12.
	DefaultNode12ImageReference = "docker.io/library/node:12"
	// DefaultNode16ImageReference is the default container image reference that
	// is used by GitHub Actions that run using node16.
	DefaultNode16ImageReference = "docker.io/library/node:16"
	// DefaultNode20ImageReference is the default container image reference that
	// is used by GitHub Actions that run using node20.
	DefaultNode20ImageReference = "docker.io/library/node:20"
	// DefaultNode24ImageReference is the default container image reference that
	// is used by GitHub Actions that run using node24.
	DefaultNode24ImageReference = "docker.io/library/node:24"
	DefaultNodeImageReference   = DefaultNode16ImageReference
)

var (
	Node10ImageReference = DefaultNode10ImageReference
	// Node12ImageReference is the container image reference that
	// is used by GitHub Actions that run using node12.
	Node12ImageReference = DefaultNode12ImageReference
	// Node16ImageReference is the container image reference that
	// is used by GitHub Actions that run using node16.
	Node16ImageReference = DefaultNode16ImageReference
	// Node20ImageReference is the container image reference that
	// is used by GitHub Actions that run using node20.
	Node20ImageReference = DefaultNode20ImageReference
	// Node24ImageReference is the container image reference that
	// is used by GitHub Actions that run using node24.
	Node24ImageReference = DefaultNode24ImageReference
	NodeImageReference   = DefaultNodeImageReference
)
