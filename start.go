package agollo

import (
	"github.com/blastbao/agollo-2/agcache"
	"github.com/blastbao/agollo-2/component"
	"github.com/blastbao/agollo-2/component/log"
	"github.com/blastbao/agollo-2/component/notify"
	"github.com/blastbao/agollo-2/component/serverlist"
	"github.com/blastbao/agollo-2/env"
	"github.com/blastbao/agollo-2/env/config"
	"github.com/blastbao/agollo-2/loadbalance/roundrobin"
	"github.com/blastbao/agollo-2/storage"
)

func init() {
	roundrobin.InitLoadBalance()
	serverlist.InitSyncServerIPList()
}

//InitCustomConfig init config by custom
func InitCustomConfig(loadAppConfig func() (*config.AppConfig, error)) {
	env.InitConfig(loadAppConfig)
}

//start apollo
func Start() error {
	return startAgollo()
}

//SetLogger 设置自定义logger组件
func SetLogger(loggerInterface log.LoggerInterface) {
	if loggerInterface != nil {
		log.InitLogger(loggerInterface)
	}
}

//SetCache 设置自定义cache组件
func SetCache(cacheFactory agcache.CacheFactory) {
	if cacheFactory != nil {
		agcache.UseCacheFactory(cacheFactory)
		storage.InitConfigCache()
	}
}

func startAgollo() error {
	//first sync
	if err := notify.SyncConfigs(); err != nil {
		return err
	}
	log.Debug("init notifySyncConfigServices finished")

	//start long poll sync config
	go component.StartRefreshConfig(&notify.ConfigComponent{})

	log.Info("agollo start finished ! ")

	return nil
}
