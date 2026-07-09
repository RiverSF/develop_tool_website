package model

import (
	"errors"

	"gorm.io/gorm"
)

type User struct {
	Id   int    `gorm:"primary_key" json:"id"`
	Name string `json:"name"`
}

func (m User) TableName() string {
	return "my_user"
}

type UserModel struct {
}

func NewUserModel() *UserModel {
	return &UserModel{}
}

func (m *UserModel) FindOrCreateByName(name string) (*User, error) {
	var user User
	err := db.Where("name = ?", name).First(&user).Error
	if err == nil {
		return &user, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	user = User{Name: name}
	if err := db.Create(&user).Error; err != nil {
		if isDuplicateEntry(err) {
			if err := db.Where("name = ?", name).First(&user).Error; err != nil {
				return nil, err
			}
			return &user, nil
		}
		return nil, err
	}
	return &user, nil
}

func (m *UserModel) GetUser(id int) *User {
	var user = &User{}
	db.Table("my_user").Where("id = ?", id).First(&user)
	return user
}

func (m *UserModel) GetUserId(name string) *User {
	var user = &User{}
	db.Table("my_user").Where("name = ?", name).First(&user)
	return user
}
