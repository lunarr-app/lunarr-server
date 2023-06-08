package db

import (
	"encoding/json"
	"strconv"

	"github.com/dgraph-io/badger/v4"
	TMDb "github.com/lunarr-app/golang-tmdb"
	"github.com/lunarr-app/lunarr-go/internal/models"
)

func CountMovies() int64 {
	var count int64
	DB.Model(&models.MovieWithFiles{}).Count(&count)
	return count
}

func CheckMovieExists(filePath string) bool {
	var count int64
	DB.Model(&models.MovieWithFiles{}).Where("location = ?", filePath).Count(&count)
	return count > 0
}

func InsertMovie(movie *TMDb.MovieDetails, path string) error {
	movieWithFiles := models.MovieWithFiles{
		TMDbID:   int32(movie.ID),
		Location: path,
	}

	// Encode movie data into JSON
	movieData, err := json.Marshal(movie)
	if err != nil {
		return err
	}

	// Store the movie data in Badger
	err = BadgerDB.Update(func(txn *badger.Txn) error {
		key := "tmdb_movie:" + strconv.FormatInt(movie.ID, 10)
		return txn.Set([]byte(key), movieData)
	})
	if err != nil {
		return err
	}

	// Create the movie record in SQLite
	err = DB.Create(&movieWithFiles).Error
	if err != nil {
		return err
	}

	return nil
}

func FindMovieByTmdbID(tmdbID int) (*models.MovieWithFiles, error) {
	var movie models.MovieWithFiles
	err := DB.Where("tmdb_id = ?", tmdbID).First(&movie).Error
	if err != nil {
		return nil, err
	}

	return &movie, nil
}

func FindMovieMetadata(tmdbID int) (*TMDb.MovieDetails, error) {
	// Retrieve the movie data from Badger using the key format "tmdb_movie:<tmdbID>"
	var movieData *TMDb.MovieDetails
	err := BadgerDB.View(func(txn *badger.Txn) error {
		key := "tmdb_movie:" + strconv.Itoa(tmdbID)
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}

		err = item.Value(func(val []byte) error {
			return json.Unmarshal(val, &movieData)
		})
		return err
	})
	if err != nil {
		return nil, err
	}

	return movieData, nil
}

func DeleteMovieByTmdbID(tmdbID int) error {
	err := DB.Delete(&models.MovieWithFiles{}, "tmdb_id = ?", tmdbID).Error
	if err != nil {
		return err
	}

	return nil
}
