package model

import (
	"time"

	"gorm.io/gorm/clause"
)

type UserKey struct {
	Id         int       `gorm:"primary_key" json:"id"`
	UserId     int       `json:"user_id"`
	BrowserKey string    `json:"browser_key"`
	UserAgent  string    `json:"user_agent"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"-"`
}

func (m UserKey) TableName() string {
	return "my_user_key"
}

type UserKeyModel struct {
}

func NewUserKeyModel() *UserKeyModel {
	return &UserKeyModel{}
}

func (m *UserKeyModel) Save(userKey *UserKey) error {
	return db.Save(userKey).Error
}

func (m *UserKeyModel) UpsertByBrowserKey(userKey *UserKey) error {
	return db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "browser_key"}},
		DoUpdates: clause.AssignmentColumns([]string{"user_agent", "updated_at"}),
	}).Create(userKey).Error
}

func (m *UserKeyModel) GetUserKeyList(userId int) []*UserKey {
	var userKeyList = []*UserKey{}
	db.Table("my_user_key").Where("user_id = ?", userId).Find(&userKeyList)
	return userKeyList
}

func (m *UserKeyModel) GetUserIdByKey(key string) *UserKey {
	var userKey = &UserKey{}
	db.Table("my_user_key").Where("browser_key = ?", key).First(&userKey)
	return userKey
}
