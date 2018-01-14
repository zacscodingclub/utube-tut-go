package models

import (
	"errors"

	"github.com/go-redis/redis"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrInvalidLogin = errors.New("invalid login")
)

func RegisterUser(u, p string) error {
	cost := bcrypt.DefaultCost
	hash, err := bcrypt.GenerateFromPassword([]byte(p), cost)
	if err != nil {
		return err
	}

	return client.Set("user:"+u, hash, 0).Err()
}

func AuthenticateUser(u, p string) error {
	hash, err := client.Get("user:" + u).Bytes()
	if err == redis.Nil {
		return ErrUserNotFound
	} else if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword(hash, []byte(p))
	if err != nil {
		return ErrInvalidLogin
	}

	return nil
}
