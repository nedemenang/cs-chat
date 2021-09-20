package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo/options"
)

type NoSQL interface {
	EnsureIndex(context.Context, string, interface{}, bool) error
	Store(context.Context, string, interface{}) error
	Update(context.Context, string, interface{}, interface{}) error
	FindAll(context.Context, string, interface{}, interface{}, *options.FindOptions) error
	FindOne(context.Context, string, interface{}, interface{}, interface{}) error
	FindCount(context.Context, string, interface{}) (int64, error)
	StartSession() (Session, error)
}

type Session interface {
	WithTransaction(context.Context, func(context.Context) error) error
	EndSession(context.Context)
}
