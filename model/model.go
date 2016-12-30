package model

import (
	"time"
)

type Model struct {
	ID        uint64 `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

type Log struct {
	Model
	Level   string
	Message string
	Type    string
}

type Site struct {
	Model
	Name     string
	Classify string
	Url      string
	Protocol string

	Contents []Content
}

type Content struct {
	Model
	Url      string `gorm:"type:text;"`
	SiteId   uint64 `gorm:"index"`
	Encoding string
	Status   int
	Code     int
	Content  []byte `gorm:"type:mediumblob"`
}

type Body struct {
	Model
	Url     string `gorm:"type:text;"`
	SiteId  uint64 `gorm:"index"`
	Status  int
	Code    int
	Content []byte
}

type Url struct {
	Model
	Url    string
	Status int
}
