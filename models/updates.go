package models

import "fmt"

type Update struct {
	key string
}

func (u *Update) GetBody() (string, error) {
	return client.HGet(u.key, "body").Result()
}

func (u *Update) GetUser() (*User, error) {
	userID, err := client.HGet(u.key, "user_id").Int64()
	if err != nil {
		return nil, err
	}

	return GetUserById(userID)
}

func NewUpdate(userID int64, b string) (*Update, error) {
	id, err := client.Incr("update:next-id").Result()
	if err != nil {
		return nil, err
	}
	key := fmt.Sprintf("update:%d", id)
	pipe := client.Pipeline()
	pipe.HSet(key, "id", id)
	pipe.HSet(key, "user_id", userID)
	pipe.HSet(key, "body", b)
	pipe.LPush("updates", id)
	pipe.LPush(fmt.Sprintf("user:%d:updates", userID), id)
	_, err = pipe.Exec()
	if err != nil {
		return nil, err
	}
	return &Update{key}, nil
}

func GetAllUpdates() ([]*Update, error) {
	updateIDs, err := client.LRange("updates", 0, 10).Result()
	if err != nil {
		return nil, err
	}

	updates := make([]*Update, len(updateIDs))
	for i, id := range updateIDs {
		key := "update:" + id
		updates[i] = &Update{key}
	}

	return updates, nil
}

func GetUpdates(userID int64) ([]*Update, error) {
	key := fmt.Sprintf("user:%d:updates", userID)
	updateIDs, err := client.LRange(key, 0, 10).Result()
	if err != nil {
		return nil, err
	}

	updates := make([]*Update, len(updateIDs))
	for i, id := range updateIDs {
		key := "update:" + id
		updates[i] = &Update{key}
	}

	return updates, nil
}

func PostUpdate(u int64, b string) error {
	_, err := NewUpdate(u, b)
	return err
}
