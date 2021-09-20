package domain

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	INACTIVE    = "INACTIVE"
	ACTIVE      = "ACTIVE"
	IN_PROGRESS = "IN_PROGRESS"
	COMPLETE    = "COMPLETE"
)

var (
	ChannelNotFound = errors.New("channel not found")
)

type (
	ChannelRepository interface {
		CreateChannel(context.Context, Channel) (Channel, error)
		GetChannelById(context.Context, string) (Channel, error)
		GetChannelsByQuery(context.Context, interface{}, int, int) ([]Channel, error)
		GetChannelsByQueryCount(context.Context, interface{}) (int64, error)
		// GetChannelsByStatus(context.Context, string) ([]Channel, error)
		UpdateChannelStatus(context.Context, Channel) error
		AddMessage(context.Context, Channel) error
	}

	StatusHistory struct {
		Status    string
		UpdatedBy string
		Timestamp int64
	}

	Message struct {
		MessageFrom string
		Message     string
		IsDeleted   bool
		Timestamp   time.Time
	}

	Channel struct {
		id            primitive.ObjectID
		userFullName  string
		userEmail     string
		repEmail      string
		currentStatus string
		statusHistory []StatusHistory
		messages      []Message
		createdAt     time.Time
		updatedAt     time.Time
	}
)

func NewChannel(id primitive.ObjectID, userEmail, currentStatus string, createdAt, updatedAt time.Time) Channel {
	return Channel{
		id:            id,
		userEmail:     userEmail,
		currentStatus: currentStatus,
		createdAt:     createdAt,
		updatedAt:     updatedAt,
	}
}

func (c *Channel) UpdateRepEmail(repEmail string) {
	c.repEmail = repEmail
}

func (c *Channel) UpdateUserFullName(fullname string) {
	c.userFullName = fullname
}

func (c *Channel) AddMessage(messageFrom, message string, timestamp time.Time) {
	c.messages = append(c.messages, Message{
		MessageFrom: messageFrom,
		Message:     message,
		Timestamp:   timestamp,
	})
}

func (c *Channel) UpdateStatus(status, updatedBy string, timestamp int64) {
	c.currentStatus = status
	c.statusHistory = append(c.statusHistory, StatusHistory{
		Timestamp: timestamp,
		UpdatedBy: updatedBy,
		Status:    status,
	})
}

func (c Channel) StatusHistory() []StatusHistory {
	return c.statusHistory
}

func (c Channel) UserFullName() string {
	return c.userFullName
}

func (c Channel) UserEmail() string {
	return c.userEmail
}

func (c Channel) RepEmail() string {
	return c.repEmail
}

func (c Channel) CurrentStatus() string {
	return c.currentStatus
}

func (c Channel) Messages() []Message {
	return c.messages
}

func (c Channel) Id() primitive.ObjectID {
	return c.id
}

func (c Channel) CreatedAt() time.Time {
	return c.createdAt
}

func (c Channel) UpdatedAt() time.Time {
	return c.updatedAt
}
