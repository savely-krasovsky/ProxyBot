package main

import (
	"context"
	"fmt"
	"github.com/L11R/go-socks5"
	"github.com/asdine/storm"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jinzhu/configor"
	"golang.org/x/net/proxy"
	"log"
	"math/rand"
	"net"
	"net/http"
	"time"
)

var (
	bot    *tgbotapi.BotAPI
	db     *storm.DB
	config Config
)

func main() {
	var err error

	// Load configuration
	err = configor.Load(&config, "_config.yml")
	if err != nil {
		log.Fatal(err)
	}

	// Random seed
	rand.Seed(time.Now().UnixNano())

	var tr http.Transport

	// When you dev it in Russia...
	if config.Proxy.Addr != "" {
		useAuth := true
		if config.Proxy.Username == "" || config.Proxy.Password == "" {
			useAuth = false
		}

		var proxyAuth *proxy.Auth
		if useAuth {
			proxyAuth = &proxy.Auth{
				User:     config.Proxy.Username,
				Password: config.Proxy.Password,
			}
		}

		tr = http.Transport{
			DialContext: func(_ context.Context, network, addr string) (net.Conn, error) {
				socksDialer, err := proxy.SOCKS5("tcp", fmt.Sprintf("%s:%d", config.Proxy.Addr, config.Proxy.Port), proxyAuth, proxy.Direct)
				if err != nil {
					return nil, err
				}

				return socksDialer.Dial(network, addr)
			},
		}
	}

	// Bot init
	bot, err = tgbotapi.NewBotAPIWithClient(config.Token, &http.Client{
		Transport: &tr,
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	// Database init
	db, err = storm.Open("users.db")
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for update := range updates {
			if update.Message == nil {
				continue
			}

			switch update.Message.Command() {
			case "start":
				go StartCommand(update)
			case "redeem":
				go RedeemCommand(update)
			case "update":
				go UpdateCommand(update)
			case "remove":
				go RemoveCommand(update)
			case "make_invitation":
				if update.Message.From.ID == config.AdminID {
					go MakeInvitationCommand(update)
				}
			}
		}
	}()

	// Create a SOCKS5 server
	conf := &socks5.Config{
		AuthMethods: append([]socks5.Authenticator{}, DatabaseAuthenticator{
			DB: db,
		}),
		ConnLimit: 3,
	}

	server, err := socks5.New(conf)
	if err != nil {
		panic(err)
	}

	// Create SOCKS5 proxy on localhost
	server.ListenAndServe("tcp", fmt.Sprintf(":%d", config.Port))
}
