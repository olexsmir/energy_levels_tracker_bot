package main

import (
	"log/slog"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	cfg := NewConfig()
	if cfg.Debug {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	db, err := NewDB(cfg.DBConn, cfg.TimeZone)
	if err != nil {
		slog.Error("failed to connect to the database", "error", err)
	}

	bot := NewTGBot(cfg.BotToken, cfg.AdminID, db, cfg.Debug)
	scheduler := NewScheduler(bot.Asker)
	httpserver := NewHTTPServer(cfg.HttpPort, db)

	// scheduler.Add("@every 5s")
	scheduler.Add("@every 1h")

	go func() {
		if err := bot.Start(); err != nil {
			slog.Error("failed to start telegram bot", "err", err)
		}
	}()
	go scheduler.Start()
	go httpserver.Start()

	slog.Info("the sever is started", "http_port", cfg.HttpPort)

	select {}
}
