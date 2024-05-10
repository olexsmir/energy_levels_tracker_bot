package main

import (
	"log/slog"
	"strconv"

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
		if update.Message != nil {
			chatID := telegoutil.ID(update.Message.Chat.ID)

			if err := b.db.Insert(update.Message.Text); err != nil {
				b.bot.SendMessage(telegoutil.Message(chatID, "failed to save your message"))
				slog.Error("failed to save", "err", err)
				return
			}

			b.bot.SetMessageReaction(&telego.SetMessageReactionParams{
				ChatID:    chatID,
				MessageID: update.Message.MessageID,
				Reaction: []telego.ReactionType{
					&telego.ReactionTypeEmoji{Type: "emoji", Emoji: "👍"},
				},
			})
		}
	}
}
