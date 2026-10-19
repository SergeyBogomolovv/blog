// Package model описывает записи БД и публичные представления API.
package model

import "time"

type Account struct {
	ID           int64      `db:"id"`
	LoginIndex   string     `db:"login_index"`
	LoginEnc     string     `db:"login_enc"`
	EmailEnc     string     `db:"email_enc"`
	NameEnc      string     `db:"name_enc"`
	BioEnc       string     `db:"bio_enc"`
	PasswordHash string     `db:"password_hash"`
	Role         string     `db:"role"`
	Banned       bool       `db:"banned"`
	TokenVersion int        `db:"token_version"`
	ConsentAt    *time.Time `db:"consent_at"`
	CreatedAt    time.Time  `db:"created_at"`
}

type User struct {
	ID        int64     `json:"id"`
	Login     string    `json:"login"`
	Name      string    `json:"name"`
	Bio       string    `json:"bio"`
	Email     string    `json:"email,omitempty"`
	Role      string    `json:"role"`
	Banned    bool      `json:"banned"`
	CreatedAt time.Time `json:"created_at"`
}

type Post struct {
	ID           int64     `db:"id" json:"id"`
	AuthorID     int64     `db:"author_id" json:"author_id"`
	Title        string    `db:"title" json:"title"`
	Body         string    `db:"body" json:"body"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	CommentCount int       `db:"comment_count" json:"comment_count"`
	Author       *User     `json:"author"`
}

type Comment struct {
	ID        int64     `db:"id" json:"id"`
	PostID    int64     `db:"post_id" json:"post_id"`
	AuthorID  int64     `db:"author_id" json:"author_id"`
	Body      string    `db:"body" json:"body"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	Author    *User     `json:"author"`
}
