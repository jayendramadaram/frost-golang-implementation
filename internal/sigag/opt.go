package sigag

import "github.com/sirupsen/logrus"

type Options struct {
	Logger *logrus.Logger
	Port   string
}
