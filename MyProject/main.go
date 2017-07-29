package main

import (
	c "MyProject/Controller"
	"MyProject/Modules"
	"gopkg.in/mgo.v2"
	"github.com/gin-gonic/gin"
	"log"
)

var err error

func main() {
	//初始化数据库
	Modules.DbSession, err = mgo.Dial("127.0.0.1")
	if err != nil {
		log.Println(err)
		panic(err)
	}
	//gin配置
	r := gin.Default()
	authorized := r.Group("/", gin.BasicAuth(gin.Accounts{
		"yangshaobo1996@gmail.com": "1+1=3YSB",
	}))
	{
		r.LoadHTMLGlob("src/MyProject/templates/*")
		r.Static("/statics", "./src/MyProject/statics")
		authorized.GET("/", c.GetIndex)
		authorized.POST("/new", c.UpLoadNewFile)
		authorized.GET("/success-upload", c.SuccessUpload)
		download := authorized.Group("/download")
		{
			download.GET("/:md5/:name", c.Download)
			download.GET("", c.DownloadList)
		}

	}
	r.Run(":8000")
}
