package hub

import (
	"code.google.com/p/go.net/context"
	"github.com/KosyanMedia/burlesque/storage"
)

type (
	Hub struct {
		storage     *storage.Storage
		subscribers []*context.Context
	}
)

func New() (h *Hub) {
	h = Hub{}

	return
}

func (h *Hub) Pub(ctx context.Context) context.Context {
	return ctx
}

func (h *Hub) Sub(ctx context.Context) context.Context {
	return ctx
}
