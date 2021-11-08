package service

import (
	"context"
	"fmt"
)

type RPC interface {
	Do(ctx context.Context, method string, out interface{}, args ...interface{}) error
}

type request struct {
	method string
	args   []interface{}
}

func (r request) Key() string {
	return fmt.Sprintf("%s|%s", r.method, r.args)
}
