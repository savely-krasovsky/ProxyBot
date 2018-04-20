package main

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"math/rand"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func GetUserWithRandomCreds(id int) *User {
	return &User{
		ID:        id,
		Username:  fmt.Sprintf("user%d", id),
		Password:  RandStringBytes(16),
		CreatedAt: time.Now(),
	}
}

func GetApplyButton(username, password string) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("Apply proxy settings", fmt.Sprintf(
				"https://t.me/socks?server=%s&port=%d&user=%s&pass=%s",
				config.Addr,
				config.Port,
				username,
				password,
			)),
		),
	)
}
