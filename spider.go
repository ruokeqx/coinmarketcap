package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/bitly/go-simplejson"
	"github.com/jinzhu/gorm"
)

var market_url = "https://api.coinmarketcap.com/data-api/v3/cryptocurrency/market-pairs/latest?slug=%s&start=1&limit=100&category=spot&sort=cmc_rank_advanced"
var chart_url = "https://api.coinmarketcap.com/data-api/v3/cryptocurrency/detail/chart?id=%d&range=%s&convertId=2787"
var historical_url = "https://api.coinmarketcap.com/data-api/v3/cryptocurrency/historical?id=%d&convertId=%d&timeStart=%d&timeEnd=%d" // 2020.01.01

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
		log.Println("http Get failed!")
		return nil
	}
	defer res.Body.Close() // 注册关闭连接
	if res.StatusCode != 200 {
		log.Printf("status code error: %d %s\n", res.StatusCode, res.Status)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("http read body failed!")
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

	if body == nil {
		return 0
	}

	// 手动解析json
	result := make(map[string]interface{})
	json.Unmarshal(body, &result)
	id := result["data"].(map[string]interface{})["id"]
	// fmt.Printf("%v", id)
	return int(id.(float64)) // 对interface{float64}转化为int

}

// 增加天数选项 后续可以定时每天爬1D
func ParserChartData(db *gorm.DB, coin_name string, chart_url string, id int, choice string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Parser Chart Data Error：%s\n", r)
		}
	}()

	url := fmt.Sprintf(chart_url, id, choice)

	js, err := simplejson.NewJson(Download(url))
	if err != nil || js == nil {
		log.Println("something wrong when call NewFromReader")
		return
	}
	// fmt.Println(js)

	points_js := js.Get("data").Get("points").MustMap()
	// points_zh_js := zh_js.Get("data").Get("points").MustMap()
	for t, _ := range points_js {
		var point CoinPointQuote
		point.Id = id
		point.Name = coin_name
		point.Time = t
		point.Price = js.Get("data").Get("points").Get(t).Get("v").GetIndex(0).MustFloat64()
		point.Volume = js.Get("data").Get("points").Get(t).Get("v").GetIndex(1).Interface().(json.Number).String()
		point.MarketCap = js.Get("data").Get("points").Get(t).Get("v").GetIndex(2).Interface().(json.Number).String()
		point.ZhPrice = js.Get("data").Get("points").Get(t).Get("c").GetIndex(0).MustFloat64()
		point.ZhVolume = js.Get("data").Get("points").Get(t).Get("c").GetIndex(1).Interface().(json.Number).String()
		point.ZhMarketCap = js.Get("data").Get("points").Get(t).Get("c").GetIndex(2).Interface().(json.Number).String()
		point.BitcoinRate = js.Get("data").Get("points").Get(t).Get("v").GetIndex(3).Interface().(json.Number).String()
		// fmt.Printf("%v\n", point)
		InsertChart(db, point)
	}
}

// 增加timeStart 可以每多少时间爬指定量 不用每次从头爬
func GetHistoryData(db *gorm.DB, coin_name string, url string, id int, timeStart int64) {
	usd_url := fmt.Sprintf(url, id, 2781, timeStart, time.Now().Unix())
	cny_url := fmt.Sprintf(url, id, 2787, timeStart, time.Now().Unix())

	usd_body := Download(usd_url)
	cny_body := Download(cny_url)

	usd_js, usd_err := simplejson.NewJson(usd_body)
	cny_js, cny_err := simplejson.NewJson(cny_body)
	if usd_err != nil || usd_js == nil || cny_err != nil || cny_js == nil {
		log.Println("something wrong when call NewFromReader")
		return
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
		quote.Id = usd_js.Get("data").Get("id").MustInt()
		// quote.Name = usd_js.Get("data").Get("name").MustString()
		quote.Name = coin_name // 接口解析出来的币名有些不一样 统一都用主页爬到的币名
		quote.Symbol = usd_js.Get("data").Get("symbol").MustString()
		quote.TimeOpen = reg_timeOpen.FindStringSubmatch(quote_usd_js.Get("timeOpen").MustString())[0]
		if quote_usd_js.Get("timeOpen").MustString() != quote_cny_js.Get("timeOpen").MustString() {
			println("Parser USD and CNY Historical Quotes Error!")
			println(quote_usd_js.Get("timeOpen").MustString(), quote_cny_js.Get("timeOpen").MustString())
			// os.Exit(-1)
			return
		}

		// quote.timeClose = quote_usd_js.Get("timeClose").MustString()
		quote.TimeHigh = strings.Replace(reg_timeHL.FindStringSubmatch(quote_usd_js.Get("timeHigh").MustString())[0], "T", " ", 1)
		quote.TimeLow = strings.Replace(reg_timeHL.FindStringSubmatch(quote_usd_js.Get("timeLow").MustString())[0], "T", " ", 1)
		quote.OpenPrice = quote_usd_js.Get("quote").Get("open").MustFloat64()
		quote.ZhOpenPrice = quote_cny_js.Get("quote").Get("open").MustFloat64()
		quote.HighPrice = quote_usd_js.Get("quote").Get("high").MustFloat64()
		quote.ZhHighPrice = quote_cny_js.Get("quote").Get("high").MustFloat64()
		quote.LowPrice = quote_usd_js.Get("quote").Get("low").MustFloat64()
		quote.ZhLowPrice = quote_cny_js.Get("quote").Get("low").MustFloat64()
		quote.ClosePrice = quote_usd_js.Get("quote").Get("close").MustFloat64()
		quote.ZhClosePrice = quote_cny_js.Get("quote").Get("close").MustFloat64()
		quote.Volume = quote_usd_js.Get("quote").Get("volume").Interface().(json.Number).String()
		quote.ZhVolume = quote_cny_js.Get("quote").Get("volume").Interface().(json.Number).String()
		quote.MarketCap = quote_usd_js.Get("quote").Get("marketCap").Interface().(json.Number).String()
		quote.ZhMarketCap = quote_cny_js.Get("quote").Get("marketCap").Interface().(json.Number).String()
		// quote.timestamp = quote_usd_js.Get("quote").Get("timestamp").MustString()
		// fmt.Printf("%v\n", quote)
		InsertHistory(db, quote)
	}
}

// func main() {
// 	spider(int64(3), "7D", int64(1577808000), false)
// }
