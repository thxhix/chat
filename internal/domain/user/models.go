package user

import "time"

// Структура из БД
type UserModel struct {
	ID        int64
	Login     string
	Password  string
	CreatedAt time.Time
}

// Бизнес слой
type User struct {
	ID           int64
	Login        string
	PasswordHash string
}
