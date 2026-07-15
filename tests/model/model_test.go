package model_test

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"develop_tools/internal/model"
)

func setupTestDB(t *testing.T) {
	t.Helper()
	dsn := "file:" + strings.ReplaceAll(t.Name(), "/", "_") + "?mode=memory&cache=shared"
	d, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{SingularTable: true},
		Logger:         logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := d.AutoMigrate(
		&model.User{},
		&model.UserKey{},
		&model.Share{},
		&model.Dsp{},
		&model.DspNotice{},
	); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	model.SetDB(d)
	t.Cleanup(func() {
		sqlDB, err := d.DB()
		if err == nil {
			_ = sqlDB.Close()
		}
		model.SetDB(nil)
	})
}

func TestTableNames(t *testing.T) {
	cases := []struct {
		got, want string
	}{
		{model.User{}.TableName(), "my_user"},
		{model.UserKey{}.TableName(), "my_user_key"},
		{model.Share{}.TableName(), "my_share"},
		{model.Dsp{}.TableName(), "my_dsp"},
		{model.DspNotice{}.TableName(), "my_dsp_notice"},
	}
	for _, c := range cases {
		if c.got != c.want {
			t.Errorf("TableName=%q want %q", c.got, c.want)
		}
	}
}

func TestShareAndDspConstants(t *testing.T) {
	if model.ShareOpCollect != 1 || model.ShareOpShare != 2 || model.ShareOpDownload != 3 {
		t.Fatalf("Share op constants unexpected: %d %d %d", model.ShareOpCollect, model.ShareOpShare, model.ShareOpDownload)
	}
	if model.ShareStatusOK != 0 || model.ShareStatusDeleted != 1 {
		t.Fatalf("Share status constants unexpected: %d %d", model.ShareStatusOK, model.ShareStatusDeleted)
	}
	if model.DspMarketOverseas != 0 || model.DspMarketCN != 1 {
		t.Fatalf("Dsp market constants unexpected: %d %d", model.DspMarketOverseas, model.DspMarketCN)
	}
}

func TestCreatedAtUpdatedAtJSONTags(t *testing.T) {
	now := time.Date(2026, 7, 15, 12, 0, 0, 0, time.UTC)
	dsp := model.Dsp{Id: 1, Name: "t", CreatedAt: now, UpdatedAt: now}
	b, err := json.Marshal(dsp)
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatal(err)
	}
	if _, ok := m["created_at"]; !ok {
		t.Fatalf("missing created_at in json: %s", b)
	}
	if _, ok := m["updated_at"]; !ok {
		t.Fatalf("missing updated_at in json: %s", b)
	}
	if _, ok := m["createtime"]; ok {
		t.Fatalf("legacy createtime still present: %s", b)
	}
	if _, ok := m["updatetime"]; ok {
		t.Fatalf("legacy updatetime still present: %s", b)
	}
}

func TestGetNoticeTypeValueAndName(t *testing.T) {
	m := model.NewDspNoticeModel()
	for name, val := range model.NoticeTypeMap {
		got, err := m.GetNoticeTypeValue(name)
		if err != nil || got != val {
			t.Fatalf("GetNoticeTypeValue(%q)=(%d,%v) want %d", name, got, err, val)
		}
		back, err := m.GetNoticeTypeName(val)
		if err != nil || back != name {
			t.Fatalf("GetNoticeTypeName(%d)=(%q,%v) want %q", val, back, err, name)
		}
	}
	if _, err := m.GetNoticeTypeValue("unknown"); err == nil {
		t.Fatal("expected error for unknown notice type")
	}
	if _, err := m.GetNoticeTypeName(999); err == nil {
		t.Fatal("expected error for unknown notice type value")
	}
}

func TestIsDuplicateEntry(t *testing.T) {
	if model.IsDuplicateEntry(nil) {
		t.Fatal("nil should not be duplicate")
	}
	if model.IsDuplicateEntry(errors.New("other")) {
		t.Fatal("generic error should not be duplicate")
	}
	dup := &mysql.MySQLError{Number: 1062, Message: "Duplicate entry"}
	if !model.IsDuplicateEntry(dup) {
		t.Fatal("MySQL 1062 should be duplicate")
	}
	other := &mysql.MySQLError{Number: 1045, Message: "Access denied"}
	if model.IsDuplicateEntry(other) {
		t.Fatal("non-1062 MySQL error should not be duplicate")
	}
}

func TestUserFindOrCreateAndTimestamps(t *testing.T) {
	setupTestDB(t)
	um := model.NewUserModel()

	u1, err := um.FindOrCreateByName("alice")
	if err != nil {
		t.Fatal(err)
	}
	if u1.Id == 0 || u1.Name != "alice" {
		t.Fatalf("unexpected user: %+v", u1)
	}
	if u1.CreatedAt.IsZero() || u1.UpdatedAt.IsZero() {
		t.Fatalf("timestamps not set: created=%v updated=%v", u1.CreatedAt, u1.UpdatedAt)
	}

	u2, err := um.FindOrCreateByName("alice")
	if err != nil {
		t.Fatal(err)
	}
	if u2.Id != u1.Id {
		t.Fatalf("FindOrCreate should return same user: %d vs %d", u1.Id, u2.Id)
	}

	got := um.GetUser(u1.Id)
	if got.Name != "alice" {
		t.Fatalf("GetUser: %+v", got)
	}
	byName := um.GetUserId("alice")
	if byName.Id != u1.Id {
		t.Fatalf("GetUserId: %+v", byName)
	}
}

func TestShareSoftDeleteAndListOrder(t *testing.T) {
	setupTestDB(t)
	sm := model.NewShareModel()

	now := time.Now()
	older := &model.Share{OpType: model.ShareOpCollect, UserId: 1, Uuid: "u1", Path: "/a", Name: "a", Token: "t1", Data: "d1", Status: model.ShareStatusOK, UpdatedAt: now.Add(-time.Minute)}
	if err := sm.Save(older); err != nil {
		t.Fatal(err)
	}
	newer := &model.Share{OpType: model.ShareOpCollect, UserId: 1, Uuid: "u2", Path: "/b", Name: "b", Token: "t2", Data: "d2", Status: model.ShareStatusOK, UpdatedAt: now}
	if err := sm.Save(newer); err != nil {
		t.Fatal(err)
	}

	list := sm.GetShareListByUserId(model.ShareOpCollect, 1)
	if len(list) != 2 {
		t.Fatalf("list len=%d want 2", len(list))
	}
	if list[0].Token != "t2" {
		t.Fatalf("expected updated_at desc, first=%+v", list[0])
	}

	older.Status = model.ShareStatusDeleted
	if err := sm.Update(older); err != nil {
		t.Fatal(err)
	}
	list = sm.GetShareListByUserId(model.ShareOpCollect, 1)
	if len(list) != 1 || list[0].Token != "t2" {
		t.Fatalf("soft-deleted should be hidden: %+v", list)
	}
	if data := sm.GetShareData("/a", "t1"); data != "" {
		t.Fatalf("deleted share data should be empty, got %q", data)
	}
	if data := sm.GetShareData("/b", "t2"); data != "d2" {
		t.Fatalf("active share data=%q", data)
	}
}

func TestDspSaveListAndNoticeCreatedAtFilter(t *testing.T) {
	setupTestDB(t)
	dm := model.NewDspModel()
	nm := model.NewDspNoticeModel()

	now := time.Now()
	d1 := &model.Dsp{Name: "cn1", UniqueKey: "k1", IsCn: model.DspMarketCN, Adm: "adm1", Price: 1.5, UpdatedAt: now.Add(-time.Minute)}
	if err := dm.Save(d1); err != nil {
		t.Fatal(err)
	}
	d2 := &model.Dsp{Name: "cn2", UniqueKey: "k2", IsCn: model.DspMarketCN, Adm: "adm2", Price: 2.5, UpdatedAt: now}
	if err := dm.Save(d2); err != nil {
		t.Fatal(err)
	}
	_ = dm.Save(&model.Dsp{Name: "os", UniqueKey: "k3", IsCn: model.DspMarketOverseas, Adm: "adm3"})

	list := dm.GetDspList(model.DspMarketCN)
	if len(list) != 2 {
		t.Fatalf("cn list len=%d want 2", len(list))
	}
	if list[0].UniqueKey != "k2" {
		t.Fatalf("expected updated_at desc, first=%+v", list[0])
	}
	if list[0].Adm != "" {
		t.Fatalf("GetDspList should Omit adm, got %q", list[0].Adm)
	}
	if dm.GetDspAdmById(d1.Id) != "adm1" {
		t.Fatalf("GetDspAdmById failed")
	}
	if dm.GetDspByUniqueKey("k1").Name != "cn1" {
		t.Fatalf("GetDspByUniqueKey failed")
	}
	if dm.GetDspIDByUniqueKey("k2") != d2.Id {
		t.Fatalf("GetDspIDByUniqueKey failed")
	}

	cutoff := time.Now().Add(-time.Hour).Format("2006-01-02 15:04:05")
	n := &model.DspNotice{DspId: d1.Id, NoticeType: 1, Ip: "1.1.1.1", Ua: "ua"}
	if err := nm.Save(n); err != nil {
		t.Fatal(err)
	}
	if n.CreatedAt.IsZero() {
		t.Fatal("notice CreatedAt not set")
	}
	got := nm.GetDspNoticeByDspId(d1.Id, 0, cutoff)
	if len(got) != 1 || got[0].Ip != "1.1.1.1" {
		t.Fatalf("GetDspNoticeByDspId: %+v", got)
	}
	future := time.Now().Add(time.Hour).Format("2006-01-02 15:04:05")
	if len(nm.GetDspNoticeByDspId(d1.Id, 0, future)) != 0 {
		t.Fatal("created_at filter should exclude older notices")
	}
}

func TestUserKeySaveAndLookup(t *testing.T) {
	setupTestDB(t)
	km := model.NewUserKeyModel()
	uk := &model.UserKey{UserId: 9, BrowserKey: "bk1", UserAgent: "Mozilla"}
	if err := km.Save(uk); err != nil {
		t.Fatal(err)
	}
	if uk.CreatedAt.IsZero() {
		t.Fatal("CreatedAt not set")
	}
	list := km.GetUserKeyList(9)
	if len(list) != 1 || list[0].BrowserKey != "bk1" {
		t.Fatalf("GetUserKeyList: %+v", list)
	}
	found := km.GetUserIdByKey("bk1")
	if found.UserId != 9 {
		t.Fatalf("GetUserIdByKey: %+v", found)
	}
}
