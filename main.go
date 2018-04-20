package main

import (
	"context"
	"fmt"
	"github.com/armon/go-socks5"
	"github.com/asdine/storm"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"golang.org/x/net/proxy"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"time"
)

var (
	bot *tgbotapi.BotAPI
	db  *storm.DB
)

func main() {
	var err error

	// Random seed
	rand.Seed(time.Now().UnixNano())

	// Get socks flag
	socksEnabled := true
	socksEnabledEnv := os.Getenv("ENABLE_SOCKS")
	if socksEnabledEnv == "" || socksEnabledEnv == "false" {
		socksEnabled = false
	}

	// Get token for Telegram bot
	token := os.Getenv("TOKEN")
	if token == "" {
		log.Fatal("TOKEN env variable not specified!")
	}

	var tr http.Transport

	// When you dev it in Russia...
	if socksEnabled {
		proxyAddr := os.Getenv("PROXY_ADDR")
		proxyPort := os.Getenv("PROXY_PORT")

		proxyUsername := os.Getenv("PROXY_USERNAME")
		proxyPassword := os.Getenv("PROXY_PASSWORD")

		useAuth := true
		if proxyUsername == "" || proxyPassword == "" {
			useAuth = false
		}

		var proxyAuth *proxy.Auth
		if useAuth {
			proxyAuth = &proxy.Auth{
				User:     proxyUsername,
				Password: proxyPassword,
			}
		}

		tr = http.Transport{
			DialContext: func(_ context.Context, network, addr string) (net.Conn, error) {
				socksDialer, err := proxy.SOCKS5("tcp", fmt.Sprintf("%s:%s", proxyAddr, proxyPort), proxyAuth, proxy.Direct)
				if err != nil {
					return nil, err
				}

				return socksDialer.Dial(network, addr)
			},
		}
	}

	// Bot init
	bot, err = tgbotapi.NewBotAPIWithClient(token, &http.Client{
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
			case "update":
				go UpdateCommand(update)
			case "remove":
				go RemoveCommand(update)
			}
		}
	}()

	// Create a SOCKS5 server
	conf := &socks5.Config{
		AuthMethods: append([]socks5.Authenticator{}, socks5.NoAuthAuthenticator{}, DatabaseAuthenticator{
			DB: db,
		}),
	}

	server, err := socks5.New(conf)
	if err != nil {
		panic(err)
	}

	socks5.PermitAll()

	// Create SOCKS5 proxy on localhost
	if err := server.ListenAndServe("tcp", ":1323"); err != nil {
		panic(err)
	}
}
