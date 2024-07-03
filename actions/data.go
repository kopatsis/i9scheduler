package actions

import (
	"context"
	"time"

	"firebase.google.com/go/auth"
	"github.com/sendgrid/sendgrid-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type SubscriptionCancellation struct {
	ID      primitive.ObjectID `bson:"_id,omitempty"`
	SubID   string             `bson:"sub_id"`
	UserID  string             `bson:"user_id"`
	EndTime time.Time          `bson:"end_time"`
}

type UserPayment struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`
	UserMongoID    string             `bson:"userid"`
	Username       string             `bson:"username"`
	Provider       string             `bson:"provider"`
	SubscriptionID string             `bson:"subid"`
	SubLength      string             `bson:"length"`
	EndDate        primitive.DateTime `bson:"end"`
	Expires        primitive.DateTime `bson:"expires"`
	Processing     bool               `bson:"processing"`
	Ending         bool               `bson:"ending"`
}

type User struct {
	ID                primitive.ObjectID `bson:"_id,omitempty"`
	Name              string             `bson:"name"`
	Username          string             `bson:"username"`
	Paying            bool               `bson:"paying"`
	Provider          string             `bson:"provider"`
	Level             float32            `bson:"level"`
	BannedExercises   []string           `bson:"bannedExer"`
	BannedStretches   []string           `bson:"bannedStr"`
	BannedParts       []int              `bson:"bannedParts"`
	PlyoTolerance     int                `bson:"plyoToler"`
	ExerFavoriteRates map[string]float32 `bson:"exerfavs"`
	ExerModifications map[string]float32 `bson:"exermods"`
	TypeModifications map[string]float32 `bson:"typemods"`
	RoundEndurance    map[int]float32    `bson:"roundendur"`
	TimeEndurance     map[int]float32    `bson:"timeendur"`
	PushupSetting     string             `bson:"pushsetting"`
	LastMinutes       float32            `bson:"lastmins"`
	LastDifficulty    int                `bson:"lastdiff"`
	Assessed          bool               `bson:"assessed"`
	Badges            []string           `bson:"badges"`    //New
	CompletedCount    int                `bson:"completed"` //New
}

func SetUserNotPaying(client *sendgrid.Client, auth *auth.Client, database *mongo.Database, userID string, email bool) error {
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objID}
	update := bson.M{
		"$set": bson.M{
			"paying":   false,
			"provider": "",
		},
	}

	collection := database.Collection("user")
	_, err = collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}

	userPaymentFilter := bson.M{"userid": userID}
	userPaymentCollection := database.Collection("userpaying")
	_, err = userPaymentCollection.DeleteOne(context.TODO(), userPaymentFilter)
	if err != nil {
		return err
	}

	var user User
	err = collection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		return err
	}

	userRecord, err := auth.GetUser(context.Background(), user.Username)
	if err != nil {
		return err
	}

	if email {
		if err := SendOver(client, userRecord.Email, user.Name); err != nil {
			return err
		}
	}

	return nil
}
