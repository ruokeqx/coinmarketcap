package main

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

// UserLogin 用户密码登陆认证
type UserLogin struct {
	ID         uint `gorm:"primary_key"`
	UserID     uint
	LoginName  string `gorm:"unique"`
	LoginEmail string `gorm:"unique"`
	LoginPhone string `gorm:"unique"`
	Password   string
	status     uint
}

// 验证码登录
func codeLogin() {

}

// AddUserLogin 添加用户账号 与 初始化个人信息
func AddUserLogin(userLogin map[string]interface{}) error {

	db := sqlInit()
	tx := db.Begin()

	// 首先创建 user
	userID, err := addUser(tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	loginInfo := UserLogin{
		UserID: userID,
	}

	if userLogin["password"] != nil {
		loginInfo.Password = userLogin["password"].(string)
	}

	if userLogin["login_name"] != nil {
		loginInfo.LoginName = userLogin["login_name"].(string)
		goto InsertLogin
	}
	if userLogin["login_phone"] != nil {
		loginInfo.LoginPhone = userLogin["login_phone"].(string)
		goto InsertLogin
	}
	if userLogin["login_email"] != nil {
		loginInfo.LoginEmail = userLogin["login_email"].(string)
	}
InsertLogin:

	if err := tx.Create(&loginInfo).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

// LoginUserLogin 采用密码方式登录
func LoginUserLogin(maps map[string]interface{}) (*UserLogin, error) {
	var user UserLogin
	db := sqlInit()
	if err := db.Where(maps).First(&user).Error; err != nil {
		fmt.Print(err)
		return nil, err
	}
	return &user, nil
}

// ExistUserLogin 判断是否存在此用户账号
func ExistUserLogin(maps map[string]interface{}) (bool, error) {
	var user UserLogin
	db := sqlInit()
	err := db.Select("id").Where(maps).First(&user).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		fmt.Print(err)
		return false, err
	}
	if user.ID > 0 {
		return true, nil
	}
	return false, nil
}

// UserLoginGetUserID 通过用户名 获取 用户ID
func UserLoginGetUserID(maps map[string]interface{}) (uint, error) {
	var user UserLogin
	db := sqlInit()
	err := db.Select("id,user_id").Where(maps).First(&user).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return 0, err
	}
	if user.ID < 1 { // 判断用户账号是否存在
		return 0, err
	}
	if user.UserID < 1 { // 判断用户信息是否存在
		return 0, err
	}

	return user.UserID, nil
}
