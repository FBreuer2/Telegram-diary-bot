package entity

import (
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type User struct {
	gorm.Model
	Name string
}

type Entry struct {
	gorm.Model
	UserName  string
	EntryType EntryType
	Created   time.Time
	Content   string
}

type EntryType string

const (
	Text     EntryType = "TEXT"
	Image    EntryType = "IMAGE"
	Document EntryType = "DOC"
	Location EntryType = "LOC"
)
