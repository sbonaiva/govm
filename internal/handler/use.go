package handler

import (
	"context"

	"github.com/sbonaiva/govm/internal/domain"
	"github.com/sbonaiva/govm/internal/gateway"
)

type UseHandler interface {
	Handle(ctx context.Context, use *domain.Use) error
}

type useHandler struct {
	httpGateway gateway.HttpGateway
}

func NewUse(httpGateway gateway.HttpGateway) UseHandler {
	return &useHandler{
		httpGateway: httpGateway,
	}
}

func (u *useHandler) Handle(ctx context.Context, use *domain.Use) error {
	panic("unimplemented")
}
