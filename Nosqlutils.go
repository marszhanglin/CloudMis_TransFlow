// Nosqlutils.go
package main

import (
	"time"

	"github.com/garyburd/redigo/redis"
)

var (
	// 定义常量
	RedisClient *redis.Pool
	REDIS_HOST  string
	REDIS_DB    int
	REDIS_MODLE string
)

func initRedis() {
	// 建立连接池
	REDIS_HOST = "ip:port"
	REDIS_DB = 0
	REDIS_MODLE = "CLOUDMIS_"
	RedisClient = &redis.Pool{
		// 从配置文件获取maxidle以及maxactive，取不到则用后面的默认值
		MaxIdle:     200,  //maxidle
		MaxActive:   1024, //maxactive
		IdleTimeout: 10 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", REDIS_HOST)
			if err != nil {
				return nil, err
			}
			// 认证
			if _, err := c.Do("AUTH", "xxxx"); err != nil {
				c.Close()
				return nil, err
			}
			//			// 选择db
			//			if _, err := c.Do("SELECT", REDIS_DB); err != nil {
			//				c.Close()
			//				return nil, err
			//			}
			return c, nil
		},
	}
}

// 设置字符串   expire（second）
func setNoSqlStrExpire(key string, expire int64, value string) error {
	rc := RedisClient.Get()
	_, err := redis.Strings(rc.Do("SETEX", key, expire, value))
	nosqlErrHandler(err)
	defer rc.Close()
	return err
}

func getNoSqlStr(key string) (string, error) {
	rc := RedisClient.Get()
	something, err := redis.String(rc.Do("get", key))
	nosqlErrHandler(err)
	glogInfo("redis get(" + key + "):" + something)
	defer rc.Close()
	return something, err
}

func nosqlErrHandler(err error) {
	if nil != err {
		glogError("redisError:" + err.Error())
		glogFlush()
	}
}
