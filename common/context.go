package common

import "context"

var ctx context.Context
var cancel context.CancelFunc

func init() {
	ctx, cancel = context.WithCancel(context.Background())
}

func GetContext() context.Context {
	return ctx
}

func GetCancel() context.CancelFunc {
	return cancel
}
