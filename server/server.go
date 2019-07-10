package server

import "context"

type ChatServer interface {
	Start(ctx context.Context, address string) error
	Close(ctx context.Context) error
}
