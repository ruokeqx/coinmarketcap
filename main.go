package main

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// timed spider
func timedSpider() {
	for {
		now := time.Now()
		next := now.Add(time.Hour * 24)
		next = time.Date(next.Year(), next.Month(), next.Day(), 8, 30, 0, 0, next.Location())
		t := time.NewTimer(next.Sub(now))
		<-t.C
		spider(int64(3), "1D", int64(time.Now().Unix()-86400))
	}
}

func main() {
	go func() {
		timedSpider()
	}()

	router := gin.Default()

	// middleware
	router.Use(CORSMiddleware())

	/*
		websocket api
			/price/latest
		param
			[1,1027,1839,52,5994]
		ret
			{ "id":"price","d":{"cr":{"id":1027,"d":18.377500,"p1h":-0.06657870953600,"p24h":0.64696978759100,"p7d":-3.12034341889800,"p30d":null,"ts":null,"as":null,"fmc":null,"mc":353455376506.3070840369470325028557196526469562217300733894480,"mc24hpc":null,"vol24hpc":null,"fmc24hpc":null,"p":3001.9120231238236978753265827306660050872345303020,"v":16843485653.9349384307861328125},"t":1633058777656},"s":"0"}
	*/
	router.GET("/price/latest", latest)
	/*
		api
			/data-api/v3/cryptocurrency/detail/chart?coinName=(?)&range=(?)&convertId=(?)
		param
			coinName	the name of the coin
			range		option:1D 7D 1M 1Y
			convertID	2787 for CNY/2781 for USD
		example
			/data-api/v3/cryptocurrency/detail/chart?coinName=XRP&range=1D&convertId=2787
		ret
			[
				{
					"Id": 52,
					"Name": "XRP",
					"Time": "1632894543",
					"Price": 0.9192637616936671,
					"Volume": "2942241125.47066879",
					"MarketCap": "42975984656.8174226046981975337800000000",
					"BitcoinRate": "0.000021778044",
					"ZhPrice": 5.943315998478,
					"ZhVolume": "19022471548.504901700572",
					"ZhMarketCap": "277852633601.712725230806"
				},
				{
					"Id": 52,
					"Name": "XRP",
					"Time": "1632898143",
					"Price": 0.924008822606362,
					"Volume": "2957768818.47787476",
					"MarketCap": "43197818338.8108602283993584149200000000",
					"BitcoinRate": "0.000021847960",
					"ZhPrice": 5.973532236385,
					"ZhVolume": "19121383857.695313024876",
					"ZhMarketCap": "279265255996.737851842149"
				}
			]
	*/
	router.GET("/data-api/v3/cryptocurrency/detail/chart", chart)
	/*
		api
			/data-api/v3/cryptocurrency/historical?coinName=(?)&timeStart=(?)&timeEnd=(?)
		param
			coinName	the name of the coin
			timeStart	timestamp of the start time
			timeEnd 	timestamp of the end time
		example
			/data-api/v3/cryptocurrency/historical?coinName=XRP&timeStart=1630686243&timeEnd=1630936389
		ret
			[
				{
					"Id": 52,
					"Name": "XRP",
					"Symbol": "XRP",
					"TimeOpen": "2021-09-04",
					"TimeHigh": "2021-09-04 00:24:03",
					"TimeLow": "2021-09-04 17:27:03",
					"OpenPrice": 1.29070149,
					"HighPrice": 1.29591721,
					"LowPrice": 1.24095727,
					"ClosePrice": 1.255779,
					"Volume": "4284360887.7100000000",
					"MarketCap": "58446890874.1300000000",
					"ZhOpenPrice": 8.3295420657,
					"ZhHighPrice": 8.3632017147,
					"ZhLowPrice": 8.0085177419,
					"ZhClosePrice": 8.1041697765,
					"ZhVolume": "27649122988.8366594660",
					"ZhMarketCap": "377187010256.2003350500"
				},
				{
					"Id": 52,
					"Name": "XRP",
					"Symbol": "XRP",
					"TimeOpen": "2021-09-06",
					"TimeHigh": "2021-09-06 21:53:09",
					"TimeLow": "2021-09-06 00:00:03",
					"OpenPrice": 1.30731654,
					"HighPrice": 1.4137667,
					"LowPrice": 1.30731654,
					"ClosePrice": 1.38941299,
					"Volume": "7403887685.1300000000",
					"MarketCap": "64666529504.6300000000",
					"ZhOpenPrice": 8.435982901,
					"ZhHighPrice": 9.1228951384,
					"ZhLowPrice": 8.435982901,
					"ZhClosePrice": 8.9657430832,
					"ZhVolume": "47776546843.3748443393",
					"ZhMarketCap": "417286648240.4222746712"
				}
			]
	*/
	router.GET("/data-api/v3/cryptocurrency/historical", historical)

	if err := router.Run(":8080"); err != nil {
		log.Fatal("failed run app: ", err)
	}
}
