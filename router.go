package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/jinzhu/gorm"
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
	json := make(map[string]interface{})
	c.BindJSON(&json)
	username := json["username"].(string)
	password := json["password"].(string)

	tmp_user := User{}
	err := db.Table("Users").Where("username = ?", username).First(&tmp_user).Error
	if err == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "User Registered",
		})
		return
	}
	// V1 基于时间
	// u1, err := uuid.NewUUID()
	// V4 基于随机数
	salt := uuid.New()

	sha256 := sha256.New()
	sha256.Write([]byte(password))
	sha256.Write([]byte(salt.String()))
	pwdhash := hex.EncodeToString(sha256.Sum(nil))

	tmp_user.Username = username
	tmp_user.PwdHash = pwdhash + "$" + salt.String()
	err = InsertUserInfo(db, &tmp_user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "Registry fail!",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"msg":  "Registry Success!",
		})
	}
}

func Login(c *gin.Context) {
	json := make(map[string]interface{})
	c.BindJSON(&json)
	username := json["username"].(string)
	password := json["password"].(string)
	tmp_user := User{}

	err := db.Table("Users").Where("username = ?", username).First(&tmp_user).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "Not Registry!",
		})
	} else {
		var stored_pass []string
		stored_pass = strings.Split(tmp_user.PwdHash, "$")

		sha256 := sha256.New()
		sha256.Write([]byte(password))
		sha256.Write([]byte(stored_pass[1]))
		pwdhash := hex.EncodeToString(sha256.Sum(nil))
		if strings.Compare(pwdhash, stored_pass[0]) != 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": http.StatusBadRequest,
				"msg":  "PassWord Error!",
			})
			return
		} else {
			token, _ := GenerateJwtToken(tmp_user.Uid, tmp_user.Username)
			c.JSON(http.StatusOK, gin.H{
				"code": http.StatusOK,
				"msg":  "Login Success!",
				"data": token,
			})
		}
	}
}

func TokenAuthMiddleware(c *gin.Context) {
	tokenString := c.Request.Header.Get("token")
	// fmt.Println(token)
	if tokenString != "" {
		claims := Claims{}
		_, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(jwtSecret), nil
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": 401,
				"msg":  "Token Auth Failed!",
			})
			return
		}
		c.Set("uid", claims.Uid)
		c.Next()
		return
	}

	c.JSON(http.StatusInternalServerError, gin.H{
		"code": 401,
		"msg":  "Token Auth Failed!",
	})
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

func historical_page_num(c *gin.Context) {
	coinName := c.Query("coinName")
	if !db.HasTable("history-" + coinName) {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "No query coin!",
		})
		return
	}
	var count int
	db.Table("history-" + coinName).Count(&count)
	if count == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "No data now!",
		})
		return
	} else {
		b, err := json.Marshal(count)
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
	coinName := c.Query("coinName")
	pages, _ := strconv.Atoi(c.Query("pages")) // pages=1&limits=1000
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
			c.JSON(http.StatusBadRequest, gin.H{
				"code": http.StatusBadRequest,
				"msg":  "query result convert to json error!",
			})
			return
		} else {
			c.Writer.Write(b)
			return
		}
	}
}

// 获取用户收藏
func like_get(c *gin.Context) {
	uid, _ := c.Get("uid")
	var tmp string
	var coinlike []string
	rows, _ := db.Raw("select name from userlike where uid = ?", uid).Rows()
	for rows.Next() {
		rows.Scan(&tmp)
		coinlike = append(coinlike, tmp)
	}
	b, err := json.Marshal(coinlike)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "query result convert to json error!",
		})
		return
	} else {
		c.Writer.Write(b)
		return
	}
}

// 添加收藏
func like_add(c *gin.Context) {

	if !db.HasTable("CoinLikes") {
		db.Table("CoinLikes").CreateTable(&CoinLike{})
	}

	uid, _ := c.Get("uid")
	cid, _ := strconv.Atoi(c.PostForm("cid"))
	coinlike := CoinLike{
		Uid: uid.(int),
		Cid: cid,
	}
	if err := db.Table("CoinLikes").Create(coinlike).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "like already exists!",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "add like success!",
	})
}

// 删除收藏
func like_del(c *gin.Context) {
	cid := c.Query("cid")
	uid, _ := c.Get("uid")
	db.Table("CoinLikes").Where("uid = ? and cid = ?", uid, cid).Delete(&CoinLike{})
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "delete like success!",
	})
}

// 添加虚拟货币
/*
	{
		"name" :,
		"id":
	}
*/
func addcoin(c *gin.Context) {
	jsonbody := make(map[string]interface{})
	c.BindJSON(&jsonbody)
	name := jsonbody["name"].(string)
	id := jsonbody["id"].(int)
	newcoin := Coin{Name: name, Id: id}

	var numm float64
	err := db.Table("coins").Create(&newcoin).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "fail",
		})
		return
	} else {
		b, _ := json.Marshal(numm)
		c.Writer.Write(b)
		return
	}
}

// 指定时间段内指定货币的价值比
func rate(c *gin.Context) {
	jsonbody := make(map[string]interface{})
	c.BindJSON(&jsonbody)
	st := time.Unix(int64(jsonbody["st"].(float64)), 0).Format("2006-01-02")
	et := time.Unix(int64(jsonbody["et"].(float64)), 0).Format("2006-01-02")
	cid := int(jsonbody["cid"].(float64))

	var numm float64
	rows, err := db.Raw("call db1.rate(?,?,?)", st, et, cid).Rows()
	for rows.Next() {
		rows.Scan(&numm)
	}
	if err != nil {
		fmt.Println("wrong")
		return
	} else {
		b, _ := json.Marshal(numm)
		c.Writer.Write(b)
		return
	}
}

// 指定时间段内指定货币平均开盘价格
func periodavgopen(c *gin.Context) {
	jsonbody := make(map[string]interface{})
	c.BindJSON(&jsonbody)
	st := time.Unix(int64(jsonbody["st"].(float64)), 0).Format("2006-01-02")
	et := time.Unix(int64(jsonbody["et"].(float64)), 0).Format("2006-01-02")
	cid := int(jsonbody["cid"].(float64))

	var numm float64
	rows, err := db.Raw("call db1.periodavgopen(?,?,?)", st, et, cid).Rows()
	for rows.Next() {
		rows.Scan(&numm)
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "wrong!",
		})
		return
	} else {
		b, _ := json.Marshal(numm)
		c.Writer.Write(b)
		return
	}
}

// 虚拟从银行随意获取货币
func getmoney(c *gin.Context) {
	// 修改account表在对应用户cid=-1(人民币) 为对应数值
	json := make(map[string]interface{})
	c.BindJSON(&json)
	num := json["num"].(float64)
	uid, _ := c.Get("uid")
	cid := int(json["cid"].(float64))

	db.AutoMigrate(&Account{})
	ac := Account{}
	acc := Account{uid.(int), cid, num}
	fmt.Println(acc)
	err := db.Table("Accounts").Where("Uid = ? and Cid = ?", uid, cid).First(&ac).Error
	if err != nil {
		db.Table("Accounts").Create(&acc)
	} else {
		acc.Cnum = ac.Cnum + num
		db.Table("Accounts").Where("Uid = ? and Cid = ?", uid, cid).Update(acc)
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "get success!",
	})
}

// 返回当前账户的货币
func account(c *gin.Context) {
	uid, _ := c.Get("uid")

	var results []MyAccount
	db.Raw("select name as coin_name,cnum from coins,accounts where coins.id = accounts.cid and accounts.uid = ?", uid.(int)).Scan(&results)

	b, err := json.Marshal(results)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "query result convert to json error!",
		})
		return
	} else {
		c.Writer.Write(b)
		return
	}

}

// 利用存储过程和游标计算当前用户账户的金额
func myaccount(c *gin.Context) {
	uid, _ := c.Get("uid")

	var numm float64
	var numm2 float64
	rows, err := db.Raw("call db1.accountsum(?)", uid.(int)).Rows()
	for rows.Next() {
		rows.Scan(&numm, &numm2)
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "wrong",
		})
		return
	} else {
		b, _ := json.Marshal([]float64{numm, numm2})
		c.Writer.Write(b)
		return
	}
}

// 创建交易
/*
	{
		"TsCreateTime" : ,
		"ExpectedTime" : ,
		"TsCid" : ,
		"TsNum" : ,
	}
*/
func create_transaction(c *gin.Context) {
	db.AutoMigrate(&Transaction{})
	ts := Transaction{}
	uid, _ := c.Get("uid")

	json := make(map[string]interface{})
	c.BindJSON(&json)
	ExpectedTime := int64(json["ExpectedTime"].(float64))
	TsCid := int(json["TsCid"].(float64))
	TsNum := json["TsNum"].(float64)

	ts.TsStatus = 0
	ts.SellerId = uid.(int)
	ts.TsCreateTime = time.Now().Unix()
	ts.ExpectedTime = ExpectedTime
	ts.TsCid = TsCid
	ts.TsNum = TsNum

	err := db.Table("Transactions").Create(&ts).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "create transaction error!",
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"msg":  "create transaction Success!",
		})
		return
	}
}

// 更新交易
/*
	{
		"TsId" : ,
		"TsCreateTime" : ,
		"ExpectedTime" : ,
		"TsCid" : ,
		"TsNum" : ,
	}
*/
func update_transaction(c *gin.Context) {
	db.AutoMigrate(&Transaction{})
	ts := Transaction{}
	uid, _ := c.Get("uid")

	json := make(map[string]interface{})
	c.BindJSON(&json)
	TsId := int(json["TsId"].(float64))
	TsCreateTime := int64(json["TsCreateTime"].(float64))
	ExpectedTime := int64(json["ExpectedTime"].(float64))
	TsCid := int(json["TsCid"].(float64))
	TsNum := json["TsNum"].(float64)

	ts.TsId = TsId
	ts.TsStatus = 0
	ts.SellerId = uid.(int)
	ts.TsCreateTime = TsCreateTime
	ts.ExpectedTime = ExpectedTime
	ts.TsCid = TsCid
	ts.TsNum = TsNum

	err := db.Table("Transactions").Update(&ts).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "update transaction error!",
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"msg":  "update transaction Success!",
		})
		return
	}
}

func timeout_transaction(c *gin.Context) {
	db.AutoMigrate(&Transaction{})
	tss := []Transaction{}
	uid, _ := c.Get("uid")
	db.Table("Transactions").Where("seller_id = ? and expected_time < ?", uid, time.Now().Unix()).Find(&tss)
	// 插入消息
	msg := Msg{}
	for i, _ := range tss {
		msg.Uid = uid.(int)
		msg.MsgType = 1
		msg.TsId = tss[i].TsId
		db.Table("msgs").Create(&msg)
	}
}

// 交易搜索 default -1 for all
func search_transaction(c *gin.Context) {
	// uid, _ := c.Get("uid")
	db.AutoMigrate(&Transaction{})
	jsonbody := make(map[string]interface{})
	c.BindJSON(&jsonbody)
	cid := int(jsonbody["cid"].(float64))

	tss := []Transaction{}
	if cid == -1 {
		db.Table("transactions").Where("ts_status = 0").Scan(&tss)
	} else {
		db.Table("transactions").Where("ts_status = 0 and ts_cid = ?", cid).Scan(&tss)
	}
	for i, _ := range tss {
		var numm float64
		rows, _ := db.Raw("call db1.coincny(?)", tss[i].TsCid).Rows()
		for rows.Next() {
			rows.Scan(&numm)
		}
		tss[i].Cost = numm * tss[i].TsNum * tss[i].Discount
	}
	b, _ := json.Marshal(tss)
	c.Writer.Write(b)
	return
}

func onetransaction(c *gin.Context) {
	db.AutoMigrate(&Transaction{})
	jsonbody := make(map[string]interface{})
	c.BindJSON(&jsonbody)
	TsId := int(jsonbody["TsId"].(float64))

	ts := Transaction{}
	if err := db.Table("transactions").Where("ts_id = ?", TsId).Scan(&ts).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "no such transaction!",
		})
		return
	}
	var numm float64
	rows, _ := db.Raw("call db1.coincny(?)", ts.TsCid).Rows()
	for rows.Next() {
		rows.Scan(&numm)
	}
	ts.Cost = numm * ts.TsNum * ts.Discount
	b, _ := json.Marshal(ts)
	c.Writer.Write(b)
	return
}

func msg(c *gin.Context) {
	uid, _ := c.Get("uid")
	msgs := []Msg{}

	db.Table("msgs").Where("uid = ?", uid).Find(&msgs)

	b, _ := json.Marshal(msgs)
	c.Writer.Write(b)
	return
}

func readmsg(c *gin.Context) {
	var msgs []int
	c.ShouldBind(&msgs)
	for _, msg := range msgs {
		db.Table("msgs").Where("msg_id = ?", msg).Delete(&Msg{})
	}
}

func close_transaction(c *gin.Context) {
	db.AutoMigrate(&Transaction{})
	uid, _ := c.Get("uid")
	json := make(map[string]interface{})
	c.BindJSON(&json)
	TsId := int(json["TsId"].(float64))

	// 更新状态和关闭时间
	// err := db.Raw("update transactions set ts_status = 2,ts_close_time = ? where seller_id = ? and ts_id = ?", time.Now().Unix(), uid, TsId).Error
	// err := db.Raw("update transactions set ts_status = 2 ,ts_close_time = 123456 where seller_id = 4 and ts_id = 9").Error
	err := db.Table("transactions").Where("seller_id = ? and ts_id = ?", uid, TsId).Updates(map[string]interface{}{"ts_status": 2, "ts_close_time": time.Now().Unix()}).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": http.StatusNotFound,
			"msg":  "update error!",
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"msg":  "update success!",
		})
		return
	}
}

// 我卖的(订单)
func mysell_transaction(c *gin.Context) {
	db.AutoMigrate(&Transaction{})
	uid, _ := c.Get("uid")

	tss := []Transaction{}
	db.Table("transactions").Where("seller_id = ?", uid).Scan(&tss)
	for i, _ := range tss {
		var numm float64
		rows, _ := db.Raw("call db1.coincny(?)", tss[i].TsCid).Rows()
		for rows.Next() {
			rows.Scan(&numm)
		}
		tss[i].Cost = numm * tss[i].TsNum * tss[i].Discount
	}
	b, _ := json.Marshal(tss)
	c.Writer.Write(b)
	return
}

// 打折
func discount(c *gin.Context) {
	db.AutoMigrate(&Transaction{})
	uid, _ := c.Get("uid")

	json := make(map[string]interface{})
	c.BindJSON(&json)
	TsId := int(json["TsId"].(float64))
	discount := json["discount"].(float64)

	if err := db.Table("transactions").Where("seller_id = ? and ts_id = ?", uid, TsId).Update(map[string]interface{}{"discount": discount}).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "fail!",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "success!",
	})
}

// 我买的(订单)
func mybuy_transaction(c *gin.Context) {
	db.AutoMigrate(&Transaction{})
	uid, _ := c.Get("uid")

	tss := []Transaction{}
	db.Table("transactions").Where("buyer_id = ?", uid).Scan(&tss)
	for i, _ := range tss {
		var numm float64
		rows, _ := db.Raw("call db1.coincny(?)", tss[i].TsCid).Rows()
		for rows.Next() {
			rows.Scan(&numm)
		}
		tss[i].Cost = numm * tss[i].TsNum * tss[i].Discount
	}
	b, _ := json.Marshal(tss)
	c.Writer.Write(b)
	return
}

// 购买
/*
	{
		"TsId": ""
	}
*/
func buy(c *gin.Context) {
	db.AutoMigrate(&Transaction{})
	ts := Transaction{}
	sac := Account{}
	bac := Account{}
	svac := Account{}
	bvac := Account{}
	uid, _ := c.Get("uid")
	json := make(map[string]interface{})
	c.BindJSON(&json)
	TsId := int(json["TsId"].(float64))

	err1 := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Table("transactions").Where("ts_id = ?", TsId).Scan(&ts).Error; err != nil {
			return err
		}
		ts.TsStatus = 3
		ts.BuyerId = uid.(int)
		ts.TsCloseTime = time.Now().Unix()
		// 获取币种最新汇率
		var numm float64
		rows, _ := db.Raw("call db1.coincny(?)", ts.TsCid).Rows()
		for rows.Next() {
			rows.Scan(&numm)
		}
		offset := numm * ts.TsNum * ts.Discount
		// 人民币账户
		if err := tx.Table("accounts").Where("uid = ? and cid = -1", ts.SellerId).Scan(&sac).Error; err != nil {
			return err
		}
		if err := tx.Table("accounts").Where("uid = ? and cid = -1", ts.BuyerId).Scan(&bac).Error; err != nil {
			return errors.New("no cny in your account")
		}
		// 虚拟货币账户
		if err := tx.Table("accounts").Where("uid = ? and cid = ?", ts.SellerId, ts.TsCid).Scan(&svac).Error; err != nil {
			svac.Cnum = 0
		}
		if err := tx.Table("accounts").Where("uid = ? and cid = ?", ts.BuyerId, ts.TsCid).Scan(&bvac).Error; err != nil {
			bvac.Cnum = 0
			tx.Table("accounts").Create(Account{uid.(int), ts.TsCid, 0})
		}
		if bac.Cnum-offset < 0 || svac.Cnum-ts.TsNum < 0 {
			return errors.New("Insufficient account balance")
		} else {
			sac.Cnum = sac.Cnum + offset
			bac.Cnum = bac.Cnum - offset
			svac.Cnum = svac.Cnum - ts.TsNum
			bvac.Cnum = bvac.Cnum + ts.TsNum
		}
		if err := tx.Table("accounts").Where("uid = ? and cid = -1", ts.SellerId).Update(&sac).Error; err != nil {
			return err
		}
		if err := tx.Table("accounts").Where("uid = ? and cid = -1", ts.BuyerId).Update(&bac).Error; err != nil {
			return err
		}
		if err := tx.Table("accounts").Where("uid = ? and cid = ?", ts.SellerId, ts.TsCid).Update("cnum", svac.Cnum).Error; err != nil {
			return err
		}
		if err := tx.Table("accounts").Where("uid = ? and cid = ?", ts.BuyerId, ts.TsCid).Update("cnum", bvac.Cnum).Error; err != nil {
			return err
		}
		if err := tx.Table("transactions").Where("ts_id = ?", ts.TsId).Update(&ts).Error; err != nil {
			return err
		}
		// 创建消息 通知卖家
		msg := Msg{}
		msg.Uid = ts.SellerId
		msg.MsgType = 2
		msg.TsId = ts.TsId
		if err := tx.Table("msgs").Create(&msg).Error; err != nil {
			return err
		}
		// 返回 nil 提交事务
		return nil
	})
	if err1 != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  err1.Error(),
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"msg":  "success!",
		})
		return
	}
}
