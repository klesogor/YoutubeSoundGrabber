package telegram

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	database   = "bot_db"
	collection = "audio_cache"
)

type audioCacheRecord struct {
	YoutubeVideoId  string
	TelegramAudioId string
}

type TelegramAudioCache interface {
	TryGetAudioId(videoId string) (string, error)
	SaveAudioIdToCache(youtubeVideoId, telegramAudioId string) error
}

type MongoCache struct {
	collection *mongo.Collection
}

func NewMongoCache(connection string) MongoCache {
	conn, err := mongo.NewClient(options.Client().ApplyURI(connection))
	if err != nil {
		panic(err)
	}
	err = conn.Connect(context.TODO())
	if err != nil {
		panic(err)
	}
	err = conn.Ping(context.TODO(), nil)
	if err != nil {
		panic(err)
	}
	collection := conn.Database(database).Collection(collection)

	return MongoCache{collection: collection}
}

func (c MongoCache) TryGetAudioId(videoId string) (string, error) {
	filter := bson.D{{Key: "YoutubeVideoId", Value: videoId}}
	var res audioCacheRecord
	err := c.collection.FindOne(context.Background(), filter).Decode(&res)
	if err != nil {
		return "", err
	}
	if res.TelegramAudioId != "" {
		return res.TelegramAudioId, nil
	}

	return "", errors.New("Cache miss")
}

func (c MongoCache) SaveAudioIdToCache(youtubeVideoId, telegramAudioId string) error {
	record := audioCacheRecord{YoutubeVideoId: youtubeVideoId, TelegramAudioId: telegramAudioId}
	_, err := c.collection.InsertOne(context.Background(), record)

	return err
}
