package main

import (
	"github.com/garyburd/redigo/redis"
)

var Conn redis.Conn

func InitDB() {
	c, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		Trace.Println(err)
		return
	}
	Conn = c
}
