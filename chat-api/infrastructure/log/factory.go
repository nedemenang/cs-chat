package log

import (
	"errors"

	"chat-api/adapter/logger"
)

const (
	InstanceLogrusLogger int = iota
)

var (
	errInvalidLoggerInstance = errors.New("invalid log instance")
)

func NewLoggerFactory(instance int) (logger.Logger, error) {
	switch instance {
	case InstanceLogrusLogger:
		return NewLogrusLogger(), nil
	default:
		return nil, errInvalidLoggerInstance
	}
}
