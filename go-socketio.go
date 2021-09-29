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

func historical(c *gin.Context) {
	db, err := sqlInit()
	if err != nil {
		fmt.Println("connect error")
	}
	defer db.Close()
	coinName := c.Query("coinName")
	tS := c.Query("timeStart")
	tE := c.Query("timeEnd")
	timeStart, err := strconv.ParseInt(tS, 10, 64)
	timeEnd, _ := strconv.ParseInt(tE, 10, 64)
	tS = time.Unix(int64(timeStart), 0).Format("2006-01-02")
	tE = time.Unix(int64(timeEnd), 0).Format("2006-01-02")
	var his []CoinHistoricalQuote
	// var hisList []string
	db.Table("history-"+coinName).Where("time_open BETWEEN ? AND ?", tS, tE).Find(&his)
	b, err := json.Marshal(his)
	if err != nil {
		fmt.Println("query result convert to json error!")
	} else {
		c.Writer.Write(b)
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

	router.GET("/data-api/v3/cryptocurrency/historical", historical)
	if err := router.Run(":8080"); err != nil {
		log.Fatal("failed run app: ", err)
	}
}
