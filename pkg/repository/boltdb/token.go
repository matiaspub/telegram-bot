package boltdb

import (
	"errors"
	"github.com/boltdb/bolt"
	"strconv"
	"telegram-bot/pkg/config"
	"telegram-bot/pkg/repository"
)

type TokenRepository struct {
	db *bolt.DB
}

func NewTokenRepository(db *bolt.DB) *TokenRepository {
	return &TokenRepository{db: db}
}

func (t *TokenRepository) Save(chatID int64, token string, bucket repository.Bucket) error {
	return t.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte(bucket)).Put(int64ToBytes(chatID), []byte(token))
	})
}

func (t *TokenRepository) Get(chatID int64, bucket repository.Bucket) (string, error) {
	var token string
	err := t.db.View(func(tx *bolt.Tx) error {
		token = string(tx.Bucket([]byte(bucket)).Get(int64ToBytes(chatID)))
		if token == "" {
			return errors.New("Empty token")
		}
		return nil
	})
	return token, err
}

func int64ToBytes(v int64) []byte {
	return []byte(strconv.FormatInt(v, 10))
}

func InitDB(cfg *config.Config) (*bolt.DB, error) {
	db, err := bolt.Open(cfg.DbFile, 0600, nil)
	if err != nil {
		return nil, err
	}

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(repository.AccessToken))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists([]byte(repository.RequestToken))
		if err != nil {
			return err
		}
		return nil
	})
	return db, err
}
