//redis服务类，使用redigo库
package redis

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	. "my-gin/libraries/config"
	"my-gin/libraries/log"
	"time"
)

// not support redis
var Redigo map[string]map[string]*redis.Pool

// init intialize redis config
func init() {

	Redigo = make(map[string]map[string]*redis.Pool, len(DefaultConfig.GetStringMap("redis")))
	for key, c := range DefaultConfig.GetStringMap("redis") {
		conArr := c.([]interface{})
		Redigo[key] = make(map[string]*redis.Pool)
		for _, config := range conArr {
			conOne := config.(map[interface{}]interface{})
			addr := conOne["addr"].(string)
			pwd := conOne["pwd"].(string)
			max_idle := conOne["max_idle"].(int)
			max_active := conOne["max_active"].(int)

			//多层map赋值需要每层创建默认空map
			Redigo[key][conOne["instance"].(string)] = newPool(addr, pwd, max_idle, max_active)
		}
	}
}

// NewPool 会返回一个*redis.Pool实例
func newPool(addr string, pwd string, max_idle int, max_active int) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     max_idle,         //最大空闲连接数，即会有这么多个连接提前等待着，但过了超时时间也会关闭
		MaxActive:   max_active,       //最大连接数，即最多的tcp连接数，一般建议往大的配置，但不要超过操作系统文件句柄个数（centos下可以ulimit -n查看）
		IdleTimeout: 20 * time.Second, //空闲连接超时时间，但应该设置比redis服务器超时时间短。否则服务端超时了，客户端保持着连接也没用。
		Wait:        true,             //如果超过最大连接，是报错，还是等待。
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial(
				"tcp",
				addr,
				redis.DialConnectTimeout(time.Second*1),
				redis.DialReadTimeout(time.Second*1),
				redis.DialWriteTimeout(time.Second*1),
			)
			if err != nil {
				log.InitLog("redis").Errorf("newPool", "msg", err.Error())
				fmt.Printf("redis连接失败:%s\n", err.Error())
				return nil, err
			}
			if _, err := c.Do("AUTH", pwd); err != nil {
				c.Close()
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
}

type RedisInstanceClass struct {
	redisConnInstance *redis.Pool
}

// 通过名称获取redis实例
func (redis *RedisInstanceClass) GetRedigoByName(redisName string, instance string) {
	redis.redisConnInstance = Redigo[redisName][instance]
}

// set插入数据
func (rediss *RedisInstanceClass) Set(key string, value interface{}) bool {
	_, err := rediss.redisConnInstance.Get().Do("SET", key, value)
	if err != nil {
		log.InitLog("redis").Errorf("SET", "msg", err)
		fmt.Println("redis set error")
		return false
	}
	return true
}

// setex插入数据与过期时间
func (rediss *RedisInstanceClass) SetEx(key string, value interface{}, seconds int) bool {
	_, err := rediss.redisConnInstance.Get().Do("SET", key, value, "EX", seconds)
	if err != nil {
		log.InitLog("redis").Errorf("SETEX", "msg", err)
		fmt.Println("redis setex error")
		return false
	}
	return true
}

// get获取数据
func (rediss *RedisInstanceClass) Get(key string) string {
	ret, err := redis.String(rediss.redisConnInstance.Get().Do("GET", key))
	if err != nil {
		log.InitLog("redis").Errorf("GET", "msg", err)
		fmt.Println("redis get error")
	}
	return ret
}

// del删除数据
func (rediss *RedisInstanceClass) Del(key interface{}) bool {
	_, err := rediss.redisConnInstance.Get().Do("DEL", key)
	if err != nil {
		log.InitLog("redis").Errorf("DEL", "msg", err)
		fmt.Println("redis del error")
		return false
	}
	return true
}

// hset插入数据
func (rediss *RedisInstanceClass) HSet(key string, field, val interface{}) bool {
	_, err := rediss.redisConnInstance.Get().Do("HSET", key, field, val)
	if err != nil {
		log.InitLog("redis").Errorf("HSET", "msg", err)
		fmt.Println("redis hset error")
		return false
	}
	return true
}

// hmset插入数据
func (rediss *RedisInstanceClass) HMSet(key string, val interface{}) bool {
	_, err := rediss.redisConnInstance.Get().Do("HMSET", redis.Args{}.Add(key).AddFlat(val)...)
	if err != nil {
		log.InitLog("redis").Errorf("HMSET", "msg", err)
		fmt.Println("redis hmset error")
		return false
	}
	return true
}

// hget获取数据
func (rediss *RedisInstanceClass) HGet(key, field string) interface{} {
	ret, err := redis.String(rediss.redisConnInstance.Get().Do("HGET", key, field))
	if err != nil {
		log.InitLog("redis").Errorf("HGET", "msg", err)
		fmt.Println("redis hget error")
	}
	return ret
}

// hmgetall获取数据
func (rediss *RedisInstanceClass) HMGetAll(key string) interface{} {
	ret, err := redis.StringMap(rediss.redisConnInstance.Get().Do("HGETALL", key))
	if err != nil {
		log.InitLog("redis").Errorf("HMGET", "msg", err)
		fmt.Println("redis hmget error")
	}
	return ret
}
