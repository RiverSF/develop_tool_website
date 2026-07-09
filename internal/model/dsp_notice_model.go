package model

import (
	"errors"
	"fmt"
)

type DspNotice struct {
	Id         int    `gorm:"primary_key" json:"id"`
	DspId      int    `json:"dsp_id"`
	NoticeType int    `json:"notice_type"`
	Ip         string `json:"ip"`
	Ua         string `json:"ua"`
	CreateTime string `gorm:"column:createtime" json:"createtime"`
}

func (m DspNotice) TableName() string {
	return "my_dsp_notice"
}

type DspNoticeModel struct {
}

func NewDspNoticeModel() *DspNoticeModel {
	return &DspNoticeModel{}
}

func (m *DspNoticeModel) GetDspNoticeByDspId(dspId, noticeId int, nowDate string) []*DspNotice {
	var dspNotice = []*DspNotice{}
	db.Table("my_dsp_notice").
		Where("dsp_id = ?", dspId).
		Where("id > ?", noticeId).
		Where("createtime > ?", nowDate).
		Order("id asc").
		Find(&dspNotice)
	return dspNotice
}

func (m *DspNoticeModel) Save(dspNotice *DspNotice) error {
	return db.Save(dspNotice).Error
}

var NoticeTypeMap = map[string]int{
	"burl":     1,
	"nurl":     2,
	"lurl":     3,
	"tpnurl":   4,
	"tplurl":   5,
	"tpburl":   6,
	"tpimpurl": 7,
	"tpclkurl": 8,
}

var noticeTypeNameMap = map[int]string{
	1: "burl",
	2: "nurl",
	3: "lurl",
	4: "tpnurl",
	5: "tplurl",
	6: "tpburl",
	7: "tpimpurl",
	8: "tpclkurl",
}

func (m *DspNoticeModel) GetNoticeTypeValue(noticeType string) (int, error) {
	if typeInt, ok := NoticeTypeMap[noticeType]; ok {
		return typeInt, nil
	}

	return 0, errors.New(fmt.Sprintf("invalid noticeType: %s", noticeType))
}

func (m *DspNoticeModel) GetNoticeTypeName(noticeTypeValue int) (string, error) {
	if name, ok := noticeTypeNameMap[noticeTypeValue]; ok {
		return name, nil
	}

	return "", errors.New(fmt.Sprintf("invalid noticeTypeValue: %d", noticeTypeValue))
}
