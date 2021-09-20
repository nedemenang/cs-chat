package repository

import (
	"context"
	"time"

	"chat-api/domain"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User schema
type userBSON struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	FirstName string             `bson:"firstName"`
	LastName  string             `bson:"lastName"`
	Email     string             `bson:"email"`
	Password  string             `bson:"password"`
	Role      string             `bson:"role"`
	CreatedAt time.Time          `bson:"createdAt,omitempty"`
	UpdatedAt time.Time          `bson:"updatedAt,omitempty"`
}

type UserNoSQL struct {
	collectionName string
	db             NoSQL
}

func NewUserNoSQL(db NoSQL) UserNoSQL {
	result := UserNoSQL{
		db:             db,
		collectionName: "users",
	}

	// for _, key := range []string{"email"} {
	// 	err := db.EnsureIndex(
	// 		context.Background(),
	// 		result.collectionName,
	// 		bson.D{{Key: key, Value: 1}},
	// 		false,
	// 	)
	// 	if err != nil {
	// 		log.Panic(err)
	// 	}
	// }
	return result
}

func (a UserNoSQL) CreateUser(ctx context.Context, user domain.User) (domain.User, error) {
	var userBSON = userBSON{
		FirstName: user.FirstName(),
		LastName:  user.LastName(),
		Email:     user.Email(),
		Role:      user.Role(),
		Password:  user.Password(),
		CreatedAt: user.CreatedAt(),
		UpdatedAt: user.UpdatedAt(),
	}

	if err := a.db.Store(ctx, a.collectionName, userBSON); err != nil {
		return domain.User{}, errors.Wrap(err, "error creating user")
	}

	return user, nil
}

func (a UserNoSQL) GetUserByEmail(ctx context.Context, emailAddress string) (domain.User, error) {
	var (
		userBSON = &userBSON{}
		query    = bson.M{"email": emailAddress}
	)
	if err := a.db.FindOne(ctx, a.collectionName, query, nil, userBSON); err != nil {
		return domain.User{}, errors.Wrap(err, "error fetching user")
	}

	user := domain.NewUser(
		userBSON.ID,
		userBSON.FirstName,
		userBSON.LastName,
		userBSON.Email,
		userBSON.Password,
		userBSON.CreatedAt,
		userBSON.UpdatedAt,
	)
	user.UpdateRole(userBSON.Role)

	return user, nil
}
