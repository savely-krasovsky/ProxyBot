package main

import (
	"math/rand"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"fmt"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func GetApplyButton(username, password string) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("Apply settings", fmt.Sprintf(
				"https://t.me/socks?server=%s&port=%d&user=%s&pass=%s",
				config.Addr,
				config.Port,
				username,
				password,
			)),
		),
	)
}