package controllers

import (
	"gopkg.in/gin-gonic/gin.v1"
	"github.com/jfeige/ltools"
	"net/http"
	"jbt/models"
	"strconv"
	"fmt"
)

func Index(context *gin.Context){

	hotWords := models.GetHotSearch()

	context.HTML(http.StatusOK,"index.html",gin.H{
		"hotWords":hotWords,
	})

}


//form提交搜索
func List(context *gin.Context){
	pagesize := 10
	words := context.Param("words")		//搜索关键字
	tmp_page := context.Param("page")		//页码

	page,_ := strconv.Atoi(tmp_page)
	if page <= 0{
		page = 1
	}
	offset := (page-1) * pagesize
	tp := context.Param("type")			//类型 time_d,size_d,rala_d

	if !ltools.InArray(tp,[]string{"time_d","size_d"}){
		tp = "time_d"
	}

	if words == ""{
		context.HTML(http.StatusOK,"index.html",nil)
	}else{
		recordList := models.Search(words,"time_d",page,offset,pagesize)
		var prev,next int
		if page > 1{
			prev = page -1
		}
		if len(recordList) == pagesize{
			next = page + 1
		}

		context.HTML(http.StatusOK,"list.html",gin.H{
			"recordList":recordList,
			"words":words,
			"page":page,
			"prev":prev,
			"next":next,
		})
	}

}

/**
	资源详情
 */
func Info(context *gin.Context){
	code := context.Param("code")
	fmt.Println("1:",code)
	if code == ""{
		//跳转到错误页面
		return
	}
	record := new(models.Record)
	err := record.Load(code)
	if err != nil{
		fmt.Println(err)
		//没有找到资源，错误提示
		return
	}
	context.HTML(http.StatusOK,"info.html",gin.H{
		"record":record,
	})
}
