package telegram

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	db         = "telegram_bot"
	collection = "audio_cache"
)

type AudioCache interface {
	GetVideoById(videoId string) (int, error)
	SetVideoId(youtubeId string, telegramID string) error
}

type MongoCahce struct {
	collection *mongo.Collection
}

func NewCache(connection string) (MongoCahce, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(connection))
	if err != nil {
		return MongoCahce{}, err
	}
	collection := client.Database(db).Collection(collection)

	return MongoCahce{collection: collection}, nil
}
