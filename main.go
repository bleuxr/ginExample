package main

import (
	"ginExample/common"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	InitConfig()
	common.InitDB()
	// common.InitDB()
	// db.Close() is removed from v2
	// defer db.Close()

	r := gin.Default()
	r = CollectRouter(r)
	port := viper.GetString("server.port")
	if port != "" {
		panic(r.Run(":" + port))
	}
	panic(r.Run())
	// r.Run() // listen and serve on 0.0.0.0:8080
}

func InitConfig() {
	workDir, _ := os.Getwd()
	viper.SetConfigName("application")
	viper.SetConfigType("yml")
	viper.AddConfigPath(workDir + "/config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}
