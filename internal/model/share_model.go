package model

import (
	"time"
)

// Share op_type / status 约定（与前端 _footer.html 一致）
const (
	ShareOpCollect  = 1 // 收藏
	ShareOpShare    = 2 // 分享
	ShareOpDownload = 3 // 下载

	ShareStatusOK      = 0 // 正常
	ShareStatusDeleted = 1 // 软删除
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
		Where("path = ? AND token = ? AND status = ?", path, token, ShareStatusOK).
		Select("data").
		Scan(&data).Error
	return
}

func (m *ShareModel) GetShareDataById(id, userId int) (share *Share) {
	share = &Share{}
	db.Table("my_share").Where("id = ? AND user_id = ? AND status = ?", id, userId, ShareStatusOK).First(&share)
	return
}

func (m *ShareModel) GetShareListByUserId(opType, userId int) (shareList []*Share) {
	shareList = []*Share{}
	db.Table("my_share").
		Where("op_type = ? AND user_id = ? AND status = ?", opType, userId, ShareStatusOK).
		Order("updated_at desc").
		Find(&shareList)
	return
}
