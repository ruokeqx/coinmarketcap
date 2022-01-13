package main

import (
	"encoding/hex"
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
	db := sqlInit()
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
	db := sqlInit()
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

/*
++++++++++++++
+ login part +
++++++++++++++
*/

func Register(c *gin.Context) {
	var mAuth auth

	// 解析 body json 数据到实体类
	if err := c.ShouldBindJSON(&mAuth); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  err,
		})
		return
	}
	db := sqlInit()
	tmp_user := UserTable{}
	db.Table("Users").Where("username = ?", mAuth.UserName).First(&tmp_user)

	// 判断是否存在
	if tmp_user.PwdHash != "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "User Registered",
		})
		return
	}

	pwdhash, err := AesEncrypt([]byte(mAuth.PassWord))
	if err != nil {
		fmt.Print(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  err,
		})
	}

	user_info := UserTable{
		Username: mAuth.UserName,
		PwdHash:  hex.EncodeToString(pwdhash),
	}
	// 注册
	InsertUserInfo(db, &user_info)

	// 注册成功之后 make token
	token, err := GenerateToken(mAuth.UserName)
	if err != nil {
		fmt.Print(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  err,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "Registry Success!",
		"data": token,
	})
}

func Login(c *gin.Context) {
	var mAuth auth
	if err := c.ShouldBindJSON(&mAuth); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  err,
		})
		return
	}

	db := sqlInit()
	tmp_user := UserTable{}
	if db == nil {
		fmt.Println("db nil")
		return
	}
	db.Table("Users").Where("username = ?", mAuth.UserName).First(&tmp_user)

	// 判断是否存在
	if tmp_user.PwdHash == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "User not Registered!",
		})
		return
	}
	pwd, _ := hex.DecodeString(tmp_user.PwdHash)
	pwd, err := AesDecrypt(pwd)
	if err != nil {
		fmt.Print(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  err,
		})
		return
	}

	// 登录失败
	if string(pwd) != mAuth.PassWord {
		fmt.Printf("Login Error:%s %s", string(pwd), mAuth.PassWord)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "PassWord Error!",
		})
		return
	}

	// 生成token
	token, merr := GenerateToken(tmp_user.Username)
	if merr != nil {
		fmt.Print(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "GenerateToken Error!",
		})
		return
	}

	// TokenList.PushBack(token) // 将生成的token存入TokenList中
	db.Table("Users").Where("username = ?", mAuth.UserName).Update("token", token)
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "Login Success!",
		"data": token,
	})
}

func TokenAuthMiddleware(c *gin.Context) {

	tmp_user := UserTable{}
	fmt.Println("TokenAuthMiddleware")

	token := c.Request.Header.Get("token") // 查找请求中是否有token
	fmt.Println(token)
	if token != "" {
		// for i := TokenList.Front(); i != nil; i = i.Next() {
		// 	if i.Value == token {
		// 		fmt.Println("Token Auth Success!")
		// 		c.Next()
		// 		return
		// 	}
		// }

		// 查询是否有用户的token是这个
		db.Table("Users").Where("token = ?", token).First(&tmp_user)
		if tmp_user.Token != "" {
			c.Next()
			return
		}
	}

	c.JSON(http.StatusInternalServerError, gin.H{
		"code": 401,
		"msg":  "Token Auth Failed!",
	})
	// Pass on to the next-in-chain
	c.Abort()
}

/*
++++++++++++++
+    page    +
++++++++++++++
*/

// example //
// func getgoods(c *gin.Context) {
// 	var good []goods

// 	name := c.Query("name")
// 	pages, _ := strconv.Atoi(c.Query("pages"))
// 	limits, _ := strconv.Atoi(c.Query("limits"))
// 	db := sqlInit()

// 	if !db.HasTable("goods") {
// 		db.Table("goods").CreateTable(&goods{})
// 	}
// 	db.Table("goods").Where("name Like ?", "%"+name+"%").Offset((pages - 1) * limits).Limit(limits).Find(&good)
// 	if len(good) == 0 {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"code": http.StatusBadRequest,
// 			"msg":  "No data now!",
// 		})
// 		return
// 	} else {
// 		b, err := json.Marshal(good)
// 		if err != nil {
// 			fmt.Println("query result convert to json error!")
// 			return
// 		} else {
// 			c.Writer.Write(b)
// 			return
// 		}
// 	}
// }

// data to paint chart
func chart_page(c *gin.Context) {
	// /data-api/v3/cryptocurrency/chart_page
	db := sqlInit()

	coinName := c.Query("coinName")
	pages, _ := strconv.Atoi(c.Query("pages"))
	limits, _ := strconv.Atoi(c.Query("limits"))
	if !db.HasTable("chart-" + coinName) {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "No query coin!",
		})
		return
	}
	var chart []CoinPointQuote
	db.Table("chart-" + coinName).Order("time").Offset((pages - 1) * limits).Limit(limits).Find(&chart)
	if len(chart) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "No data now!",
		})
		return
	} else {
		b, err := json.Marshal(chart)
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
func historical_page(c *gin.Context) {
	db := sqlInit()
	coinName := c.Query("coinName")
	pages, _ := strconv.Atoi(c.Query("pages"))
	limits, _ := strconv.Atoi(c.Query("limits"))
	if !db.HasTable("history-" + coinName) {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "No query coin!",
		})
		return
	}
	var his []CoinHistoricalQuote
	db.Table("history-" + coinName).Order("time_open").Offset((pages - 1) * limits).Limit(limits).Find(&his)
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
			return
		} else {
			c.Writer.Write(b)
			return
		}
	}
}

// 获取用户收藏
func like_get(c *gin.Context) {
	db := sqlInit()
	// 根据token获取用户名
	token := c.Request.Header.Get("token")
	tmp_user := UserTable{}
	db.Table("Users").Where("token = ?", token).First(&tmp_user)
	username := tmp_user.Username

	// 查找用户所有收藏的币
	var coinlike []CoinLike
	db.Table("CoinLikes").Where("username = ?", username).Find(&coinlike)
	if len(coinlike) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "No data now!",
		})
		return
	} else {
		b, err := json.Marshal(coinlike)
		if err != nil {
			fmt.Println("query result convert to json error!")
		} else {
			c.Writer.Write(b)
			return
		}
	}
}

// 添加收藏
func like_add(c *gin.Context) {

	if !db.HasTable("CoinLikes") {
		db.Table("CoinLikes").CreateTable(&CoinLike{})
	}

	token := c.Request.Header.Get("token")
	tmp_user := UserTable{}
	db.Table("Users").Where("token = ?", token).First(&tmp_user)
	username := tmp_user.Username

	coinlike := CoinLike{
		Username: username,
		Coinname: c.PostForm("coinName"),
	}
	if err := db.Table("CoinLikes").Create(coinlike).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "like already exists!",
		})
	}

	// fmt.Println(coinlike, "insert success!")

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "add like success!",
	})
}

// 删除收藏
func like_del(c *gin.Context) {
	coinName := c.Query("coinName")
	token := c.Request.Header.Get("token")
	tmp_user := UserTable{}
	db.Table("Users").Where("token = ?", token).First(&tmp_user)
	if coinName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "no coinname!",
		})
		return
	}
	fmt.Println(token, coinName)
	db.Table("CoinLikes").Where("username = ? and coinname = ?", tmp_user.Username, coinName).Delete(&CoinLike{})
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "delete like success!",
	})
}
