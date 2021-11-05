package main

import (
	"ginExample/common"

	"github.com/gin-gonic/gin"
)

func main() {
	common.InitDB()
	// common.InitDB()
	// db.Close() is removed from v2
	// defer db.Close()

	r := gin.Default()
	r = CollectRouter(r)

	r.Run() // listen and serve on 0.0.0.0:8080
}
