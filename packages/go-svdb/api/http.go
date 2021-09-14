package main

import (
	"net/http"

	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"mempool.com/foundation/bitdb/common"
)

func startHttp(config *ConfigData) {
	r := mux.NewRouter()
	//获取节点总体数据：ok
	r.HandleFunc("/api/blockchaininfo/{type}", common.Logging(GetBlockChainInfo))
	//根据高度，获取某个区块信息,ok
	r.HandleFunc("/api/blockByHeight/{type}/{height}", common.Logging(GetBlockByHeight))
	//根据区块hash值，获取区块信息,ok
	r.HandleFunc("/api/blockByHash/{type}/{hash}", common.Logging(GetBlockByHash))
	//获取最近的区块：
	r.HandleFunc("/api/bestblock/{type}/{start}/{limit}", common.Logging(GetBestBlock))
	//根据区块hash获取预览交易数据
	r.HandleFunc("/api/block/browser_txs/{type}/{hash}/{start}/{end}", common.Logging(GetBlockBrowserTxsByHash))
	//根据交易hash获blockByHash取交易数据：
	r.HandleFunc("/api/tx/{type}/{txid}", common.Logging(GetTxInfo))
	//根据交易hash获取浏览器需要交易数据：
	r.HandleFunc("/api/browser_tx/{type}/{txid}", common.Logging(GetBrowserTxInfo))
	//根据地址获取浏览器需要tx交易预览数据：
	r.HandleFunc("/api/address/browser_txs/{type}/{address}/{start}/{end}", common.Logging(GetAddressBrowserTxs))
	//根据地址获取地址余额：
	r.HandleFunc("/api/unspentinfo/{type}/{address}", common.Logging(ListUnspentInfo))
	//根据地址获取可以花费的交易
	r.HandleFunc("/api/utxo/{type}/{address}", common.Logging(ListUnspentUtxo))
	//获取区块链概况
	r.HandleFunc("/api/blockchain/state/{type}", common.Logging(GetChainState))
	//获取全网难度
	r.HandleFunc("/api/blockchain/difficulty/{type}", common.Logging(GetChainDifficulty))
	//获取用户的余额
	r.HandleFunc("/api/monitored_address/balance/type/{type}/appid/{appid}/ownerid/{ownerid}", common.Logging(GetBalance))
	//获取历史交易流水
	r.HandleFunc("/api/monitored_address/inventory/type/{type}/appid/{appid}/ownerid/{ownerid}/start/{start}/end/{end}", common.Logging(GetInventory))
	//获取用户的utxo
	r.HandleFunc("/api/monitored_address/utxo/type/{type}/appid/{appid}/ownerid/{ownerid}/value/{value}", common.Logging(GetMonitoredAddrUtxo))
	//添加监听地址
	r.HandleFunc("/api/monitored_address/add/type/{type}/appid/{appid}/ownerid/{ownerid}/addr/{addr}", common.Logging(AddMonitoredAddr))
	r.HandleFunc("/api/monitored_address/sync/type/{type}/appid/{appid}/ownerid/{ownerid}/addr/{addr}", common.Logging(SyncMonitoredAddr))
	//获取监听地址owner
	r.HandleFunc("/api/monitored_address/get/type/{type}/ownerid/{ownerid}", common.Logging(GetMonitoredAddrOwner))
	//测试用的地址-newtx
	r.HandleFunc("/api/test/onnewtx/txhashstr/{txhashstr}/index/{index}/isinput/{isinput}/addr/{addr}/value/{value}", common.Logging(OnNewTx))
	//测试用的地址-newblock
	r.HandleFunc("/api/test/onnewblock/height/{height}/hash/{hash}", common.Logging(OnNewBlock))
	//用来根据ownerid检查余额
	r.HandleFunc("/api/test/CheckOwneridBalance/type/{type}/appid/{appid}/ownerid/{ownerid}", common.Logging(CheckOwneridBalance))
	r.HandleFunc("/api/GetOwneridCount/type/{type}/appid/{appid}", common.Logging(GetOwneridCount))
	glog.Fatal(http.ListenAndServe(config.Listen, r))
}
