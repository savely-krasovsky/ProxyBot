package main

import (
	"fmt"
	"github.com/asdine/storm"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"strings"
	"time"
)

// Shows all creds, if there is no user in DB, generates new one
func StartCommand(update tgbotapi.Update) {
	var user User

	// Try to get user
	err := db.One("ID", update.Message.From.ID, &user)
	if err != nil {
		log.Println(err)
	}

	// If not found generate new creds and save
	if err == storm.ErrNotFound {
		var users []User
		err := db.All(&users)
		if err != nil {
			log.Println(err)
			return
		}

		if len(users) > config.Limit {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, `Sorry, bot has reached the limit of users! Contact your administrator.`)
			bot.Send(msg)

			return
		}

		if !config.Private {
			err = db.Save(GetUserWithRandomCreds(update.Message.From.ID))
			if err != nil {
				log.Println(err)
				return
			}
		} else {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, `Bot is in private mode. Contact your administrator to get invitation code.`)
			bot.Send(msg)

			return
		}
	} else if err != nil {
		return
	}

	// If we didn't find user before, but already save
	if user == (User{}) {
		err := db.One("ID", update.Message.From.ID, &user)
		if err != nil {
			log.Println(err)
			return
		}
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf(
		`<b>Your credentials:</b>

Server: <code>%s</code>
Port: <code>%d</code>
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
	msg.ReplyMarkup = GetApplyButton(user.Username, user.Password)
	bot.Send(msg)
}

// Updates user, if there are not arguments (username and password) generates defaults
func UpdateCommand(update tgbotapi.Update) {
	var user User

	// Try to get user
	err := db.One("ID", update.Message.From.ID, &user)
	if err != nil {
		log.Println(err)
	}

	if err == storm.ErrNotFound {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, `Sorry, I can't find you, try to register first: /start`)
		bot.Send(msg)

		return
	}

	if update.Message.CommandArguments() == "" {
		user = *GetUserWithRandomCreds(update.Message.From.ID)
	} else {
		user.ID = update.Message.From.ID
		user.CreatedAt = time.Now()

		info := strings.Split(update.Message.CommandArguments(), " ")

		if len(info) == 1 {
			user.Username = info[0]
			user.Password = RandStringBytes(16)
		} else {
			user.Username = info[0]
			user.Password = info[1]
		}
	}

	err = db.Save(&user)
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
		`<b>Your updated credentials:</b>

Username: <code>%s</code>
Password: <code>%s</code>

Created at: <code>%s</code>`,
		user.Username,
		user.Password,
		user.CreatedAt.Format("02.01 / 15:04:05 MST"),
	))
	msg.ParseMode = "HTML"
	msg.ReplyMarkup = GetApplyButton(user.Username, user.Password)
	bot.Send(msg)
}

// Just removes user from DB
func RemoveCommand(update tgbotapi.Update) {
	err := db.DeleteStruct(&User{ID: update.Message.From.ID})
	if err != nil {
		log.Println(err)
	}

	// If not found, send message
	if err == storm.ErrNotFound {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, `Your profile already removed from database!`)
		bot.Send(msg)

		return
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, `Your profile has been removed from database, bye!`)
	bot.Send(msg)
}

func MakeInvitationCommand(update tgbotapi.Update) {
	invID := RandStringBytes(8)

	err := db.Set("invitations", invID, false)
	if err != nil {
		log.Println(err)
		return
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Invitation code: <code>%s</code>", invID))
	msg.ParseMode = "HTML"
	bot.Send(msg)
}

func RedeemCommand(update tgbotapi.Update) {
	if update.Message.CommandArguments() != "" {
		var isRedeemed bool
		err := db.Get("invitations", update.Message.CommandArguments(), &isRedeemed)
		if err != nil {
			log.Println(err)
			return
		}

		if !isRedeemed {
			err = db.Save(GetUserWithRandomCreds(update.Message.From.ID))
			if err != nil {
				log.Println(err)
				return
			}

			err := db.Set("invitations", update.Message.CommandArguments(), true)
			if err != nil {
				log.Println(err)
				return
			}

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, `You invitation code is correct! Click /start to get you credentials.`)
			bot.Send(msg)

			return
		} else {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, `Sorry but your code already redeemed! Contact your administrator to get the new one.`)
			bot.Send(msg)

			return
		}
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, `You need pass code as param right after <code>/redeem</code> command!`)
	msg.ParseMode = "HTML"
	bot.Send(msg)
}
