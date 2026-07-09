package model

import (
	"time"
)

type Share struct {
	Id        int       `gorm:"primary_key" json:"id"`
	OpType    int       `json:"op_type"`
	UserId    int       `json:"user_id"`
	Uuid      string    `json:"uuid"`
	Path      string    `json:"path"`
	Name      string    `json:"name"`
	Token     string    `json:"token"`
	Data      string    `json:"data"`
	Status    int       `json:"status"`
	UpdatedAt time.Time `json:"updated_at"`

	UpdatedDate string `gorm:"-" json:"updated_date"`
}

func (m Share) TableName() string {
	return "my_share"
}

type ShareModel struct{}

func NewShareModel() *ShareModel {
	return &ShareModel{}
}

func (m *ShareModel) Save(share *Share) error {
	return db.Create(share).Error
}

func (m *ShareModel) Update(share *Share) error {
	return db.Save(share).Error
}

func (m *ShareModel) GetShareData(path, token string) (data string) {
	_ = db.Table("my_share").
		Where("path = ? AND token = ? AND status = 0", path, token).
		Select("data").
		Scan(&data).Error
	return
}

func (m *ShareModel) GetShareDataById(id, userId int) (share *Share) {
	share = &Share{}
	db.Table("my_share").Where("id = ? AND user_id = ? AND status = 0", id, userId).First(&share)
	return
}

func (m *ShareModel) GetShareListByUserId(opType, userId int) (shareList []*Share) {
	shareList = []*Share{}
	db.Table("my_share").
		Where("op_type = ? AND user_id = ? AND status = 0", opType, userId).
		Order("updated_at desc").
		Find(&shareList)
	return
}

