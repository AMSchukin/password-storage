package db

import "time"

type PasswordData struct {
	tableName struct{} `pg:"data"`
	Key       string   `pg:"key,pk"`
	Value     string   `pg:"value"`
}

type Session struct {
	tableName  struct{}  `pg:"session"`
	SessionId  string    `pg:"session_id,pk"`
	TelegramId int64     `pg:"telegram_chat_id"`
	ExpiredAt  time.Time `pg:"expired_at"`
}
