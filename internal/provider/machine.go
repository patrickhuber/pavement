package provider

import "context"

type Machine interface {
	Create(ctx context.Context) error
	Delete(ctx context.Context) error
}
