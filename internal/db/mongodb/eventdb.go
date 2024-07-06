package mongodb

import (
	"context"
	"go-worker-rabbitmq-kafka-event/internal/appctx"
	"go-worker-rabbitmq-kafka-event/internal/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

const (
	database  = "EVENT_DB"
	eventColl = "EVENTS"
)

type EventDB interface {
	AddEvent(ctx context.Context, event models.Event) (string, error)
	GetEvents(ctx context.Context) ([]*models.Event, error)
	GetEventByMessage(ctx context.Context, message string) (*models.Event, error)
}

type eventDB struct {
	mongoConn *mongo.Client
}

func NewEventDB(mongoConn *mongo.Client) EventDB {
	return &eventDB{
		mongoConn: mongoConn,
	}
}

func (db eventDB) AddEvent(ctx context.Context, event models.Event) (string, error) {
	logger := appctx.FromContext(ctx)
	eventDB := db.mongoConn.Database(database)
	eventColl := eventDB.Collection(eventColl)

	result, err := eventColl.InsertOne(context.Background(), event)
	if err != nil {
		return "", err
	}

	idStr, _ := result.InsertedID.(primitive.ObjectID)

	logger.Info("Event successfully inserted in MongoDB!")
	return idStr.Hex(), nil
}

func (db eventDB) GetEvents(ctx context.Context) ([]*models.Event, error) {
	logger := appctx.FromContext(ctx)
	eventDB := db.mongoConn.Database(database)
	eventColl := eventDB.Collection(eventColl)

	pipeline := []bson.M{
		{"$sort": bson.M{"_id": -1}},
	}

	cursor, err := eventColl.Aggregate(ctx, pipeline)
	if err != nil {
		logger.Warn("GetEvents", zap.String("collection_error", err.Error()))
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	events := make([]*models.Event, 0, len(results))

	for _, result := range results {

		// birthdate := ""
		// if result["DATE"] != nil && result["DATE"].(primitive.DateTime) != 0 {
		// 	birthdate = time.UnixMilli(int64(result["DATE"].(primitive.DateTime))).Format("2006-01-02")
		// }

		events = append(events, &models.Event{
			Message:   result["MESSAGE"].(string),
			Date:      result["DATE"].(time.Time),
			CreatedBy: result["CREATED_BY"].(string),
		})
	}

	return events, nil
}

func (db eventDB) GetEventByMessage(ctx context.Context, message string) (*models.Event, error) {
	logger := appctx.FromContext(ctx)
	eventDB := db.mongoConn.Database(database)
	eventColl := eventDB.Collection(eventColl)

	filters := bson.M{
		"MESSAGE": message,
	}

	pipeline := []bson.M{
		{"$match": filters},
		{"$sort": bson.M{"_id": -1}},
		{"$limit": 1},
	}

	cursor, err := eventColl.Aggregate(ctx, pipeline)
	if err != nil {
		logger.Warn("GetEventByMessage", zap.String("collection_error", err.Error()))
		return nil, err
	}
	defer cursor.Close(ctx)

	var result models.Event
	if cursor.Next(ctx) {
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}
	} else {
		return nil, nil
	}

	return &result, nil
}
