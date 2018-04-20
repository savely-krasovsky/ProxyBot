package main

import (
	"fmt"
	"github.com/asdine/storm"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"strings"
	"time"
)

func StartCommand(update tgbotapi.Update) {
	var user User

	// try to get user
	err := db.One("ID", update.Message.From.ID, &user)
	if err != nil {
		log.Println(err)
	}

	// if not found generate new creds and save
	if err == storm.ErrNotFound {
		u := fmt.Sprintf("user%d", update.Message.From.ID)
		p := RandStringBytes(16)
		t := time.Now()

		err := db.Save(&User{
			ID:        update.Message.From.ID,
			Username:  u,
			Password:  p,
			CreatedAt: t,
		})
		if err != nil {
			log.Println(err)
			return
		}
	} else if err != nil {
		return
	}

	// if we didn't find user before, but already save
	if user == (User{}) {
		err := db.One("ID", update.Message.From.ID, &user)
		if err != nil {
			log.Println(err)
			return
		}
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf(
		`Your credentials:
Server: <code>%s</code>
Port: <code>%s</code>
Username: <code>%s</code>
Password: <code>%s</code>

Created at: <code>%s</code>`,
		config.Addr,
		config.Port,
		user.Username,
		user.Password,
		user.CreatedAt.Format("02.01 / 15:04:05 MST"),
	))
	msg.ParseMode = "HTML"
	bot.Send(msg)
}

func UpdateCommand(update tgbotapi.Update) {
	var user User
	user.ID = update.Message.From.ID
	user.CreatedAt = time.Now()

	if update.Message.CommandArguments() == "" {
		user.Username = fmt.Sprintf("user%d", update.Message.From.ID)
		user.Password = RandStringBytes(16)
	} else {
		info := strings.Split(update.Message.CommandArguments(), " ")

		if len(info) == 1 {
			user.Username = info[0]
			user.Password = RandStringBytes(16)
		} else {
			user.Username = info[0]
			user.Password = info[1]
		}
	}

	err := db.Save(&user)
	if err != nil {
		log.Println(err)
	}

	if err == storm.ErrAlreadyExists {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, `User with this username already exists.`)
		bot.Send(msg)

		return
	} else if err != nil {
		return
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf(
		`Your updated credentials:
Username: <code>%s</code>
Password: <code>%s</code>

Created at: <code>%s</code>`,
		user.Username,
		user.Password,
		user.CreatedAt.Format("02.01 / 15:04:05 MST"),
	))
	msg.ParseMode = "HTML"
	bot.Send(msg)
}

func RemoveCommand(update tgbotapi.Update) {
	err := db.DeleteStruct(&User{ID: update.Message.From.ID})
	if err != nil {
		log.Println(err)
	}

	// if not found, send message
	if err == storm.ErrNotFound {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, `Your profile already removed from database!`)
		bot.Send(msg)

		return
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, `Your profile has been removed from database, bye!`)
	bot.Send(msg)
}
