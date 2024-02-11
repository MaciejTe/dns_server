package loader

import "context"

type Loader interface {
	Load(ctx context.Context) error
}
