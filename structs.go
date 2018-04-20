package main

import "time"

type User struct {
	ID        int    `storm:"id"`
	Username  string `storm:"unique"`
	Password  string
	CreatedAt time.Time `storm:"index"`
}