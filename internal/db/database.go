package db

import (
	"os"
	"path/filepath"

	"github.com/lunarr-app/lunarr-go/internal/config"
	"github.com/lunarr-app/lunarr-go/internal/models"
	"github.com/lunarr-app/lunarr-go/internal/util"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDatabase() {
	util.Logger.Info().Msg("Connecting to the SQLite database")

	// Extract the folder path from the database location URI
	dbURI := config.Get().Database.URI
	dbPath := filepath.Dir(dbURI)

	// Create the data folder if it doesn't exist
	err := os.MkdirAll(dbPath, 0755)
	if err != nil {
		util.Logger.Fatal().Err(err).Msg("Failed to create data folder for database")
	}

	// Connect to the SQLite database
	db, err := gorm.Open(sqlite.Open(dbURI), &gorm.Config{})
	if err != nil {
		util.Logger.Fatal().Err(err).Msg("Failed to connect to the SQLite database")
	}

	// Set the database connection in the DB variable
	DB = db

	// AutoMigrate the tables
	MigrateTables()

	// Create text search index
	err = db.Exec("CREATE FULLTEXT INDEX idx_movie_text_search ON movie_with_files (movie_title, movie_original_title, movie_belongs_to_collection_name, files)").Error
	if err != nil {
		util.Logger.Fatal().Err(err).Msg("Failed to create text search index")
	}

	util.Logger.Info().Msg("Database initialization complete")
}

func MigrateTables() {
	err := DB.AutoMigrate(
		&models.UserAccount{},
		&models.MovieWithFiles{},
	)
	if err != nil {
		util.Logger.Fatal().Err(err).Msg("Failed to perform auto migration")
	}
}
