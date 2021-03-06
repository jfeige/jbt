package models

import (
	"errors"
	"github.com/garyburd/redigo/redis"
	"github.com/jfeige/lconfig"
)

var (
	BtPageSize int
	AppPort string
)

var (
	lcf  lconfig.LConfigInterface
	err  error
	conn *connect
)

type connect struct {
	pool *redis.Pool
}

//读取配置文件，初始化数据库和redis连接池
func InitBaseConfig(file string) error {
	lcf, err = lconfig.NewConfig(file)
	if err != nil {
		return err
	}
	//redis配置
	err = initRedisConfig()
	if err != nil {
		return err
	}
	//其他配置参数
	err = initDefaultConfig()
	if err != nil {
		return err
	}

	conn = &connect{}
	conn.pool = initRedisPool()

	return nil
}

/**
获取redis连接
*/
func (this *connect) GetRedisConn() redis.Conn {
	return this.pool.Get()
}

/**
读取其他配置
*/
func initDefaultConfig() error {
	AppPort = lcf.String("app_port")
	if AppPort == ""{
		return errors.New("config parameters:app_port is error!")
	}

	BtPageSize, _ = lcf.Int("pagesize")
	if BtPageSize <= 0 {
		return errors.New("Can't not find default parameters:pagesize")
	}
	return nil
}
