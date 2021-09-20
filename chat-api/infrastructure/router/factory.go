package router

import (
	"chat-api/adapter/repository"
	"errors"
	"time"

	"chat-api/adapter/logger"
	"chat-api/adapter/validator"
)

type Server interface {
	Listen()
}

type Port int64

var (
	errInvalidWebServerInstance = errors.New("invalid router server instance")
)

const (
	InstanceGin int = iota
	InstanceSNS
)

func NewWebServerFactory(
	instance int,
	log logger.Logger,
	dbNoSQL repository.NoSQL,
	validator validator.Validator,
	port Port,
	ctxTimeout time.Duration,
) (Server, error) {
	switch instance {
	case InstanceGin:
		return newGinServer(log, dbNoSQL, validator, port, ctxTimeout), nil
	default:
		return nil, errInvalidWebServerInstance
	}
}
