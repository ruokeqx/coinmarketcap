package main

import (
	"fmt"
	"sync"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Coin struct {
	Name string `gorm:"primary_key"`
	Id   int
}

type CoinPointQuote struct {
	Id          int
	Name        string
	Time        string `gorm:"primary_key"`
	Price       float64
	Volume      string
	MarketCap   string
	BitcoinRate string
	ZhPrice     float64
	ZhVolume    string
	ZhMarketCap string
}

type CoinHistoricalQuote struct {
	Id       int
	Name     string
	Symbol   string
	TimeOpen string `gorm:"primary_key"`
	// timeClose    string
	TimeHigh     string
	TimeLow      string
	OpenPrice    float64
	HighPrice    float64
	LowPrice     float64
	ClosePrice   float64
	Volume       string
	MarketCap    string
	ZhOpenPrice  float64
	ZhHighPrice  float64
	ZhLowPrice   float64
	ZhClosePrice float64
	ZhVolume     string
	ZhMarketCap  string

	// timestamp string
}

type UserTable struct {
	Username string `json:"username"`
	PwdHash  string `json:"password"`
}

var (
	once sync.Once
	db   *gorm.DB
)

func sqlInit() *gorm.DB {
	// 创建数据库连接
	// db, err = gorm.Open("mysql", "ruokeqx:ruokeqx666@(121.196.208.97:3306)/ruokeqx?charset=utf8mb4&parseTime=True&loc=Local")
	once.Do(func() {
		var err error
		db, err = gorm.Open("mysql", "root:root@(192.168.0.165:3306)/db1?charset=utf8mb4&parseTime=True&loc=Local")
		if err != nil {
			panic("Connect database error!")
		}
	})

	// defer db.Close()
	return db
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

// 存入chart数据
func InsertChart(db *gorm.DB, point CoinPointQuote) {
	// db.AutoMigrate(&CoinHistoricalQuote{})
	if !db.HasTable("chart-" + point.Name) {
		db.Table("chart-" + point.Name).CreateTable(&CoinPointQuote{})
	}
	th := CoinPointQuote{}
	db.Table("chart-"+point.Name).Where("Time = ?", point.Time).First(&th)
	if th.Time == "" {
		db.Table("chart-" + point.Name).Create(point)
		fmt.Println(point, "insert success!")
	} else {
		fmt.Println(point, "already exists!")
	}
}

// 存入历史数据
func InsertHistory(db *gorm.DB, quote CoinHistoricalQuote) {
	// db.AutoMigrate(&CoinHistoricalQuote{})
	if !db.HasTable("history-" + quote.Name) {
		db.Table("history-" + quote.Name).CreateTable(&CoinHistoricalQuote{})
	}
	th := CoinHistoricalQuote{}
	db.Table("history-"+quote.Name).Where("time_open = ?", quote.TimeOpen).First(&th)
	if th.TimeOpen == "" {
		db.Table("history-" + quote.Name).Create(quote)
		fmt.Println(quote, "insert success!")
	} else {
		fmt.Println(quote, "already exists!")
	}
}

func InsertUserInfo(db *gorm.DB, user *UserTable) {
	if !db.HasTable("Users") {
		db.Table("Users").CreateTable(&UserTable{})
	}

	db.Table("Users").Create(user)
	fmt.Println(user, "insert success!")

}
