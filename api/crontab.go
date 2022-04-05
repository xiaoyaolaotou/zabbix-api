package api

import (
	"fmt"
	"go.uber.org/zap"
	"time"
	"zbx-api/config"
)

// 定时任务
func Cron() {
	serversConfig := config.InitServersConfig()
	t :=time.Tick(time.Second * time.Duration(serversConfig.Interval))
	for{
		select {
		case  <-t:
			fmt.Println("tick---定时任务触发----")
			zap.L().Info("tick---定时任务触发----")
			go SetZabbixToRedis()
		}
	}
}
