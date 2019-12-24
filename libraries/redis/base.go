// redis服务类，使用redigo库
package redis

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"math/rand"
	. "my-gin/libraries/config"
	"my-gin/libraries/log"
	"time"
)

// not support redis
var Redigo map[string]map[string][]*redis.Pool

// init intialize redis config
func init() {
	Redigo = make(map[string]map[string][]*redis.Pool, len(DefaultConfig.GetStringMap("redis")))
	for redisKey, redisConf := range UnmarshalConfig.Redis {
		RedisPool := make(map[string][]*redis.Pool)
		for key, value := range redisConf {
			for _, v := range value {
				RedisPool[key] = append(RedisPool[key], newPool(v.Addr, v.Pwd, v.Max_idle, v.Max_active))
			}
		}
		Redigo[redisKey] = RedisPool
	}
}

// NewPool 会返回一个*redis.Pool实例
func newPool(addr string, pwd string, max_idle int, max_active int) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     max_idle,         // 最大空闲连接数，即会有这么多个连接提前等待着，但过了超时时间也会关闭
		MaxActive:   max_active,       // 最大连接数，即最多的tcp连接数，一般建议往大的配置，但不要超过操作系统文件句柄个数（centos下可以ulimit -n查看）
		IdleTimeout: 50 * time.Second, // 空闲连接超时时间，但应该设置比redis服务器超时时间短。否则服务端超时了，客户端保持着连接也没用。
		// MaxConnLifetime: 5 * time.Second,  //关闭早于此持续时间的连接。如果该值为零，则该池不会根据使用期限关闭连接。
		Wait: true, // 如果超过最大连接，是报错，还是等待。
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

// 指定获取redis连接,短连接操作
func Init(redisName string) map[string]*redis.Pool {
	config := UnmarshalConfig.Redis[redisName]
	if config == nil {
		log.InitLog("redis").Errorf("newPool", "redis config not exist：", redisName)
		panic("redis config not exist：" + redisName)
	}
	RedisPool := make(map[string]*redis.Pool)
	randObj := rand.New(rand.NewSource(time.Now().UnixNano()))
	for key, value := range config {
		randNum := randObj.Intn(len(value))
		conf := value[randNum]
		RedisPool[key] = newPool(conf.Addr, conf.Pwd, conf.Max_idle, conf.Max_active)
	}
	return RedisPool
}

type RedisInstanceClass struct {
	redisConnInstance map[string]*redis.Pool
}

// 通过名称获取redis实例
func (rediss *RedisInstanceClass) GetRedigoByName(redisName string) {
	// 单次获取，短连接操作
	// rediss.redisConnInstance = Init(redisName)

	// 获取全局连接池配置，长连接操作
	if Redigo[redisName] == nil {
		log.InitLog("redis").Errorf("newPool", "redis config not exist：", redisName)
		panic("redis config not exist：" + redisName)
	}
	rediss.redisConnInstance = make(map[string]*redis.Pool)
	randObj := rand.New(rand.NewSource(time.Now().UnixNano()))
	for key, val := range Redigo[redisName] {
		randNum := randObj.Intn(len(val))
		rediss.redisConnInstance[key] = val[randNum]
	}
}

func (rediss *RedisInstanceClass) GetMasterPool() redis.Conn {
	return rediss.redisConnInstance["master"].Get()
}

func (rediss *RedisInstanceClass) GetSlavePool() redis.Conn {
	return rediss.redisConnInstance["slave"].Get()
}

// set插入数据
func (rediss *RedisInstanceClass) Set(key string, value interface{}) bool {
	pool := rediss.GetMasterPool()
	defer pool.Close()
	_, err := pool.Do("SET", key, value)
	if err != nil {
		log.InitLog("redis").Errorf("SET", "msg", err)
		fmt.Println("redis set error")
		return false
	}
	return true
}

// setex插入数据与过期时间
func (rediss *RedisInstanceClass) SetEx(key string, value interface{}, seconds int) bool {
	pool := rediss.GetMasterPool()
	defer pool.Close()
	_, err := pool.Do("SET", key, value, "EX", seconds)
	if err != nil {
		log.InitLog("redis").Errorf("SETEX", "msg", err)
		fmt.Println("redis setex error")
		return false
	}
	return true
}

// get获取数据
func (rediss *RedisInstanceClass) Get(key string) string {
	pool := rediss.GetSlavePool()
	defer pool.Close()
	ret, err := redis.String(pool.Do("GET", key))
	if err != nil {
		log.InitLog("redis").Errorf("GET", "msg", err)
		fmt.Println("redis get error")
	}
	return ret
}

// del删除数据
func (rediss *RedisInstanceClass) Del(key interface{}) bool {
	pool := rediss.GetMasterPool()
	defer pool.Close()
	_, err := pool.Do("DEL", key)
	if err != nil {
		log.InitLog("redis").Errorf("DEL", "msg", err)
		fmt.Println("redis del error")
		return false
	}
	return true
}

// hset插入数据
func (rediss *RedisInstanceClass) HSet(key string, field, val interface{}) bool {
	pool := rediss.GetMasterPool()
	defer pool.Close()
	_, err := pool.Do("HSET", key, field, val)
	if err != nil {
		log.InitLog("redis").Errorf("HSET", "msg", err)
		fmt.Println("redis hset error")
		return false
	}
	return true
}

// hmset插入数据
func (rediss *RedisInstanceClass) HMSet(key string, val interface{}) bool {
	pool := rediss.GetMasterPool()
	defer pool.Close()
	_, err := pool.Do("HMSET", redis.Args{}.Add(key).AddFlat(val)...)
	if err != nil {
		log.InitLog("redis").Errorf("HMSET", "msg", err)
		fmt.Println("redis hmset error")
		return false
	}
	return true
}

// hget获取数据
func (rediss *RedisInstanceClass) HGet(key, field string) string {
	pool := rediss.GetSlavePool()
	defer pool.Close()
	ret, err := redis.String(pool.Do("HGET", key, field))
	if err != nil {
		log.InitLog("redis").Errorf("HGET", "msg", err)
		fmt.Println("redis hget error")
	}
	return ret
}

// hmgetall获取数据
func (rediss *RedisInstanceClass) HMGetAll(key string) map[string]string {
	pool := rediss.GetSlavePool()
	defer pool.Close()
	ret, err := redis.StringMap(pool.Do("HGETALL", key))
	if err != nil {
		log.InitLog("redis").Errorf("HMGET", "msg", err)
		fmt.Println("redis hmget error")
	}
	return ret
}

// zadd插入数据
func (rediss *RedisInstanceClass) ZAdd(key string, score int, member string) bool {
	pool := rediss.GetMasterPool()
	defer pool.Close()
	if _, err := pool.Do("ZADD", key, score, member); err != nil {
		log.InitLog("redis").Errorf("ZADD", "msg", err)
		fmt.Println("redis zadd error：" + err.Error())
		return false
	}
	return true
}

// zadd删除数据
func (rediss *RedisInstanceClass) ZRem(key string, member string) bool {
	pool := rediss.GetMasterPool()
	defer pool.Close()
	if _, err := pool.Do("ZREM", key, member); err != nil {
		log.InitLog("redis").Errorf("ZREM", "msg", err)
		fmt.Println("redis zrem error：" + err.Error())
		return false
	}
	return true
}

// zadd获取指点范围正序
func (rediss *RedisInstanceClass) ZRange(key string, start int, stop int) []string {
	pool := rediss.GetSlavePool()
	defer pool.Close()
	ret, err := redis.Strings(pool.Do("ZRANGE", key, start, stop))
	if err != nil {
		log.InitLog("redis").Errorf("ZRANGE", "msg", err)
		fmt.Println("redis zrange error：" + err.Error())
	}
	return ret
}

// zadd获取指点范围倒序
func (rediss *RedisInstanceClass) ZRevrange(key string, start int, stop int) []string {
	pool := rediss.GetSlavePool()
	defer pool.Close()
	ret, err := redis.Strings(pool.Do("ZREVRANGE", key, start, stop))
	if err != nil {
		log.InitLog("redis").Errorf("ZREVRANGE", "msg", err)
		fmt.Println("redis zrevrange error：" + err.Error())
	}
	return ret
}
