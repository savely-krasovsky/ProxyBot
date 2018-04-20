package main

import "time"

type User struct {
	ID        int    `storm:"id"`
	Username  string `storm:"unique"`
	Password  string
	CreatedAt time.Time `storm:"index"`
}

type Config struct {
	Token       string `required:"true"`
	Addr        string `required:"true"`
	Port        int    `default:"1080"`
	Proxy       struct {
		Addr     string
		Port     int `default:"1080"`
		Username string
		Password string
	}
}
