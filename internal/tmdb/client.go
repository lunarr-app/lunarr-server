package tmdb

import (
	tmdb "github.com/cyruzin/golang-tmdb"
	"github.com/lunarr-app/lunarr-go/internal/config"
	"github.com/lunarr-app/lunarr-go/internal/util"
)

var TmdbClient *tmdb.Client

// IMPORTANT: The following access token is for production usage only and should NOT be shared or used in third-party repositories.
const accessToken = "eyJhbGciOiJIUzI1NiJ9.eyJhdWQiOiIzYzM0NTExNGUxNmZiNjM2NWFiMmQxZjA5Y2I5MjlhNyIsIm5iZiI6MTcyNzU5Mzc3MS41ODU4NjYsInN1YiI6IjVlMzVhMzdmNzZlZWNmMDAxNThmNjliZSIsInNjb3BlcyI6WyJhcGlfcmVhZCJdLCJ2ZXJzaW9uIjoxfQ.Kkgw_UxLJXe5OZeFGtNF-a3dLrdaSjDwMFNj56qiDDc"

func InitTMDBClient() {
	util.Logger.Info().Msg("Initializing TMDB client...")
	cfg := config.Get()

	var client *tmdb.Client
	var err error

	if cfg.TMDb.APIKey != "" {
		client, err = tmdb.Init(cfg.TMDb.APIKey)
	} else {
		client, err = tmdb.InitV4(accessToken)
	}

	if err != nil {
		util.Logger.Fatal().Err(err).Msg("Failed to initialize TMDB client")
	}

	TmdbClient = client
	util.Logger.Info().Msg("TMDB client initialized successfully")
}
