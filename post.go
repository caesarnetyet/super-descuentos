package main

import (
	"github.com/google/uuid"
	"time"
)

type Post struct {
	ID           uuid.UUID `json:"id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	Url          string    `json:"url"`
	Author       User      `json:"author"`
	Likes        int       `json:"likes"`
	ExpireTime   time.Time `json:"expire_time"`
	CreationTime time.Time `json:"creation_time"`
}
