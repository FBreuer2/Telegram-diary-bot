package BotDatabase

import (
	"encoding/json"
	"time"

	"github.com/FBreuer2/telegram-diary/entity"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jinzhu/gorm"
)

type BotDatabase struct {
	handle *gorm.DB
}

func New(filePath string) (*BotDatabase, error) {
	db, err := gorm.Open("sqlite3", filePath)

	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&entity.User{})
	db.AutoMigrate(&entity.Entry{})

	return &BotDatabase{
		handle: db,
	}, nil
}

func (bD *BotDatabase) Close() {
	bD.handle.Close()
	return
}

func (bD *BotDatabase) Create(username string) bool {
	newUser := &entity.User{
		Name: username,
	}

	// XXX: Proper error handling here
	bD.handle.Create(&newUser)
	return true
}

func (bD *BotDatabase) Exists(username string) bool {
	var user entity.User

	bD.handle.Where("name = ?", username).First(&user)

	return (user.Name != "")
}

func (bD *BotDatabase) AddText(text string, username string) {
	var entry = &entity.Entry{
		UserName:  username,
		EntryType: entity.Text,
		Created:   time.Now(),
		Content:   text,
	}

	bD.addEntry(entry)
}

// XXX: Figure out how to get a text here
func (bD *BotDatabase) AddAndDownloadImage(photos *[]tgbotapi.PhotoSize, username string) {

}

func (bD *BotDatabase) AddLocation(location *tgbotapi.Location, username string) {
	val, _ := json.Marshal(location)

	var entry = &entity.Entry{
		UserName:  username,
		EntryType: entity.Location,
		Created:   time.Now(),
		Content:   string(val),
	}

	bD.addEntry(entry)
	return
}

func (bD *BotDatabase) addEntry(newEntry *entity.Entry) {
	// XXX: Proper error handling here
	bD.handle.Create(&newEntry)
	return
}
