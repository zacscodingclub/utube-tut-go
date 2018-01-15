package models

import (
	"errors"
	"fmt"

	"github.com/go-redis/redis"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrInvalidLogin = errors.New("invalid login")
)

type User struct {
	key string
}

func NewUser(u string, h []byte) (*User, error) {
	id, err := client.Incr("user:next-id").Result()
	if err != nil {
		return nil, err
	}

	key := fmt.Sprintf("user:%d", id)
	pipe := client.Pipeline()
	pipe.HSet(key, "id", id)
	pipe.HSet(key, "username", u)
	pipe.HSet(key, "hash", h)
	pipe.HSet("user:by-username", u, id)

	_, err = pipe.Exec()
	if err != nil {
		return nil, err
	}

	return &User{key}, nil
}

func (u *User) GetUserId() (int64, error) {
	return client.HGet(u.key, "id").Int64()
}
func (u *User) GetUsername() (string, error) {
	return client.HGet(u.key, "username").Result()
}

func (u *User) GetHash() ([]byte, error) {
	return client.HGet(u.key, "hash").Bytes()
}

func (u *User) Authenticate(p string) error {
	hash, err := u.GetHash()
	if err != nil {
		return err
	}
	err = bcrypt.CompareHashAndPassword(hash, []byte(p))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return ErrInvalidLogin
	}

	return err
}

func GetUserById(id int64) (*User, error) {
	key := fmt.Sprintf("user:%d", id)
	return &User{key}, nil
}

func GetUserByUsername(username string) (*User, error) {
	id, err := client.HGet("user:by-username", username).Int64()
	if err == redis.Nil {
		return nil, ErrUserNotFound
	} else if err != nil {
		return nil, err
	}

	return GetUserById(id)
}

func RegisterUser(u, p string) error {
	cost := bcrypt.DefaultCost
	hash, err := bcrypt.GenerateFromPassword([]byte(p), cost)
	if err != nil {
		return err
	}
	_, err = NewUser(u, hash)
	return err
}

func AuthenticateUser(u, p string) (*User, error) {
	user, err := GetUserByUsername(u)
	if err != nil {
		return nil, err
	}

	return user, user.Authenticate(p)
}
