package user

import (
	"fmt"
	"strconv"

	"gopkg.in/redis.v5"
)

type UserRepository struct {
	rds *redis.Client
}

type User struct {
	Id     string
	ChatId int64
	Name   string
}

func NewUserRepository(host, password string, db int) *UserRepository {
	return &UserRepository{redis.NewClient(&redis.Options{
		Addr:     host,
		Password: password,
		DB:       db,
	})}
}

func (ur *UserRepository) AddUser(user *User) error {
	err := ur.rds.Set(user.Id+"_chatid", user.ChatId, 0).Err()
	if err != nil {
		return err
	}
	err = ur.rds.Set(user.Id+"_username", user.Name, 0).Err()
	if err != nil {
		return err
	}
	err = ur.rds.Set(strconv.FormatInt(user.ChatId, 10), user.Id, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

func (ur *UserRepository) GetUser(userId string) (*User, error) {
	chatId, err := ur.rds.Get(userId + "_chatid").Result()
	if err == redis.Nil {
		return nil, fmt.Errorf("%v", "there is not a user registered given the userId")
	}
	if err != nil {
		return nil, err
	}
	userName, err := ur.rds.Get(userId + "_username").Result()
	if err == redis.Nil {
		return nil, fmt.Errorf("%v", "there is not a user registered given the userId")
	}
	if err != nil {
		return nil, err
	}
	pChatId, _ := strconv.ParseInt(chatId, 10, 64)
	return &User{Id: userId, ChatId: pChatId, Name: userName}, nil
}

func (ur *UserRepository) GetUserId(chatId int64) (string, error) {
	userId, err := ur.rds.Get(strconv.FormatInt(chatId, 10)).Result()
	if err == redis.Nil {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return userId, nil
}
