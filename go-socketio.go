package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
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
	db.Table("history-"+coinName).Where("time_open BETWEEN ? AND ?", tS, tE).Find(&his)
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

func main() {
	router := gin.Default()

	server := socketio.NewServer(nil)
	server.OnConnect("/", func(s socketio.Conn) error {
		s.Emit("reply", "hello")
		return nil
	})

	server.OnEvent("/", "notice", func(s socketio.Conn, msg string) {
		fmt.Print(msg)
		s.Emit("res", "have "+msg)
	})

	go server.Serve()
	defer server.Close()

	router.Use(CORSMiddleware())
	router.GET("/socket.io/*any", gin.WrapH(server))
	router.POST("/socket.io/*any", gin.WrapH(server))
	// /data-api/v3/cryptocurrency/detail/chart?coinName=(?)&range=(?)&convertId=(?)
	router.GET("/data-api/v3/cryptocurrency/detail/chart", chart)
	// /data-api/v3/cryptocurrency/historical?coinName=(?)&timeStart=(?)&timeEnd=(?)
	router.GET("/data-api/v3/cryptocurrency/historical", historical)
	if err := router.Run(":8080"); err != nil {
		log.Fatal("failed run app: ", err)
	}
}
