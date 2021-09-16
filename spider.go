package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"time"
)

var market_url = "https://api.coinmarketcap.com/data-api/v3/cryptocurrency/market-pairs/latest?slug=%s&start=1&limit=100&category=spot&sort=cmc_rank_advanced"
var historical_url = "https://api.coinmarketcap.com/data-api/v3/cryptocurrency/historical?id=%d&convertId=2787&timeStart=1626393600&timeEnd=1631750400"

type Coin struct {
	name  string
	price string
}

var jar, err = cookiejar.New(nil) // 设置全局cookie管理器
func download(url string) {
	println(url)
	client := &http.Client{
		Jar:     jar,              // Jar 域自动管理Cookie
		Timeout: 15 * time.Second, // 设置15秒超时
	}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.108 Safari/537.36")

	res, err := client.Do(req)
	if err != nil {
		log.Fatal("http Get failed!")
		return
	}
	defer res.Body.Close() // 注册关闭连接
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}
	body, err := ioutil.ReadAll(res.Body)

	var result interface{}
	json.Unmarshal(body, &result)

	m := result.(map[string]interface{})
	println(m)
}

func main() {
	url := fmt.Sprintf(market_url, "Bitcoin")
	download(url)
}
