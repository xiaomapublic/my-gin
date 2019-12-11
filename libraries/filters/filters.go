package filters

import (
	"github.com/gin-contrib/cache"
	"github.com/gin-contrib/cache/persistence"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"my-gin/libraries/config"
	"time"
)

func RegisterSession() gin.HandlerFunc {
	sessionConfig := config.UnmarshalConfig.Redis["session"]["master"][0]
	store, _ := sessions.NewRedisStore(
		10,
		"tcp",
		sessionConfig.Addr,
		sessionConfig.Pwd,
		[]byte("secret"))
	return sessions.Sessions("session", store)
}

func RegisterCache() gin.HandlerFunc {
	sessionConfig := config.UnmarshalConfig.Redis["session"]["master"][0]
	var cacheStore persistence.CacheStore
	cacheStore = persistence.NewRedisCache(sessionConfig.Addr, sessionConfig.Pwd, time.Minute)
	return cache.Cache(&cacheStore)
}
