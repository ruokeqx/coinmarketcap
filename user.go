package main

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
)

// Model 基类
type Model struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at "`
	// * 代表 null
	DeletedAt *time.Time `json:"deleted_at"`
}

// User 用户信息表
// type UserTable struct {
// 	Model
// 	Nickname string `json:"nickname"`
// }

var models Model

// GetUser 获取用户信息
func GetUser(id uint) (*UserTable, error) {
	var user UserTable
	db := sqlInit()
	err := db.Where("id = ? ", id).First(&user).Error
	if err != nil {
		fmt.Print(err)
		return nil, err
	}
	return &user, nil
}

// AddUser 新增用户信息
func addUser(tx *gorm.DB) (uint, error) {
	var user UserTable
	if err := tx.Create(&user).Error; err != nil {
		fmt.Print(err)
		return 0, err
	}
	return user.ID, nil
}

// ExistUserByID 检查是否存在此用户
func ExistUserByID(id uint) (bool, error) {
	var user UserTable
	db := sqlInit()
	err := db.Select("id").Where("id = ? ", id).First(&user).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		fmt.Print(err)
		return false, err
	}
	if user.ID > 0 {
		return true, nil
	}
	return false, nil
}

// User 用户
type User struct {
	ID uint

	UserName string
	Password string
	Code     string
}

// Register 注册用户
func (u *User) Register() error {
	maps := make(map[string]interface{})

	// 加密
	password, err := Encrypt(u.Password)
	if err != nil {
		return err
	}

	maps["password"] = password

	// 创建 用户信息 与 用户密码
	return AddUserLogin(maps)
}

// PwdLogin 登录用户
func (u *User) PwdLogin() (string, error) {

	user, err := u.getUserLoginInfo()
	if err != nil {
		return "", err
	}

	inputPwd := u.Password
	hashPwd := user.Password

	// 比较 密码
	if err := Compare(inputPwd, hashPwd); err != nil {
		return "", err
	}

	if err := existUserInfo(user.UserID); err != nil {
		return "", err
	}

	// 生成token
	token, merr := GenerateToken(u.UserName)
	if merr != nil {
		fmt.Print(err)
		return "", err
	}

	return token, nil
}

// Register 账号密码注册
// @Summary 账号密码注册
// @accept application/x-www-form-urlencoded
// @Tags auth
// @Produce  json
// @Param auth body api.auth true "账号密码登录/注册"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /auth/register [post]
