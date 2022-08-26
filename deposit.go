package forge

import (
	"context"
	"fmt"
)

type Deposit interface {
	fmt.GoStringer

	Store(context.Context, Ore, *Metal) error
	Retrieve(context.Context, Ore) (*Metal, error)
}
