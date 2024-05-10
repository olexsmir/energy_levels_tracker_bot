package main

import "os"

type Config struct {
	HttpPort string
	BotToken string
	Debug    bool
	AdminID  string
	DBConn   string
}

func NewConfig() *Config {
	return &Config{
		HttpPort: os.Getenv("HTTP_PORT"),
		BotToken: os.Getenv("BOT_TOKEN"),
		Debug:    os.Getenv("DEBUG") == "true",
		AdminID:  os.Getenv("ADMIN_ID"),
		DBConn:   os.Getenv("DBCONN"),
	}
}
