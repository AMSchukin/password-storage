package repository

import (
	"fmt"
	"github.com/go-pg/pg/v10"
	"github.com/labstack/gommon/log"
	"password-storage/internal/pkg/models/db"
	"time"
)

func NewDBConn() (con *pg.DB) {
	address := fmt.Sprintf("%s:%s", "localhost", "5432")
	options := &pg.Options{
		User:     "admin",
		Password: "12345",
		Addr:     address,
		Database: "postgres",
		PoolSize: 50,
	}
	con = pg.Connect(options)
	if con == nil {
		log.Error("cannot connect to postgres")
	}
	return
}

func CreateNewPasswordData(pg *pg.DB, pd db.PasswordData) error {
	tx, err := pg.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Model(&pd).Insert()
	if err != nil {
		err := tx.Rollback()
		if err != nil {
			return err
		}
		return err
	}
	return tx.Commit()
}

func CreateNewSession(pg *pg.DB, s db.Session) error {
	tx, err := pg.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Model(&s).Insert()
	if err != nil {
		err := tx.Rollback()
		if err != nil {
			return err
		}
		return err
	}
	return tx.Commit()
}

func GetSessionExpiredAt(pg *pg.DB, sessionId string) time.Time {
	session := db.Session{}
	err := pg.Model(&session).Where("session_id = ?0", sessionId).Select()
	if err != nil {
		return time.Time{}
	}
	return session.ExpiredAt
}

func GetSessionExpiredAtByTelegramId(pg *pg.DB, telegramId int64) time.Time {
	session := db.Session{}
	err := pg.Model(&session).Where("telegram_chat_id = ?0", telegramId).Order("expired_at desc").Limit(1).Select()
	if err != nil {
		return time.Time{}
	}
	return session.ExpiredAt
}

func GetAllPasswordData(pg *pg.DB) ([]db.PasswordData, error) {
	pd := make([]db.PasswordData, 10)
	err := pg.Model(&pd).Select()
	if err != nil {
		return nil, err
	}
	return pd, err
}

func GetPasswordDataByKey(pg *pg.DB, key interface{}) (db.PasswordData, error) {
	pd := db.PasswordData{}
	err := pg.Model(&pd).Where("key = ?0", key).Select()
	if err != nil {
		return db.PasswordData{}, err
	}
	return pd, err
}
