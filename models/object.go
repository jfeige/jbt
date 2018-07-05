package models

import (
	"github.com/garyburd/redigo/redis"
	"errors"
	log "github.com/alecthomas/log4go"
	"strings"
	"html/template"
)

type Record struct {
	Title string
	Size string
	Count string
	Time string
	Hot string
	Murl string
	Turl string
	Code string				//该资源的唯一值，base64(Murl)
	Files string			//文件列表
}

func (this *Record) Load(code string)error{
	key := "record:" + code
	rcon := conn.pool.Get()
	defer rcon.Close()

	exists, err := redis.Int(rcon.Do("EXISTS", key))
	if err == nil && exists == 1 {
		values, err := redis.Values(rcon.Do("HGETALL", key))
		if err == nil {
			if len(values) == 0 {
				log.Error("record not exists!code:%s",code)
				return errors.New("record not exists!")
			} else {
				err = redis.ScanStruct(values, this)
				if err == nil {
					return nil
				}
				return err
			}
		}else{
			return err
		}
	}else{
		log.Error("record not exists!code:%s",code)
		return errors.New("record not exists!")
	}
	return errors.New("known record!")
}

/**
	文件列表
 */
func (this *Record) GetFiles()[]string{

	files := strings.Split(this.Files,"&&")
	return files
}

/**
	获取磁力链接
 */
func (this *Record) GetMurl()template.URL{

	return template.URL(this.Murl)

}

/**
	获取迅雷链接
 */
func (this *Record) GetTurl()template.URL{

	return template.URL(this.Turl)

}