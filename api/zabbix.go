package api

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"zbx-api/config"
	"zbx-api/models"

	"github.com/canghai908/zabbix-go"
)


type HostGetResult struct {
	Host       string `json:"host"`
	HostID     string `json:"hostid"`
	Interfaces []map[string]string
}


// 获取主机ip/hostid信息
func GetHosts(api *zabbix.API)[]HostGetResult{
	params := map[string]interface{}{
		"output":           []string{"hostid", "host"}, //需求数据，监控项的name 和最新的值
		"selectInterfaces": []string{"ip","interfaceid"}}             //群组名称

	response, err := api.CallWithError("host.get", params)
	var result []HostGetResult
	ret,err:=json.Marshal(response.Result)
	if err!=nil {
		//zap.l.Fatal("json 解析错误",zap.Error(err))
		log.Fatal(err)
	}
	err= json.Unmarshal(ret,&result)
	return result
}


// 获取CPU MEM
func ItemStatGet(api *zabbix.API,key string) zabbix.Items{
	params := map[string]interface{}{
		"output":      []string{"hostid", "name", "lastvalue"}, //需求数据，监控项的name 和最新的值
		"search":      map[string]string{"key_": key}}          //监控项

	ret, err := api.ItemsGet(params)
	if err != nil {
		log.Fatal(err)
		//zap.L().Fatal("获取监控数据失败", zap.Error(err))
	}
	return ret
}


// 通过API 获取数据
func ItemStat(api *zabbix.API,statConfigMap map[string]string)(string,error){
	statMap :=make(map[string]zabbix.Items)
	for k,v:=range statConfigMap{
		statMap[k] = ItemStatGet(api,v)
	}

	resultMap:=make(map[string]map[string]interface{})
	for _,v:=range GetHosts(api){
		resultMap[v.HostID] = make(map[string]interface{})
		resultMap[v.HostID] ["hostId"] = v.HostID
		resultMap[v.HostID] ["ip"] = v.Interfaces[0]["ip"]
	}
	for k,_:=range statConfigMap{
		for _,v2:=range statMap[k] {
			_,ok:=resultMap[v2.HostId]
			if ok{
				resultMap[v2.HostId][k] = v2.Lastvalue
			}
		}
	}

	var resultList []interface{}
	for _,v:=range resultMap{
		resultList = append(resultList,v)
	}
	data,err:=json.Marshal(resultList)
	if err != nil{
		return "", err
	}
	return string(data),nil


}

// 获取key
func GetZabbixKey(c *gin.Context)  {
	serverId:=c.Param("id")
	data:=models.RdsClient.HGet("zabbix",serverId)
	var startList []map[string]string
	err:=json.Unmarshal([]byte(data.Val()),&startList)
	if err!=nil{
		c.JSON(http.StatusBadGateway,startList)	
	}else{
		c.JSON(http.StatusOK,startList)
	}
}


// 设置key

func SetZabbixToRedis(){
	configMap :=config.InitServersConfig()
	apiMap :=make(map[string]*zabbix.API) // 存储API
	for _,v:=range configMap.Servers{
		url := fmt.Sprintf("http://%s/api_jsonrpc.php", v.Host)
		api:=zabbix.NewAPI(url)
		api.Login(v.User,v.Password)
		apiMap[v.ID] = api
	}

	for _,v:=range configMap.Servers{
		itemStat,err:=ItemStat(apiMap[v.ID],configMap.Items)
		if err!=nil{
			log.Fatal(err)
		}
		models.RdsClient.HSet("zabbix",v.ID,itemStat)
	}

}




