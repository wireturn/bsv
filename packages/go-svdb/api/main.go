package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"

	"github.com/golang/glog"
	"mempool.com/foundation/bitdb/common"
	go_util "mempool.com/foundation/go-util"
)

type ConfigData struct {
	EnableStatistics        bool
	StatisticsFlushInterval int
	Listen                  string
	Fullnode                common.FullnodeConfig
	Mongo                   map[string]common.MongoConfig
	IsTestEnvironment       bool
	SyncCenterIp            string
}

func getConfig() (result ConfigData, ok bool, errMsg string) {
	// 配置数据
	var configData ConfigData
	// 解析命令行参数
	configFilePath := flag.String("config", "./config.json", "Path of config file")
	flag.Parse()
	// 读取配置文件
	configJSON, err := ioutil.ReadFile(*configFilePath)
	if err != nil {
		errMsg = "read config failed: " + err.Error()
		return configData, false, errMsg
	}
	err = json.Unmarshal(configJSON, &configData)
	if err != nil {
		errMsg = "parse config failed: " + err.Error()
		return configData, false, errMsg
	}
	return configData, true, "ok"
}

type GlobalHandle struct {
	conf      *ConfigData
	mongo     *common.MongoManager
	fullnode  *common.FullnodeManager
	productor *common.KafkaManager
}

var gHandle GlobalHandle
var glbBitdbIp string

func InitHandle(fc common.FullnodeConfig, mc map[string]common.MongoConfig) (*common.FullnodeManager, *common.MongoManager) {
	pFullnodeManager := &common.FullnodeManager{}
	pMongoManager := &common.MongoManager{}
	// pKafkaProductorManager := &common.KafkaManager{}
	pFullnodeManager.Init()
	pMongoManager.Init()
	// pKafkaProductorManager.Init()
	for coinType, fullnode := range fc {
		_, err := pFullnodeManager.AddFullnodeInfo(coinType, fullnode)
		if err != nil {
			panic(err)
		}

		pMongoManager.AddServiceHub(
			coinType,
			mc[coinType].AddrService,
			mc[coinType].BlockHeaderService,
			mc[coinType].AddrIndexService,
			mc[coinType].HeightAddrService,
			mc[coinType].HeightHashService,
			mc[coinType].MonitoredService,
			mc[coinType].SyncMonitoredAddrService)
		// pKafkaProductorManager.AddKafkaProductor(kafkaUrl, coinType)
	}
	return pFullnodeManager, pMongoManager
}

func (p *GlobalHandle) init() {
	p.fullnode, p.mongo = InitHandle(p.conf.Fullnode, p.conf.Mongo)
	go_util.Init(p.conf.IsTestEnvironment)
}

func main() {
	config, ok, note := getConfig()
	if !ok {
		glog.Error("Read Config file Error," + note)
		return
	}
	bitdbIp, err := common.GetIntranetIp()
	if err != nil {
		panic(err)
	}
	glbBitdbIp = bitdbIp
	gHandle.conf = &config
	gHandle.init()
	common.GetStatisticInstance().Init(config.EnableStatistics, config.StatisticsFlushInterval)
	startHttp(&config)
	//select {} //阻塞进程
}
