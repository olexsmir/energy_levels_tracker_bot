package main

import (
	"os"
	"strconv"
)

type Config struct {
	HttpPort string
	BotToken string
	Debug    bool
	AdminID  string
	DBConn   string
	TimeZone int
}

func NewConfig() *Config {
	tz, err := strconv.Atoi(os.Getenv("TIMEZONE"))
	if err != nil {
		panic(err)
	}

	return &Config{
		HttpPort: os.Getenv("HTTP_PORT"),
		BotToken: os.Getenv("BOT_TOKEN"),
		Debug:    os.Getenv("DEBUG") == "true",
		AdminID:  os.Getenv("ADMIN_ID"),
		DBConn:   os.Getenv("DBCONN"),
		TimeZone: tz,
	}
}
