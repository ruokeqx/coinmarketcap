package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, Content-Length, X-CSRF-Token, Token, session, Origin, Host, Connection, Accept-Encoding, Accept-Language, X-Requested-With")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Request.Header.Del("Origin")

		c.Next()
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// latest data api function
func latest(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  err,
		})
		// log.Println("upgrade:", err)
		return
	}
	defer ws.Close()
	var message []byte
	for {
		_, message, err = ws.ReadMessage()
		fmt.Println(string(message))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": http.StatusInternalServerError,
				"msg":  err,
			})
			// log.Println("read:", err)
			return
		}
		if message != nil && string(message[0]) == "[" {
			break
		}
	}
	// coinmarketcap websocket
	client, _, err := websocket.DefaultDialer.Dial("wss://stream.coinmarketcap.com/price/latest", nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  err,
		})
		// log.Println("websocket:", err)
		return
	}
	// command := "{\"method\":\"subscribe\",\"id\":\"price\",\"data\":{\"cryptoIds\":[1,52,1027,5994],\"index\":null}}"
	command := fmt.Sprintf("{\"method\":\"subscribe\",\"id\":\"price\",\"data\":{\"cryptoIds\":%s,\"index\":null}}", message)
	err = client.WriteMessage(websocket.TextMessage, []byte(command))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  err,
		})
		// log.Println("Subscription failed:", err)
		return
	}
	for {
		_, json_bytes, _ := client.ReadMessage()
		fmt.Println(string(json_bytes))
		err = ws.WriteMessage(websocket.TextMessage, json_bytes) // websocket.TextMessage/websocket.BinaryMessage
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": http.StatusInternalServerError,
				"msg":  err,
			})
			return
		}
	}
}

// data to paint chart
func chart(c *gin.Context) {
	// /data-api/v3/cryptocurrency/detail/chart?coinName=(?)&range=(?)&convertId=(?)
	db, err := sqlInit()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "Database connect error!",
		})
		return
	}
	defer db.Close()
	coinName := c.Query("coinName")
	if !db.HasTable("chart-" + coinName) {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "No query coin!",
		})
		return
	}
	queryRange := c.Query("range")
	convertId := c.Query("convertId")
	fmt.Println(convertId)
	var tE int64
	var tS int64
	var chart []CoinPointQuote
	var retChart []CoinPointQuote
	switch queryRange {
	case "1D":
		// fmt.Println("1D")
		tE = time.Now().Unix()
		tS = tE - 60*60*24
		db.Table("chart-"+coinName).Where("time BETWEEN ? AND ?", tS, tE).Order("time").Find(&chart)
		retChart = chart
	case "7D":
		// fmt.Println("7D")
		tE = time.Now().Unix()
		tS = tE - 60*60*24*7
		db.Table("chart-"+coinName).Where("time BETWEEN ? AND ?", tS, tE).Order("time").Find(&chart)
		retChart = chart
	case "1M":
		// fmt.Println("1M")
		tE = time.Now().Unix()
		tS = tE - 60*60*24*30
		db.Table("chart-"+coinName).Where("time BETWEEN ? AND ?", tS, tE).Order("time").Find(&chart)
		for index, cont := range chart {
			if index%12 == 0 {
				retChart = append(retChart, cont)
			}
		}
	case "1Y":
		// fmt.Println("1Y")
		tE = time.Now().Unix()
		tS = tE - 60*60*24*365
		db.Table("chart-"+coinName).Where("time BETWEEN ? AND ?", tS, tE).Order("time").Find(&chart)
		for index, cont := range chart {
			if index%288 == 0 {
				retChart = append(retChart, cont)
			}
		}
	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "Now such option!",
		})
		return
	}
	if len(retChart) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "No data now!",
		})
		return
	} else {
		b, err := json.Marshal(retChart)
		if err != nil {
			fmt.Println("query result convert to json error!")
			return
		} else {
			c.Writer.Write(b)
			return
		}
	}
}

// historical data
func historical(c *gin.Context) {
	db, err := sqlInit()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "Database connect error!",
		})
		return
	}
	defer db.Close()
	coinName := c.Query("coinName")
	if !db.HasTable("history-" + coinName) {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "No query coin!",
		})
		return
	}
	tS := c.Query("timeStart")
	tE := c.Query("timeEnd")
	timeStart, _ := strconv.ParseInt(tS, 10, 64)
	timeEnd, _ := strconv.ParseInt(tE, 10, 64)
	tS = time.Unix(int64(timeStart), 0).Format("2006-01-02")
	tE = time.Unix(int64(timeEnd), 0).Format("2006-01-02")
	var his []CoinHistoricalQuote
	// var hisList []string
	db.Table("history-"+coinName).Where("time_open BETWEEN ? AND ?", tS, tE).Order("time_open").Find(&his)
	if len(his) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "No data now!",
		})
		return
	} else {
		b, err := json.Marshal(his)
		if err != nil {
			fmt.Println("query result convert to json error!")
		} else {
			c.Writer.Write(b)
			return
		}
	}
}
