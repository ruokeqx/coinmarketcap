package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type auth struct {
	// UserName 用户名
	UserName string `json:"username" example:"zhangsan" validate:"required,gte=5,lte=30"`
	// PassWord 密码
	PassWord string `json:"password" example:"zhangsan" validate:"required,gte=5,lte=30"`
}

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
	var tmp_user UserTable
	db.Where("name = ?", mAuth.UserName).First(&tmp_user)

	// 判断是否存在
	if tmp_user.PwdHash != "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "User Registered",
		})
		return
	}

	pwdhash, err := Encrypt(mAuth.PassWord)
	if err != nil {
		fmt.Print(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  err,
		})
	}

	user_info := UserTable{
		Username: mAuth.UserName,
		PwdHash:  pwdhash,
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
	var tmp_user UserTable
	db.Where("name = ?", mAuth.UserName).First(&tmp_user)

	// 判断是否存在
	if tmp_user.PwdHash == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "User not Registered!",
		})
		return
	}

	pwdhash, err := Encrypt(mAuth.PassWord)
	if err != nil {
		fmt.Print(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  err,
		})
	}

	// 登录失败
	if pwdhash != tmp_user.PwdHash {
		fmt.Printf("Login Error:%s %s", pwdhash, tmp_user.PwdHash)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "PassWord Error!",
		})
	}

	// 生成token
	token, merr := GenerateToken(tmp_user.Username)
	if merr != nil {
		fmt.Print(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "GenerateToken Error!",
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "Login Success!",
		"data": token,
	})
}

var jwtSecret = []byte("JwtSecret")

// Claims 声明
type Claims struct {
	LoginName []byte `json:"loginname"`
	jwt.StandardClaims
}

// GenerateToken 生成 token
func GenerateToken(loginName string) (string, error) {
	var err error
	aesLoginName, err := AesEncrypt([]byte(loginName))
	if err != nil {
		return "", err
	}

	// 现在的时间
	nowTime := time.Now()
	// 过期的时间
	expireTime := nowTime.Add(3 * time.Hour)
	// 初始化 声明
	claims := Claims{
		aesLoginName, jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "aims",
		},
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 获取完整签名之后的 token
	return tokenClaims.SignedString(jwtSecret)
}

func main() {

	// gin server
	router := gin.Default()

	// middleware
	router.Use(CORSMiddleware())
	router.GET("/price/latest", latest)
	router.GET("/data-api/v3/cryptocurrency/detail/chart", chart)
	router.GET("/data-api/v3/cryptocurrency/historical", historical)

	router.POST("/register", Register)
	router.POST("/login", Login)

	if err := router.Run(":8080"); err != nil {
		log.Fatal("failed run app: ", err)
	}
}
