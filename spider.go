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
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"golang.org/x/sync/semaphore"

	"github.com/bitly/go-simplejson"
)

var market_url = "https://api.coinmarketcap.com/data-api/v3/cryptocurrency/market-pairs/latest?slug=%s&start=1&limit=100&category=spot&sort=cmc_rank_advanced"
var historical_url = "https://api.coinmarketcap.com/data-api/v3/cryptocurrency/historical?id=%d&convertId=2787&timeStart=1626393600&timeEnd=1631750400"

type CoinQuote struct {
	id         int
	name       string
	symbol     string
	timeOpen   string
	timeClose  string
	timeHigh   string
	timeLow    string
	openPrice  float64
	highPrice  float64
	lowPrice   float64
	closePrice float64
	volume     float64
	marketCap  float64
	timestamp  string
}

var jar, err = cookiejar.New(nil) // 设置全局cookie管理器
func download(url string) []byte {
	fmt.Println(url)
	client := &http.Client{
		Jar:     jar,              // Jar 域自动管理Cookie
		Timeout: 15 * time.Second, // 设置15秒超时
	}
	req, err := http.NewRequest("GET", url, nil)
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

func ParserHistoryData(body []byte) {

	js, err := simplejson.NewJson(body)
	if err != nil || js == nil {
		log.Fatal("something wrong when call NewFromReader")
	}
	// fmt.Println(js)
	quotes_js := js.Get("data").Get("quotes").MustArray()
	for i, _ := range quotes_js {

		quote_js := js.Get("data").Get("quotes").GetIndex(i)

		var quote CoinQuote
		quote.id = js.Get("data").Get("id").MustInt()
		quote.name = js.Get("data").Get("name").MustString()
		quote.symbol = js.Get("data").Get("symbol").MustString()
		quote.timeOpen = quote_js.Get("timeOpen").MustString()
		quote.timeClose = quote_js.Get("timeClose").MustString()
		quote.timeHigh = quote_js.Get("timeHigh").MustString()
		quote.timeLow = quote_js.Get("timeLow").MustString()
		quote.openPrice = quote_js.Get("quote").Get("open").MustFloat64()
		quote.highPrice = quote_js.Get("quote").Get("high").MustFloat64()
		quote.lowPrice = quote_js.Get("quote").Get("low").MustFloat64()
		quote.closePrice = quote_js.Get("quote").Get("close").MustFloat64()
		quote.volume = quote_js.Get("quote").Get("volume").MustFloat64()
		quote.marketCap = quote_js.Get("quote").Get("marketCap").MustFloat64()
		quote.timestamp = quote_js.Get("quote").Get("timestamp").MustString()
		fmt.Printf("%v\n", quote)
	}
}

func main() {

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

	s := semaphore.NewWeighted(1) // 并发限制为3
	var w sync.WaitGroup
	for coin := coins.Front(); coin != nil; coin = coin.Next() {
		w.Add(1)
		go func(coin_name string) {
			s.Acquire(context.Background(), 1)
			url := fmt.Sprintf(market_url, coin_name)
			id := GetId(download(url))
			if id == 0 {
				return
			}

			url = fmt.Sprintf(historical_url, id)
			ParserHistoryData(download(url))

			s.Release(1)
			w.Done()
		}(coin.Value.(string))
	}
	runtime.Gosched() // 使当前goroutine让出执行时机
	time.Sleep(time.Hour)
}
