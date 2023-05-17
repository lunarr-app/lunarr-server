package db

import (
	"context"
	"time"

	TMDb "github.com/lunarr-app/golang-tmdb"
	"github.com/lunarr-app/lunarr-go/internal/util"
	"go.mongodb.org/mongo-driver/bson"
)

type MovieWithFiles struct {
	Movie *TMDb.MovieDetails `bson:"movie"`
	Files []string           `bson:"files"`
}

func CheckMovieExists(filePath string) bool {
	filter := bson.M{
		"files": filePath,
	}

	var result bson.M
	err := MoviesLists.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		return false
	}

	return result != nil
}

func InsertMovie(movie *TMDb.MovieDetails, file string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	movieWithFiles := MovieWithFiles{
		Movie: movie,
		Files: []string{file},
	}

	_, err := MoviesLists.InsertOne(ctx, movieWithFiles)
	if err != nil {
		util.Logger.Error().Err(err).Msg("Failed to insert movie into MongoDB")
		return err
	}

	util.Logger.Info().Msgf("Movie inserted successfully: %s", movie.Title)
	return nil
}
