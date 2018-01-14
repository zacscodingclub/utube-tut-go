package models

func GetComments() ([]string, error) {
	return client.LRange("comments", 0, 10).Result()
}

func PostComment(c string) error {
	return client.LPush("comments", c).Err()
}
