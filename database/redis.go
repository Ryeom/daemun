package database

import (
	"github.com/Ryeom/daemun/log"
	"github.com/gomodule/redigo/redis"
	"github.com/spf13/viper"
	"time"
)

const (
	ForeverTTL = "-1"

	Test            = "test"
	GWConfig        = "gw-config"
	GWwp            = "gw-finance-wp"
	GWfp            = "gw-finance-fp"
	Session         = "session"
	ACRoute         = "achilles-route"
	ACInfo          = "achilles-info"
	GWBigdata       = "gw-big-data-open"
	GWInterPlatform = "gw-interplatform"
)

func NewRedisConnection(platform, target string, index int) *redis.Pool {
	ip := viper.GetString("redis-list-" + platform + "." + target)
	return NewConnection(ip, index)
}

func NewConnection(server string, db int) *redis.Pool {
	return newRedisPool(server+":6379", db, "")
}
func newRedisPool(server string, db int, pw string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     10,
		IdleTimeout: 2 * time.Second,
		MaxActive:   10,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server, redis.DialPassword(pw), redis.DialDatabase(db))
			if err != nil {
				return nil, err
			}
			return c, err
		},
	}
}

func ClosePool(p *redis.Pool) {
	if p != nil {
		p.Close()
	}
}

func AddString(targetRedis *redis.Pool, key string, expireTime string, value string) {
	if targetRedis == nil {
		return
	}
	conn := targetRedis.Get()
	defer conn.Close()
	var err error

	if expireTime == ForeverTTL {
		_, err = conn.Do("SET", key, value)
	} else {
		_, err = conn.Do("SETEX", key, expireTime, value)
	}
	if err != nil {
		log.Logger.Errorf("String Insert ERROR expire time : %s Error : %v", expireTime, err)
	}
}

func AddDefaultValue(targetRedis *redis.Pool, values map[string]interface{}) {
	if targetRedis == nil {
		return
	}
	for key, value := range values {
		if IsExist(targetRedis, key) {
			continue
		}
		switch value.(type) {
		case string:
			v := value.(string)
			AddString(targetRedis, key, "-1", v)
		case map[string]string:
			AddHash(targetRedis, key, value)
		}
	}
}

func IsExist(targetRedis *redis.Pool, key string) bool {
	if targetRedis == nil {
		return false
	}
	conn := targetRedis.Get()
	defer conn.Close()

	ok, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		return false
	}
	return ok
}
func AddList(targetRedis *redis.Pool, key, value string) {
	if targetRedis == nil {
		return
	}
	conn := targetRedis.Get()
	defer conn.Close()
	_, err := conn.Do("RPUSH", key, value)
	if err != nil {
		log.Logger.Errorf("List Insert ERROR : %v", err)
	}
}

func AddExpire(targetRedis *redis.Pool, key string, ttl int) {
	if targetRedis == nil {
		return
	}
	conn := targetRedis.Get()
	defer conn.Close()
	_, err := conn.Do("EXPIRE", key, ttl)
	if err != nil {
		log.Logger.Errorf("EXPIRE Insert ERROR : %v", err)
	}
}

func AddHash(targetRedis *redis.Pool, key string, value interface{}) {
	if targetRedis == nil {
		return
	}
	conn := targetRedis.Get()
	defer conn.Close()
	_, err := conn.Do("HSET", redis.Args{}.Add(key).AddFlat(value)...)
	if err != nil {
		log.Logger.Errorf("HASH Insert ERROR : %v", err)
	}
}

func RemoveList(targetRedis *redis.Pool, key, value string) bool {
	if targetRedis == nil {
		return false
	}
	conn := targetRedis.Get()
	_, err := conn.Do("LREM", key, 1, value)
	if err != nil {
		log.Logger.Error(err)
		return false
	}
	return true
}

func GetString(targetRedis *redis.Pool, keyword string) string {
	if targetRedis == nil {
		return ""
	}
	conn := targetRedis.Get()
	defer conn.Close()
	str, err := redis.String(conn.Do("GET", keyword))
	if err != nil {
		log.Logger.Error(err)
		return ""
	}
	return str
}
func GetExpireTime(targetRedis *redis.Pool, keyword string) int {
	if targetRedis == nil {
		return -9
	}
	conn := targetRedis.Get()
	defer conn.Close()
	ttl, err := redis.Int(conn.Do("TTL", keyword))
	if err != nil {
		log.Logger.Error(err)
		return -9
	}
	return ttl
}
func GetList(targetRedis *redis.Pool, key string) []string {
	if targetRedis == nil {
		return nil
	}
	conn := targetRedis.Get()
	defer conn.Close()
	value, err := redis.Strings(conn.Do("LRANGE", key, 0, -1))
	if err != nil {
		log.Logger.Error(err)
	}
	return value
}

func GetType(targetRedis *redis.Pool, key string) string {
	if targetRedis == nil {
		return ""
	}
	conn := targetRedis.Get()
	defer conn.Close()
	value, err := redis.String(conn.Do("TYPE", key))
	if err != nil {
		log.Logger.Error(err)
	}
	return value
}

func ScanKeyList(targetRedis *redis.Pool, key string) []string {
	if targetRedis == nil {
		return nil
	}
	conn := targetRedis.Get()
	defer conn.Close()
	nextCursor := 0
	var result []string
	for {
		arr, err := redis.Values(conn.Do("SCAN", nextCursor, "MATCH", key, "COUNT", 1000))
		if err != nil {
			return result
		} else {
			nextCursor, _ = redis.Int(arr[0], nil)
			tempUserList, _ := redis.Strings(arr[1], nil)
			result = append(result, tempUserList...)
		}
		if nextCursor == 0 {
			break
		}
	}
	return result
}

func GetHash(targetRedis *redis.Pool, key string) map[string]string {
	if targetRedis == nil {
		return nil
	}
	conn := targetRedis.Get()
	defer conn.Close()
	obj, err := redis.StringMap(conn.Do("HGETALL", key))
	if err != nil {
		log.Logger.Error(err)
	}
	return obj
}

func GetSet(targetRedis *redis.Pool, key string) []string {
	if targetRedis == nil {
		return nil
	}
	conn := targetRedis.Get()
	defer conn.Close()
	obj, err := redis.Strings(conn.Do("SMEMBERS", key))
	if err != nil {
		log.Logger.Error(err)
	}
	return obj
}
