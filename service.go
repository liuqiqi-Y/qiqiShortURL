package main

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/mattheath/base62"
)

const (
	URLIDKEY          = "next.url.id"
	ShortURLkey       = "shortURL:%s->url"
	URLHashKey        = "urlhash:%s->shortURL"
	ShortURLDetailKey = "shortURL:%s->detail"
)

type URLDetail struct {
	URL                 string `json:"url"`
	CreatedAt           string `json:"created_at"`
	ExpirationInMinutes int    `json:"expiration_in_minutes"`
}

func Shorten(url string, exp int) (string, error) {
	h := toSha1(url)
	d, err := Conn.Do("GET", fmt.Sprintf(URLHashKey, h))
	if d != nil && err == nil {
		return redis.String(d, err)
	}
	if err != nil {
		return "", err
	}
	if d == nil && err == nil {
		_, err := Conn.Do("Incr", URLIDKEY)
		if err != nil {
			return "", err
		}
		s, err := redis.String(Conn.Do("GET", URLIDKEY))
		if err != nil {
			return "", err
		}
		d, _ := strconv.Atoi(s)
		eid := base62.EncodeInt64(int64(d))
		t := strconv.Itoa((exp * 60))
		_, err = Conn.Do("SET", fmt.Sprintf(ShortURLkey, eid), url, "EX", t)
		if err != nil {
			return "", err
		}
		_, err = Conn.Do("SET", fmt.Sprintf(URLHashKey, h), eid, "EX", t)
		if err != nil {
			return "", err
		}
		detail, err := json.Marshal(
			&URLDetail{
				URL:                 url,
				CreatedAt:           time.Now().String(),
				ExpirationInMinutes: exp,
			})
		if err != nil {
			return "", err
		}
		_, err = Conn.Do("SET", fmt.Sprintf(ShortURLDetailKey, eid), detail, "EX", t)
		if err != nil {
			return "", err
		}
		return eid, nil
	}
	return "", nil
}

func toSha1(url string) string {
	data := []byte(url)
	s := fmt.Sprintf("%x", sha1.Sum(data))
	return s
}

func URLInfo(eid string) (interface{}, error) {
	d, err := Conn.Do("GET", fmt.Sprintf(ShortURLDetailKey, eid))
	if d == nil && err == nil {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return redis.String(d, err)
}
func ShortURLToURL(eid string) (string, error) {
	d, err := Conn.Do("GET", fmt.Sprintf(ShortURLkey, eid))
	if d == nil && err == nil {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return redis.String(d, err)
}
