package sigag

import "go.uber.org/zap"

type Options struct {
	Logger *zap.Logger
	Port   string
}
