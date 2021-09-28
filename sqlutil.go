package main

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Coin struct {
	Name string
	Id   int
}

func sqlInit() (db *gorm.DB, err error) {
	// 创建数据库连接
	db, err = gorm.Open("mysql", "ruokeqx:ruokeqx666@(121.196.208.97:3306)/ruokeqx?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		fmt.Print("Connect database error!")
		return
	}
	// defer db.Close()
	return db, err
}

// 存入coin_name及其对应的id
func InsertCoin(db *gorm.DB, coin_name string, id int) {
	db.AutoMigrate(&Coin{})
	tc := Coin{Name: coin_name, Id: id}
	cc := Coin{}
	db.Where("name = ?", tc.Name).First(&cc)
	if cc.Name == "" {
		db.Create(tc)
		fmt.Println(tc, "insert success!")
	} else {
		fmt.Println(tc, "already exists!")
	}
}
