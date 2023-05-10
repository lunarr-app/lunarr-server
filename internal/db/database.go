package db

import (
	"context"
	"time"

	"github.com/lunarr-app/lunarr-go/internal/common"
	"github.com/lunarr-app/lunarr-go/internal/config"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var UsersAccounts *mongo.Collection
var MoviesLists *mongo.Collection
var TvShowsLists *mongo.Collection
var WatchHistory *mongo.Collection

func init() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	common.Logger.Info().Msg("Creating new MongoClient instance")
	clientOptions := options.Client().ApplyURI(config.Get().Database.URI)
	mongoClient, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		common.Logger.Fatal().Err(err).Msg("Failed to connect to MongoDB")
	}
	err = mongoClient.Ping(ctx, nil)
	if err != nil {
		common.Logger.Fatal().Err(err).Msg("Failed to ping MongoDB")
	}

	common.Logger.Info().Msg("Exporting MongoDB collections as typed objects")
	UsersAccounts = mongoClient.Database("lunarr").Collection("users.accounts")
	MoviesLists = mongoClient.Database("lunarr").Collection("movies.lists")
	TvShowsLists = mongoClient.Database("lunarr").Collection("tv_shows.lists")
	WatchHistory = mongoClient.Database("lunarr").Collection("watch.history")

	common.Logger.Info().Msg("Creating indexes on relevant fields for improved query performance")
	UsersAccounts.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.M{"username": 1}, Options: options.Index().SetName("username_index").SetUnique(true)},
		{Keys: bson.M{"api_key": 1}, Options: options.Index().SetName("api_key_index").SetUnique(true)},
	})
	WatchHistory.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.M{"user_id": 1, "tmdb_id": 1},
		Options: options.Index().
			SetName("user_watch_history_index").
			SetUnique(false),
	})

	common.Logger.Info().Msg("Creating text indexes for improved searching performance")
	MoviesLists.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys: bson.M{
				"tmdb.title":                      "text",
				"tmdb.original_title":             "text",
				"tmdb.belongs_to_collection.name": "text",
				"files":                           "text",
			},
			Options: options.Index().
				SetName("movie_text_search_index").
				SetDefaultLanguage("english").
				SetWeights(bson.M{
					"tmdb.title":                      10,
					"tmdb.original_title":             5,
					"tmdb.belongs_to_collection.name": 2,
					"files":                           1,
				}),
		},
	})
	TvShowsLists.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys: bson.M{
				"tmdb.name":          "text",
				"tmdb.original_name": "text",
				"tmdb.tagline":       "text",
				"files":              "text",
			},
			Options: options.Index().
				SetName("tvshow_text_search_index").
				SetDefaultLanguage("english").
				SetWeights(bson.M{
					"tmdb.name":          10,
					"tmdb.original_name": 5,
					"tmdb.tagline":       2,
					"files":              1,
				}),
		},
	})
}