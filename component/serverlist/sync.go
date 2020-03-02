package serverlist

import (
	"time"

	"github.com/blastbao/agollo-2/component"
	"github.com/blastbao/agollo-2/component/log"
	"github.com/blastbao/agollo-2/env"
	"github.com/blastbao/agollo-2/env/config"
	"github.com/blastbao/agollo-2/protocol/http"
)

const (
	//refresh ip list
	refreshIPListInterval = 20 * time.Minute //20m
)

func init() {
	go component.StartRefreshConfig(&SyncServerIPListComponent{})
}

//InitSyncServerIPList 初始化同步服务器信息列表
func InitSyncServerIPList() {

}

//SyncServerIPListComponent set timer for update ip list
//interval : 20m
type SyncServerIPListComponent struct {
}

//Start 启动同步服务器列表
func (s *SyncServerIPListComponent) Start() {
	SyncServerIPList(nil)
	log.Debug("syncServerIpList started")

	t2 := time.NewTimer(refreshIPListInterval)
	for {
		select {
		case <-t2.C:
			SyncServerIPList(nil)
			t2.Reset(refreshIPListInterval)
		}
	}
}

//SyncServerIPList sync ip list from server
//then
//1.update agcache
//2.store in disk
func SyncServerIPList(newAppConfig *config.AppConfig) error {
	appConfig := env.GetAppConfig(newAppConfig)
	if appConfig == nil {
		panic("can not find apollo config!please confirm!")
	}

	_, err := http.Request(env.GetServicesConfigURL(appConfig), &env.ConnectConfig{}, &http.CallBack{
		SuccessCallBack: env.SyncServerIPListSuccessCallBack,
	})

	return err
}
