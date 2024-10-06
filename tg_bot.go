package main

import (
	"fmt"
	"log/slog"
	"strconv"
	"time"
	"unicode"

	"github.com/mymmrac/telego"
	"github.com/mymmrac/telego/telegoutil"
)

type TGBot struct {
	token   string
	debug   bool
	adminID telego.ChatID
	bot     *telego.Bot
	db      *DB
}

func NewTGBot(tok, adminID string, db *DB, debug bool) *TGBot {
	id, _ := strconv.Atoi(adminID)
	return &TGBot{
		token:   tok,
		debug:   debug,
		adminID: telegoutil.ID(int64(id)),
		db:      db,
	}
}

func (b *TGBot) Start() error {
	bot, err := telego.NewBot(b.token, telego.WithDefaultLogger(b.debug, true))
	if err != nil {
		return err
	}

	b.bot = bot

	updates, err := bot.UpdatesViaLongPolling(nil)
	if err != nil {
		return err
	}

	b.handleMessages(updates)
	return nil
}

func (b *TGBot) Stop() error {
	b.bot.StopLongPolling()
	return nil
}

func (b *TGBot) Asker() {
	b.bot.SendMessage(telegoutil.Message(b.adminID, "what is your energy level, right now?"))
}

func (b *TGBot) handleMessages(updates <-chan telego.Update) {
	for update := range updates {
		slog.Debug("new message", "user_id", update.Message.From.ID, "message", update.Message.Text)
		if update.Message != nil && update.Message.Chat.ID == b.adminID.ID {
			chatID := telegoutil.ID(update.Message.Chat.ID)

			if err := validateMsg(update.Message.Text); err != nil {
				_, _ = b.bot.SendMessage(telegoutil.Message(chatID, err.Error()))
				slog.Error("message is not valid", "err", err)
				return
			}

			t := time.Now().In(location)
			if err := b.db.Insert(update.Message.Text, t.Hour(), t); err != nil {
				_, _ = b.bot.SendMessage(telegoutil.Message(chatID, "failed to save your message"))
				slog.Error("failed to save", "err", err)
				return
			}

			if err := b.bot.SetMessageReaction(&telego.SetMessageReactionParams{
				ChatID:    chatID,
				MessageID: update.Message.MessageID,
				Reaction: []telego.ReactionType{
					&telego.ReactionTypeEmoji{Type: "emoji", Emoji: "👍"},
				},
			}); err != nil {
				slog.Error("failed to react to user message")
				return
			}
		}
	}
}

func validateMsg(inp string) error {
	slog.Debug("input message", "inp", inp)

	if !unicode.IsDigit(rune(inp[0])) {
		return fmt.Errorf("the message must start with the rateing")
	}

	i, err := strconv.Atoi(string(inp[0]))
	if err != nil {
		return err
	}

	if i < 0 || i > 5 {
		return fmt.Errorf("the rating must be between 0 and 5")
	}

	return nil
}
