package main

import (
	"bufio"
	"container/list"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"golang.org/x/sync/semaphore"

	"github.com/bitly/go-simplejson"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var market_url = "https://api.coinmarketcap.com/data-api/v3/cryptocurrency/market-pairs/latest?slug=%s&start=1&limit=100&category=spot&sort=cmc_rank_advanced"
var chart_url = "https://api.coinmarketcap.com/data-api/v3/cryptocurrency/detail/chart?id=%d&range=1D"
var historical_url = "https://api.coinmarketcap.com/data-api/v3/cryptocurrency/historical?id=%d&convertId=%d&timeStart=1626393600&timeEnd=1631750400"

type Coin struct {
	Name string
	Id   int
}

type CoinPointQuote struct {
	name        string
	time        string
	price       float64
	volume      string
	bitcoinRate float64
}

type CoinHistoricalQuote struct {
	id       int
	name     string
	symbol   string
	timeOpen string
	// timeClose    string
	timeHigh     string
	timeLow      string
	openPrice    float64
	highPrice    float64
	lowPrice     float64
	closePrice   float64
	volume       string
	marketCap    string
	zhOpenPrice  float64
	zhHighPrice  float64
	zhLowPrice   float64
	zhClosePrice float64
	zhVolume     string
	zhMarketCap  string

	// timestamp string
}

var jar, _ = cookiejar.New(nil) // 设置全局cookie管理器
func Download(tourl string) []byte {
	fmt.Println(tourl)
	// proxy 不用的话就注释掉
	proxy := func(_ *http.Request) (*url.URL, error) {
		return url.Parse("http://127.0.0.1:8088")
	}
	transport := &http.Transport{Proxy: proxy}

	client := &http.Client{
		Transport: transport,
		Jar:       jar,              // Jar 域自动管理Cookie
		Timeout:   15 * time.Second, // 设置15秒超时
	}
	req, _ := http.NewRequest("GET", tourl, nil)
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.108 Safari/537.36")

	res, err := client.Do(req)
	if err != nil {
		log.Fatal("http Get failed!")
		return nil
	}
	defer res.Body.Close() // 注册关闭连接
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal("http read body failed!")
		return nil
	}
	// fmt.Println(string(body))
	return body
}

func GetId(body []byte) int {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("获取ID时捕获到的错误：%s\n", r)
		}
	}()

	// 手动解析json
	result := make(map[string]interface{})
	json.Unmarshal(body, &result)
	id := result["data"].(map[string]interface{})["id"]
	// fmt.Printf("%v", id)
	return int(id.(float64)) // 对interface{float64}转化为int

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

func ParserChartData(coin_name string, chart_url string, id int) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Parser Chart Data Error：%s\n", r)
		}
	}()

	url := fmt.Sprintf(chart_url, id)
	zh_url := url + "&convertId=2787"
	println(url, zh_url)
	js, err := simplejson.NewJson(Download(url))
	zh_js, zh_err := simplejson.NewJson(Download(zh_url))
	if err != nil || js == nil || zh_err != nil || zh_js == nil {
		log.Fatal("something wrong when call NewFromReader")
	}
	// fmt.Println(js)

	points_js := js.Get("data").Get("points").MustMap()
	// points_zh_js := zh_js.Get("data").Get("points").MustMap()
	for t, _ := range points_js {
		var point CoinPointQuote
		point.name = coin_name
		point.time = t
		point.price = zh_js.Get("data").Get("points").Get(t).Get("c").GetIndex(0).MustFloat64()
		point.volume = zh_js.Get("data").Get("points").Get(t).Get("c").GetIndex(1).Interface().(json.Number).String()
		point.bitcoinRate = js.Get("data").Get("points").Get(t).Get("v").GetIndex(3).MustFloat64()
		fmt.Printf("%v\n", point)
	}
}

func GetHistoryData(url string, id int) {
	usd_url := fmt.Sprintf(url, id, 2781)
	cny_url := fmt.Sprintf(url, id, 2787)

	usd_body := Download(usd_url)
	cny_body := Download(cny_url)

	usd_js, usd_err := simplejson.NewJson(usd_body)
	cny_js, cny_err := simplejson.NewJson(cny_body)
	if usd_err != nil || usd_js == nil || cny_err != nil || cny_js == nil {
		log.Fatal("something wrong when call NewFromReader")
	}
	// fmt.Println(js)
	quotes_js := usd_js.Get("data").Get("quotes").MustArray()
	for i, _ := range quotes_js {
		var quote CoinHistoricalQuote
		reg_timeOpen := regexp.MustCompile("[0-9]*-[0-9]*-[0-9]*")
		reg_timeHL := regexp.MustCompile("[^.]*")
		if reg_timeOpen == nil || reg_timeHL == nil {
			fmt.Println("regexp Error!")
			return
		}

		quote_usd_js := usd_js.Get("data").Get("quotes").GetIndex(i)
		quote_cny_js := cny_js.Get("data").Get("quotes").GetIndex(i)
		quote.id = usd_js.Get("data").Get("id").MustInt()
		quote.name = usd_js.Get("data").Get("name").MustString()
		quote.symbol = usd_js.Get("data").Get("symbol").MustString()
		quote.timeOpen = reg_timeOpen.FindStringSubmatch(quote_usd_js.Get("timeOpen").MustString())[0]
		if quote_usd_js.Get("timeOpen").MustString() != quote_cny_js.Get("timeOpen").MustString() {
			println("Parser USD and CNY Historical Quotes Error!")
			println(quote_usd_js.Get("timeOpen").MustString(), quote_cny_js.Get("timeOpen").MustString())
			os.Exit(-1)
		}

		// quote.timeClose = quote_usd_js.Get("timeClose").MustString()
		quote.timeHigh = strings.Replace(reg_timeHL.FindStringSubmatch(quote_usd_js.Get("timeHigh").MustString())[0], "T", " ", 1)
		quote.timeLow = strings.Replace(reg_timeHL.FindStringSubmatch(quote_usd_js.Get("timeLow").MustString())[0], "T", " ", 1)
		quote.openPrice = quote_usd_js.Get("quote").Get("open").MustFloat64()
		quote.zhOpenPrice = quote_cny_js.Get("quote").Get("open").MustFloat64()
		quote.highPrice = quote_usd_js.Get("quote").Get("high").MustFloat64()
		quote.zhHighPrice = quote_cny_js.Get("quote").Get("high").MustFloat64()
		quote.lowPrice = quote_usd_js.Get("quote").Get("low").MustFloat64()
		quote.zhLowPrice = quote_cny_js.Get("quote").Get("low").MustFloat64()
		quote.closePrice = quote_usd_js.Get("quote").Get("close").MustFloat64()
		quote.zhClosePrice = quote_cny_js.Get("quote").Get("close").MustFloat64()
		quote.volume = quote_usd_js.Get("quote").Get("volume").Interface().(json.Number).String()
		quote.zhVolume = quote_cny_js.Get("quote").Get("volume").Interface().(json.Number).String()
		quote.marketCap = quote_usd_js.Get("quote").Get("marketCap").Interface().(json.Number).String()
		quote.zhMarketCap = quote_cny_js.Get("quote").Get("marketCap").Interface().(json.Number).String()
		// quote.timestamp = quote_usd_js.Get("quote").Get("timestamp").MustString()
		fmt.Printf("%v\n", quote)
	}
}

func main() {
	// 创建数据库连接
	db, err := gorm.Open("mysql", "ruokeqx:ruokeqx666@(121.196.208.97:3306)/ruokeqx?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		fmt.Print("Connect database error!")
		return
	}
	defer db.Close()

	// 从coins.txt中读取coin名称并保存到列表
	var coins = list.New()
	file, err := os.Open("./coins.txt") // Open用于读取文件  默认具有Read的文件描述符
	if err != nil {
		fmt.Println("File Open Error:%v\n", err)
		return
	}
	defer file.Close() //滞后关闭
	reader := bufio.NewReader(file)
	for {
		coin_name, err := reader.ReadString('\n') // 读到一个换行就结束
		if err == io.EOF {
			break
		}
		coin_name = strings.Trim(coin_name, "\r\n") // 去除前后换行符,这里巨坑
		// coin_name = strings.Replace(coin_name, "\n", "", 1)
		// fmt.Println(coin_name)
		coins.PushBack(coin_name)
	}
	fmt.Println("File Read Success")

	// 遍历列表  并发获取id和数据
	s := semaphore.NewWeighted(1) // 并发限制为3
	var w sync.WaitGroup          // 等待组
	for coin := coins.Front(); coin != nil; coin = coin.Next() {
		w.Add(1) // 每启动一个新任务  等待组加一
		go func(coin_name string) {
			// 获取id
			s.Acquire(context.Background(), 1)
			url := fmt.Sprintf(market_url, coin_name)
			id := GetId(Download(url))
			if id == 0 {
				s.Release(1)
				w.Done()
				return
			}

			//存储coin_name及其对应id
			InsertCoin(db, coin_name, id)

			// 获取图表数据
			ParserChartData(coin_name, chart_url, id)
			// os.Exit(0)

			// 获取历史数据
			GetHistoryData(historical_url, id)

			s.Release(1) // 释放信号量锁
			w.Done()     // 设置等待组完成一项任务
		}(coin.Value.(string))
	}
	w.Wait() // 等待所有任务的完成  即计数器值为0
}
