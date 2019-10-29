package models

import (
	"../helpers"
	"github.com/jinzhu/gorm"
)

// Token struct a gorm model
type Token struct {
	gorm.Model
	Token  string `gorm:"type:varchar(255)"`
	User   User   `gorm:"foreignkey:UserID;unique;not null"`
	UserID uint
}

// IfExists token IfExist control method
func (token Token) IfExists() bool {
	conn, _ := Database()
	GetToken := Token{}
	res := conn.Find(&GetToken, token)
	if res.RecordNotFound() {
		return false
	}
	return true
}

// Get get token method
func (token Token) Get() (GetToken Token) {
	conn, _ := Database()
	conn.Where(&token).First(&GetToken)
	return
}

// CreateNew create new token
func (NewToken Token) CreateNew() Token {
	conn, _ := Database()
	NewToken.Token = helpers.RandomString(50)
	conn.Create(&NewToken)
	return NewToken
}

func (UpdatedToken Token) Update() Token {
	conn, _ := Database()
	conn.First(&UpdatedToken, UpdatedToken.ID)
	UpdatedToken.Token = helpers.RandomString(50)
	conn.Save(&UpdatedToken)

	return UpdatedToken
}

// User struct a gorm model
type User struct {
	gorm.Model
	Email    string `gorm:"type:varchar(100);unique;not null"`
	Password string `gorm:"type:varchar(255)" json:"-"`
}

// IfExists user if exists method
func (user User) IfExists() bool {
	conn, _ := Database()
	GetUser := User{}
	Password := ""

	if user.Password != "" {
		Password = user.Password
		user.Password = ""
	}
	res := conn.Find(&GetUser, user)

	if res.RecordNotFound() {

		return false
	}
	if Password != "" {
		if GetUser.Password == Password {
			return true
		}
	}

	return false
}

// Get user get method
func (user User) Get() (GetUser User) {
	conn, _ := Database()
	if user.Password != "" {
		user.Password = ""
	}
	conn.Where(&user).Find(&GetUser)

	return
}

// CreateNew for create a new user record
func (NewUser User) CreateNew() User {
	conn, _ := Database()
	token := Token{}
	if NewUser.Password != "" {
		hashedString, _ := helpers.HashPassword(NewUser.Password)
		NewUser.Password = hashedString
	}
	conn.Create(&NewUser)

	token.User = NewUser
	token.CreateNew()

	return NewUser
}

// Update for update user records
func (UpdatedUser User) Update() User {

	conn, _ := Database()
	conn.First(&UpdatedUser, UpdatedUser.ID)
	if UpdatedUser.Password != "" {
		UpdatedUser.Password, _ = helpers.HashPassword(UpdatedUser.Password)
		token := Token{UserID: UpdatedUser.ID}
		token.Update()
	}
	conn.Save(&UpdatedUser)

	return UpdatedUser
}
