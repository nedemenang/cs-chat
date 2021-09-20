package repository

import (
	"context"
	"log"
	"time"

	"chat-api/domain"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type StatusHistory struct {
	Status    string `bson:"status"`
	UpdatedBy string `bson:"updatedBy"`
	Timestamp int64  `bson:"timestamp"`
}

type Messages struct {
	MessageFrom string    `bson:"messageFrom"`
	Message     string    `bson:"message"`
	Timestamp   time.Time `bson:"timestamp"`
}

type channelBSON struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	RepEmail      string             `bson:"repEmail"`
	UserEmail     string             `bson:"userEmail"`
	CurrentStatus string             `bson:"currentStatus"`
	StatusHistory []StatusHistory    `bson:"statusHistory"`
	Messages      []Messages         `bson:"messages"`
	CreatedAt     time.Time          `bson:"createdAt,omitempty"`
	UpdatedAt     time.Time          `bson:"updatedAt,omitempty"`
}

type ChannelNoSQL struct {
	collectionName string
	db             NoSQL
}

func NewChannelNoSQL(db NoSQL) ChannelNoSQL {
	result := ChannelNoSQL{
		db:             db,
		collectionName: "channels",
	}

	for _, key := range []string{"userEmail", "repEmail", "currentStatus"} {
		err := db.EnsureIndex(
			context.Background(),
			result.collectionName,
			bson.D{{Key: key, Value: 1}},
			false,
		)
		if err != nil {
			log.Panic(err)
		}
	}
	return result
}

func (a ChannelNoSQL) CreateChannel(ctx context.Context, channel domain.Channel) (domain.Channel, error) {
	var channelBSON = channelBSON{
		ID:            channel.Id(),
		RepEmail:      channel.RepEmail(),
		UserEmail:     channel.UserEmail(),
		CurrentStatus: channel.CurrentStatus(),
		CreatedAt:     channel.CreatedAt(),
		UpdatedAt:     channel.UpdatedAt(),
	}

	for _, status := range channel.StatusHistory() {
		channelBSON.StatusHistory = append(channelBSON.StatusHistory, StatusHistory{
			Status:    status.Status,
			UpdatedBy: status.UpdatedBy,
			Timestamp: status.Timestamp,
		})
	}

	if err := a.db.Store(ctx, a.collectionName, channelBSON); err != nil {
		return domain.Channel{}, errors.Wrap(err, "error creating channel")
	}

	return channel, nil
}

func (a ChannelNoSQL) GetChannelById(ctx context.Context, id string) (domain.Channel, error) {

	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.Channel{}, errors.Wrap(err, "error converting id")
	}

	var (
		channelBSON = &channelBSON{}
		query       = bson.M{"_id": idHex}
	)

	if err := a.db.FindOne(ctx, a.collectionName, query, nil, channelBSON); err != nil {
		switch err {
		case mongo.ErrNoDocuments:
			return domain.Channel{}, domain.ErrUserNotFound
		default:
			return domain.Channel{}, errors.Wrap(err, "error fetching user")
		}
	}
	channel := domain.NewChannel(
		channelBSON.ID,
		channelBSON.UserEmail,
		channelBSON.CurrentStatus,
		channelBSON.CreatedAt,
		channelBSON.UpdatedAt,
	)
	channel.UpdateRepEmail(channelBSON.RepEmail)
	for _, status := range channelBSON.StatusHistory {
		channel.UpdateStatus(status.Status, status.UpdatedBy, status.Timestamp)
	}

	for _, message := range channelBSON.Messages {
		channel.AddMessage(message.MessageFrom, message.Message, message.Timestamp)
	}

	return channel, nil
}

func (a ChannelNoSQL) GetChannelsByQueryCount(ctx context.Context, query interface{}) (int64, error) {

	count, err := a.db.FindCount(ctx, a.collectionName, query)
	if err != nil {
		switch err {
		case mongo.ErrNilDocument:
			return 0, errors.Wrap(domain.ErrUserNotFound, "error counting channels")
		default:
			return 0, errors.Wrap(err, "error counting channels")
		}
	}
	return count, nil
}

func (a ChannelNoSQL) GetChannelsByQuery(ctx context.Context, query interface{}, start, limit int) ([]domain.Channel, error) {
	// lookupStage := bson.D{{"$lookup", bson.D{{"from", "podcasts"}, {"localField", "podcast"}, {"foreignField", "_id"}, {"as", "podcast"}}}}
	// unwindStage := bson.D{{"$unwind", bson.D{{"path", "$podcast"}, {"preserveNullAndEmptyArrays", false}}}}

	// showLoadedCursor, err := episodesCollection.Aggregate(ctx, mongo.Pipeline{lookupStage, unwindStage})
	// if err != nil {
	// 	panic(err)
	// }
	// var showsLoaded []bson.M
	// if err = showLoadedCursor.All(ctx, &showsLoaded); err != nil {
	// 	panic(err)
	// }
	// fmt.Println(showsLoaded)

	findOptions := options.Find()
	findOptions.SetSkip(int64(start))
	findOptions.SetLimit(int64(limit))
	var channelBSONs = make([]channelBSON, 0)
	if err := a.db.FindAll(ctx, a.collectionName, query, &channelBSONs, findOptions); err != nil {
		switch err {
		case mongo.ErrNilDocument:
			return []domain.Channel{}, errors.Wrap(domain.ErrUserNotFound, "error listing channels")
		default:
			return []domain.Channel{}, errors.Wrap(err, "error listing channels")
		}
	}

	var channels = make([]domain.Channel, 0)

	for _, channelBSON := range channelBSONs {
		var channel = domain.NewChannel(
			channelBSON.ID,
			channelBSON.UserEmail,
			channelBSON.CurrentStatus,
			channelBSON.CreatedAt,
			channelBSON.UpdatedAt,
		)
		channel.UpdateRepEmail(channelBSON.RepEmail)
		channels = append(channels, channel)
	}

	return channels, nil
}

func (a ChannelNoSQL) UpdateChannelStatus(ctx context.Context, channel domain.Channel) error {

	statusHistory := make([]StatusHistory, 0)
	for _, doc := range channel.StatusHistory() {
		statusHistory = append(statusHistory, StatusHistory{
			Status:    doc.Status,
			UpdatedBy: doc.UpdatedBy,
			Timestamp: doc.Timestamp,
		})
	}
	var (
		query  = bson.M{"_id": channel.Id()}
		update = bson.M{"$set": bson.M{"statusHistory": statusHistory, "repEmail": channel.RepEmail(), "currentStatus": channel.CurrentStatus()}}
	)

	if err := a.db.Update(ctx, a.collectionName, query, update); err != nil {
		switch err {
		case mongo.ErrNilDocument:
			return errors.Wrap(domain.ErrUserNotFound, "error updating status")
		default:
			return errors.Wrap(err, "error updating status")
		}
	}
	return nil
}

func (a ChannelNoSQL) AddMessage(ctx context.Context, channel domain.Channel) error {

	messages := make([]Messages, 0)
	for _, message := range channel.Messages() {
		messages = append(messages, Messages{
			MessageFrom: message.MessageFrom,
			Message:     message.Message,
			Timestamp:   message.Timestamp,
		})
	}

	var (
		query  = bson.M{"_id": channel.Id()}
		update = bson.M{"$set": bson.M{"messages": messages}}
	)

	if err := a.db.Update(ctx, a.collectionName, query, update); err != nil {
		switch err {
		case mongo.ErrNilDocument:
			return errors.Wrap(domain.ErrUserNotFound, "error adding message")
		default:
			return errors.Wrap(err, "error adding message")
		}
	}
	return nil
}
