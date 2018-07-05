package models

import (
	"encoding/base64"
	"github.com/garyburd/redigo/redis"
	log "github.com/alecthomas/log4go"
	"github.com/PuerkitoBio/goquery"
	"github.com/jfeige/ltools"
	"strings"
	"sync"
	"strconv"
	"time"
	"fmt"
)
var(
	 wg sync.WaitGroup
)
/**
	根据关键字搜索
 */
func Search(words,tp string,page,offset,pagesize int)[]*Record{
	rcon := conn.pool.Get()
	defer rcon.Close()
	w := base64.StdEncoding.EncodeToString([]byte(words))
	key := "list:" + w + ":" + tp

	ret := make([]*Record,0)

	result,err := redis.Strings(rcon.Do("ZREVRANGE",key,offset,(offset + pagesize-1)))
	if len(result) == 0 || err != nil{
		if err != nil{
			log.Error("Search has error:%v",err)
			return ret
		}
		//未命中缓存，从网络查询
		SearchFromNetWork(key,tp,words,page)

		//单线程查询其他页码 1~50页
		for i := page+1; i <= page+50;i++{
			go SearchFromNetWork(key,tp,words,i)
		}

		//读取缓存
		result,err = redis.Strings(rcon.Do("ZREVRANGE",key,offset,(offset + pagesize-1)))
	}
	for _,record := range result{
		datas := strings.Split(record,"&*&")
		file_title := datas[0]
		file_size := datas[1]
		file_count := datas[2]
		create_time := datas[3]
		file_hot := datas[4]
		magnet_url := datas[5]
		thunder_url := datas[6]
		files := datas[7]
		code := ltools.ToMd5(magnet_url)

		obj := new(Record)
		obj.Title = file_title
		obj.Size = file_size
		obj.Count = file_count
		obj.Time = create_time
		obj.Hot = file_hot
		obj.Murl = magnet_url
		obj.Turl = thunder_url
		obj.Code = code
		obj.Files = files

		//每次用户请求放入redis，保证资源的有效性
		go saveToRedis(obj)

		ret = append(ret,obj)
	}

	//写入搜索统计
	go countSearch(words)

	return ret

}

/**
	从网络查询数据
 */
func SearchFromNetWork(key,tp,words string,page int){
	var url = "http://www.btsou8.net/list/" + words + "/" +strconv.Itoa(page)+ "/" + tp
	doc,err := goquery.NewDocument(url)
	if err != nil{
		log.Error("SearchFromNetWork has error:%v",err)
		return
	}
	doc.Find(".T1").Each(func(i int, selection *goquery.Selection) {
		url,_ := selection.Find("[name=file_title]").Attr("href")
		if url != ""{
			wg.Add(1)
			go SearchContent(key,tp,url,page)

		}
	})
	wg.Wait()

}

/**
	查询消息内容
 */
func SearchContent(key,tp,url string,page int){
	defer wg.Done()
	var data string
	url = "http://www.btsou8.net" + url
	doc,err := goquery.NewDocument(url)
	if err != nil{
		log.Error("SearchContent has error:%v",err)
		return
	}
	title := doc.Find(".T2").Text()

	selection := doc.Find("dl p")
	fmt.Println("----------",selection.Eq(0).Text())
	if selection.Eq(0).Text() == ""{
		return
	}
	if selection.Eq(1).Text() == ""{
		return
	}
	if selection.Eq(2).Text() == ""{
		return
	}
	if selection.Eq(3).Text() == ""{
		return
	}
	if selection.Eq(4).Text() == ""{
		return
	}
	if selection.Eq(5).Text() == ""{
		return
	}
	if selection.Eq(6).Text() == ""{
		return
	}
	file_size := strings.Split(selection.Eq(0).Text(),"：")[1]
	file_count := strings.Split(selection.Eq(1).Text(),"：")[1]
	create_time := strings.Split(selection.Eq(2).Text(),"：")[1]
	file_hot := strings.Split(selection.Eq(4).Text(),"：")[1]
	magnet_url,_ := selection.Eq(5).Find("a").Attr("href")
	thunder_url,_:= selection.Eq(6).Find("a").Attr("href")

	//创建时间   文件大小

	data = title + "&*&" + file_size + "&*&" + file_count + "&*&" + create_time + "&*&" + file_hot + "&*&" + magnet_url + "&*&" + thunder_url

	filters := []string{"txt","TXT","url","</a>","mht","html"}

	args := make([]interface{},0)
	args = append(args,key)
	var files_ret string
	fs := make([]string,0)
	//文件列表
	doc.Find(".flist li").Each(func(i int, selection *goquery.Selection) {
		isfilter := false
		c,_ := selection.Html()
		c = strings.Replace(c,"<span>",":",1)
		c = strings.Replace(c,"</span>","",1)
		files := strings.Split(c,":")
		for _,v := range filters{
			if strings.HasSuffix(files[0],v){
				isfilter = true
				break
			}
		}
		if !isfilter{
			fs = append(fs,c)
		}
	})
	files_ret = strings.Join(fs,"&&")
	files_ret = strings.Replace(files_ret,":"," ",-1)
	data = data + "&*&" + files_ret
	var score float64
	if tp == "time_d"{
		//创建时间
		ctime := strings.Replace(create_time,"-","",-1)
		score,_ = strconv.ParseFloat(ctime,64)
	}else if tp == "size_d"{
		values := strings.Split(file_size," ")
		//文件大小
		if len(values) == 2{
			size,_ := strconv.ParseFloat(values[0],64)
			if values[1] == "MB"{
				score = size * 1000000
			}else if values[1] == "GB"{
				score = size * 1000000000
			}else if values[1] == "KB" {
				score = size*1000
			}else if values[1] == "字节"{
				score = size
			}else{
				return
			}
		}else{
			return
		}
	}else{
		score,_ = strconv.ParseFloat(file_hot,64)
	}
	rcon := conn.pool.Get()
	defer rcon.Close()
	_,err = rcon.Do("ZADD",key,score,data)
}

/**
	把单个资源放入redis
 */
func saveToRedis(obj *Record){
	rcon := conn.pool.Get()
	defer rcon.Close()

	key := "record:" + obj.Code

	rcon.Do("HMSET",redis.Args{}.Add(key).AddFlat(obj)...)

	rcon.Do("expire",key,86400)
}

/**
	搜索统计
 */
func countSearch(words string){
	rcon := conn.pool.Get()
	defer rcon.Close()

	//热门搜索关键词
	key := "hotWords"
	rcon.Do("zincrby",key,1,words)


	//最近搜索关键词
	key = "justWords"
	rcon.Do("ZADD",key,time.Now().UnixNano()/1e6,words)
	cnt,_ := redis.Int(rcon.Do("ZCARD",key))
	if cnt > 20{
		rcon.Do("zremrangebyrank",key,0,0)
	}
}

/**
	获取热搜前8的关键词
 */
func GetHotSearch()[]string{
	rcon := conn.pool.Get()
	defer rcon.Close()

	key := "hotWords"
	words,_ := redis.Strings(rcon.Do("zrevrange",key,0,7))

	return words
}

/**
	获取最近搜索前20的关键词
 */
func GetJustSearch()[]string{
	rcon := conn.pool.Get()
	defer rcon.Close()

	key := "justWords"
	words,_ := redis.Strings(rcon.Do("zrevrange",key,0,19))

	return words
}