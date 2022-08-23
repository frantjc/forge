package forge

import "context"

type Volume interface {
	Remove(context.Context) error
}
