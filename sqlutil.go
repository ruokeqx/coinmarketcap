package main

import (
	"fmt"
	"sync"

	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Coin struct {
	Name string
	Id   int `gorm:"primary_key"`
}
type CoinLike struct {
	Uid int
	Cid int
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

type User struct {
	Uid      int    `json:"uid"`
	Username string `json:"username"`
	PwdHash  string `json:"password"`
}

// Claims 声明
type Claims struct {
	Uid      int    `json:"uid"`
	Username string `json:"loginname"`
	jwt.StandardClaims
}

type Transaction struct {
	TsId         int     `gorm:"primary_key;AUTO_INCREMENT"` // 交易编号
	TsStatus     int     `json:"TsStatus"`                   // 交易状态 0 created/1 timeout/2 closed
	SellerId     int     `json:"SellerId"`                   // 卖家id
	BuyerId      int     `json:"BuyerId"`                    // 买家id
	TsCreateTime int64   `json:"TsCreaTime"`                 // 交易创建时间
	ExpectedTime int64   `json:"ExpectedTime"`               // 期待交易日期
	TsCloseTime  int64   `json:"TsCloseTime"`                // 交易关闭时间
	TsCid        int     `json:"TsCid"`                      // 交易货币
	Discount     float64 `gorm:"default:1;"`                 // 折扣
	TsNum        float64 `json:"TsNum"`                      // 交易量
	Cost         float64 `json:"Cost"`
}

type Msg struct {
	MsgId   int `gorm:"primary_key;AUTO_INCREMENT"`
	Uid     int
	MsgType int // 0 discount 1 time out 2 finished
	TsId    int
}

type Account struct {
	Uid  int     `gorm:"primary_key"`
	Cid  int     `gorm:"primary_key"`
	Cnum float64 `json:"Cnum"`
}

type MyAccount struct {
	CoinName string  `json:"CoinName"`
	Cnum     float64 `json:"Cnum"`
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
		db, err = gorm.Open("mysql", "root:root@(127.0.0.1:3306)/db1?charset=utf8mb4&parseTime=True&loc=Local&multiStatements=true&autocommit=true")
		// rawdb, _ = sql.Open("mysql", "root:root@(127.0.0.1:3306)/db1?charset=utf8mb4&parseTime=True&loc=Local&multiStatements=true&autocommit=true")
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

func InsertUserInfo(db *gorm.DB, user *User) error {
	if !db.HasTable("Users") {
		db.Table("Users").CreateTable(&User{})
	}

	err := db.Table("Users").Create(user).Error
	if err != nil {
		return err
	}
	return nil
}

// func main() {
// 	password := "ruokeqx"
// 	salt := uuid.New()

// 	sha256 := sha256.New()
// 	sha256.Write([]byte(password))
// 	fmt.Println(salt)
// 	sha256.Write([]byte(salt.String()))
// 	fmt.Println(hex.EncodeToString(sha256.Sum(nil)) + "$" + salt.String())

// 	var salted_hash []string
// 	salted_hash_password := "e3a3c3934225b513e03ff0a7f68a605fd8dcf63c1d8c30cde3717c26d9f58935$93d2d0b0-df17-44e3-a875-0aef60457719"
// 	salted_hash = strings.Split(salted_hash_password, "$")
// 	fmt.Println(salted_hash[0], salted_hash[1])
// }
