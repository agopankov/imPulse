package database

import (
	"context"
	"github.com/agopankov/imPulse/client/internal/emailsender"
	"github.com/agopankov/imPulse/client/internal/emailverify"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

type MongoDB struct {
	client *mongo.Client
}

func NewMongoDB(uri string) (*MongoDB, error) {
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	return &MongoDB{client: client}, nil
}

func (m *MongoDB) SendVerificationEmail(emailAddress string, firstBotID int64, secondBotID int64, postmarkToken string) {
	verificationCode := emailverify.GenerateVerificationCode(6)

	collection := m.client.Database("impulse").Collection("users")
	item := Verification{
		Email:       emailAddress,
		Code:        verificationCode,
		FirstBotID:  firstBotID,
		SecondBotID: secondBotID,
	}
	_, err := collection.InsertOne(context.Background(), item)
	if err != nil {
		log.Fatalf("Got error inserting item: %s", err)
	}

	sender := emailsender.NewEmailSender(postmarkToken)
	sender.SendEmail(emailAddress, "Your verification code", "Your verification code is: "+verificationCode)
}

func (m *MongoDB) VerifyCode(emailAddress string, code string) bool {
	collection := m.client.Database("impulse").Collection("users")

	var item Verification
	err := collection.FindOne(context.Background(), bson.M{"email": emailAddress}).Decode(&item)
	if err != nil {
		log.Fatalf("Error occurred while fetching data from MongoDB %v", err)
	}

	if code == item.Code {
		_, err := collection.UpdateOne(context.Background(), bson.M{"email": emailAddress}, bson.M{"$set": bson.M{"lastverified": time.Now()}})
		if err != nil {
			log.Fatalf("Got error updating LastVerified: %s", err)
		}
		return true
	} else {
		return false
	}
}

func (m *MongoDB) ShouldSendVerificationEmail(emailAddress string) bool {
	databases, err := m.ListDatabases()
	if err != nil {
		log.Fatalf("Failed to get database list: %v", err)
	}
	log.Printf("Databases: %v", databases)

	collection := m.client.Database("impulse").Collection("users")

	var item Verification
	err = collection.FindOne(context.Background(), bson.M{"email": emailAddress}).Decode(&item)

	if err == mongo.ErrNoDocuments {
		return true
	} else if err != nil {
		log.Fatalf("Error occurred while fetching data from MongoDB %v", err)
	}

	if item.LastVerified.IsZero() || time.Since(item.LastVerified) > 240*time.Hour {
		return true
	}

	return false
}

func (m *MongoDB) ListDatabases() ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return m.client.ListDatabaseNames(ctx, bson.M{})
}

func (m *MongoDB) GetAllUsers() ([]Verification, error) {
	collection := m.client.Database("impulse").Collection("users")

	cursor, err := collection.Find(context.Background(), bson.D{})
	if err != nil {
		return nil, err
	}

	var users []Verification
	if err = cursor.All(context.Background(), &users); err != nil {
		return nil, err
	}

	return users, nil
}
