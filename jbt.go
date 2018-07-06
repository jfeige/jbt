package main

import (
	"gopkg.in/gin-gonic/gin.v1"
	log "github.com/alecthomas/log4go"
	"flag"
	"runtime"
	"net/http"
	"fmt"
	"jbt/models"
	cont "jbt/controllers"
)

var(
	logfile = flag.String("log","./conf/jbt-log.xml","log4go file path!")
	configfile = flag.String("config","./conf/jbt.ini","config file path")
)

func init(){
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main(){
	log.LoadConfiguration(*logfile)
	defer log.Close()

	err := models.InitBaseConfig(*configfile)
	if err != nil{
		fmt.Printf("InitBaseConfig has error:%v\n",err)
		log.Error("InitBaseConfig has error:%v\n", err)
		return
	}

	gin.SetMode(gin.DebugMode) //全局设置环境，此为开发环境，线上环境为gin.ReleaseMode

	router := initRouter()

	err = http.ListenAndServe(models.AppPort, router)
	if err != nil {
		fmt.Printf("http.ListenAndServe has error:%v\n", err)
		log.Error("http.ListenAndServe has error:%v\n", err)
		return
	}
}

func initRouter() *gin.Engine {

	router := gin.Default()
	router.LoadHTMLGlob("views/*")

	router.Static("/static","./static")


	//首页
	router.GET("/",cont.Index)

	//搜索 form提交
	router.POST("/list/:words/:type/:page",cont.List)
	router.GET("/list/:words/:type/:page",cont.List)

	//资源详情
	router.GET("/hash/:code",cont.Info)

	return router

}
