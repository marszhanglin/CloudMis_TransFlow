// main
package main

import (
	"CloudMis_TransFlow/fconf"
	//"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	//https://github.com/go-playground/validator/blob/v9/_examples/simple/main.go
	"gopkg.in/go-playground/validator.v9"
)

var (
	db          *gorm.DB
	e           *echo.Echo
	projectName string
	validate    *validator.Validate
	iniData     *fconf.Config
)

func main() {
	//initIni()
	projectName = "cloudmis/transflow"
	validate = validator.New()
	db = initDB()
	defer db.Close()
	initLog()
	defer glogFlush()

	initRedis()

	initEcho()
	echoRoute()
	glogFlush()

	//e.Logger.Fatal(e.Start(":" + iniData.String("flow.cnf.webPort")))
	e.Logger.Fatal(e.Start(":10001"))
	//e.Logger.Fatal(e.StartAutoTLS(":"))
}

func initDB() (db *gorm.DB) { 
	db, err := gorm.Open("mysql", "root:password@tcp(localhost:8066)/cloudmis_transflow_mc?charset=utf8&parseTime=True&loc=Local")

	if err != nil {
		panic(err)
	}

	db.DB().SetMaxIdleConns(100)
	db.DB().SetMaxOpenConns(1000)
	db.DB().SetConnMaxLifetime(1000)
	db.LogMode(true)
	return db
}

func initEcho() {
	//e := initEcho()
	e = echo.New()
	// Middleware
	e.Use(middleware.Logger()) //打印到控制台
	//打印到文件
	e.Use(middleware.Recover())
}

func getTimeUUID() int64 {
	uuid := time.Now().Unix()*100000 + rand.Int63n(100000)
	fmt.Println(uuid)
	return uuid
}
