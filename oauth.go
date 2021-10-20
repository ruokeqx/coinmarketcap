package main

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

// Oauth 第三方登录认证
type Oauth struct {
	ID               uint `gorm:"primary_key"`
	UserID           uint
	OauthType        uint   `json:"oauth_type"`
	OauthID          string `json:"oauth_id" gorm:"unique"`
	OauthAccessToken string `json:"access_token"`
	OauthExpires     string `json:"expires_at"`
	Status           uint   `json:"status"`
}

// AddUserOauth 添加用户账号 与 初始化个人信息
func AddUserOauth(userOatuh map[string]interface{}) error {

	db := sqlInit()
	tx := db.Begin()

	// 首先创建 user
	userID, err := addUser(tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	fmt.Println(userOatuh)

	oauthInfo := Oauth{
		UserID:           userID,
		OauthID:          userOatuh["oauth_id"].(string),
		OauthType:        userOatuh["oauth_type"].(uint),
		OauthAccessToken: userOatuh["access_token"].(string),
		OauthExpires:     userOatuh["expires"].(string),
	}
	if err := tx.Create(&oauthInfo).Error; err != nil {
		tx.Rollback()
		fmt.Print(err)
		return err
	}
	tx.Commit()
	return nil
}

// LoginUserOauth 采用密码方式登录
// func LoginUserOauth(maps map[string]interface{}) (*Oauth, error) {
// 	var user Oauth
// 	if err := DB.Where(maps).First(&user).Error; err != nil {
// 		return nil, err
// 	}
// 	return &user, nil
// }

// ExistUserOauth 判断用户账号是否存在
func ExistUserOauth(maps map[string]interface{}) (bool, error) {
	var user Oauth
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
