package controller

import (
	"context"
	pkgtx "personaapp/pkg/tx"
)

func init() {
	//TODO: add validators
}

type Storage interface {
	BeginTx(ctx context.Context) (pkgtx.Tx, error)
	NoTx() pkgtx.Tx
}

type Controller struct {
	s Storage
}

func New(s Storage) *Controller {
	return &Controller{s: s}
}
