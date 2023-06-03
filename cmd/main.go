package main

import (
	"fmt"
	"log"

	"github.com/lunarr-app/lunarr-go/internal/config"
	"github.com/lunarr-app/lunarr-go/internal/db"
	"github.com/lunarr-app/lunarr-go/internal/server"
	"github.com/lunarr-app/lunarr-go/internal/tmdb"
	"github.com/lunarr-app/lunarr-go/internal/util"
)

func main() {
	// Initialize logger instance
	util.InitLogger()

	// Parse command-line flags
	config.InitConfig()
	config.ParseFlags()

	// Initialize the database
	db.InitDatabase()

	// Initialize the TMDB client
	tmdb.InitTMDBClient()

	// Create a new instance of the server
	app := server.New()

	// Start the server on the specified port
	cfg := config.Get()
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	log.Fatal(app.Listen(addr))
}
