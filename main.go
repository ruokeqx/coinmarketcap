package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// middleware
	router.Use(CORSMiddleware())

	// /price/latest
	// param
	// [1,1027,1839,52,5994]
	// ret
	//
	router.GET("/price/latest", latest)
	// /data-api/v3/cryptocurrency/detail/chart?coinName=(?)&range=(?)&convertId=(?)
	router.GET("/data-api/v3/cryptocurrency/detail/chart", chart)
	// /data-api/v3/cryptocurrency/historical?coinName=(?)&timeStart=(?)&timeEnd=(?)
	router.GET("/data-api/v3/cryptocurrency/historical", historical)

	if err := router.Run(":8080"); err != nil {
		log.Fatal("failed run app: ", err)
	}
}
