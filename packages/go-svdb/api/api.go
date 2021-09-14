package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcutil"
	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"github.com/kataras/iris/core/errors"
	"mempool.com/foundation/bitdb/common"
	go_util "mempool.com/foundation/go-util"
)

const (
	BIT_BTC = "BTC"
	BIT_BCH = "BCH"
	BIT_BSV = "BSV"
)

const (
	ADDRTX_txid    = "txid"
	ADDRTX_isinput = "isinput"
	ADDRTX_index   = "index"
)

const (
	ListUnspentInfo_ErrorCode  = -100
	CancelOrder_ErrorCode      = -110
	GetTxInfo_ErrorCode        = -120
	GetBrowserTxInfo_ErrorCode = -130
	GetChainInfo_ErrorCode     = -150

	GetBlockByHeight_ErrorCode     = -160
	GetBlockByHash_ErrorCode       = -170
	GetBestBlock_ErrorCode         = -180
	GetBlockBrowserTxsv_ErrorCode  = -190
	GetAddressBrowserTxs_ErrorCode = -200

	GetChainState_ErrorCode           = -210
	GetMonitoredAddrUtxo_ErrorCode    = -220
	GetChangedMonitoredAddr_ErrorCode = -230
	GetBalance_ErrorCode              = -240
	GetInventory_ErrorCode            = -250
	GetChainDifficulty_ErrorCode      = -260
	GetOwnerBalance_ErrorCode         = -270
)

func (p *GlobalHandle) GetRpcCli(bittype string) *rpcclient.Client {
	cli := p.fullnode.GetFullnode(bittype)
	if cli == nil {
		glog.Error("Get fullnode Handle fail, type:" + bittype)
	}
	return cli
}

func (p *GlobalHandle) MongoHandle(bittype string) (*common.MongoServiceHub, error) {
	handle, ok := p.mongo.Servicehub[bittype]
	if !ok {
		return nil, errors.New("mongo handle not found." + bittype)
	}
	return handle, nil
}

func StructToMapDemo(obj interface{}) map[string]interface{} {
	obj1 := reflect.TypeOf(obj)
	obj2 := reflect.ValueOf(obj)
	var data = make(map[string]interface{})
	for i := 0; i < obj1.NumField(); i++ {
		data[obj1.Field(i).Name] = obj2.Field(i).Interface()
	}
	return data
}

//addrStr: '1AiBYt8XbsdyPAELFpcSwRpu45eb2bArMf',
//  balance: 5000000000,
//  totalReceived: 5000000000,
//  totalSent: 0,
//  txApperances: 1,
//  unconfirmedBalance: 0,
//  unconfirmedTxApperances: 0

func ListUnspentInfo(rsp http.ResponseWriter, req *http.Request) {
	errcode := ListUnspentInfo_ErrorCode
	vars := mux.Vars(req)
	mh, err := gHandle.MongoHandle(vars["type"])
	if err != nil {
		rsp.Write(common.MakeResponseWithErr(errcode-1, err))
		return
	}
	fragment := common.StartTrace("ListUnspentInfo db1")
	rst, err := mh.GetAddrUnspentInfo(vars["address"])
	fragment.StopTrace(0)
	if err != nil {
		rsp.Write(common.MakeResponseWithErr(errcode-2, err))
		return
	}

	rsp.Write(common.MakeOkRspByData(rst))
}

func ListUnspentUtxo(rsp http.ResponseWriter, req *http.Request) {
	//vars:= mux.Vars(req)
	//addressType := vars["type"]
	//address := vars["address"]

	utxoList := make([]btcjson.ListUnspentResult, 0, 100)
	rsp.Write(common.MakeOkRspByData(utxoList))
}

func GetChainDifficulty(rsp http.ResponseWriter, req *http.Request) {
	errcode := GetChainDifficulty_ErrorCode
	vars := mux.Vars(req)
	addressType := vars["type"]
	diff, err := gHandle.fullnode.GetFullnode(addressType).GetDifficulty()
	if err != nil {
		rsp.Write(common.MakeResponseWithErr(errcode, err))
		return
	}

	rsp.Write(common.MakeOkRspByData(diff))
}

func GetChainState(rsp http.ResponseWriter, req *http.Request) {
	errcode := GetChainState_ErrorCode
	vars := mux.Vars(req)
	addressType := vars["type"]
	fragment := common.StartTrace("GetChainState db1")
	state, err := gHandle.mongo.Servicehub[addressType].GetChainState()
	fragment.StopTrace(0)
	if err != nil {
		rsp.Write(common.MakeResponseWithErr(errcode, err))
		return
	}

	rsp.Write(common.MakeOkRspByData(state))
}

func AddMonitoredAddr(rsp http.ResponseWriter, req *http.Request) {
	errcode := GetChangedMonitoredAddr_ErrorCode
	vars := mux.Vars(req)
	addressType := vars["type"]
	appid := vars["appid"]
	owneridstr := vars["ownerid"]
	addr := vars["addr"]
	ownerid, err := strconv.Atoi(owneridstr)
	if err != nil {
		rsp.Write(common.MakeResponseWithErr(errcode, err))
		return
	}
	go NotifyMonitoredAddr(addressType, appid, ownerid, addr)
	fragment := common.StartTrace("AddMonitoredAddr db1")
	err = gHandle.mongo.Servicehub[addressType].AddMonitoredAddr(appid, ownerid, addr)
	fragment.StopTrace(0)
	if err != nil {
		rsp.Write(common.MakeResponseWithErr(errcode, err))
		return
	}
	rsp.Write(common.MakeOkRspByData("ok"))
}

func SyncMonitoredAddr(rsp http.ResponseWriter, req *http.Request) {
	//这个函数同AddMonitoredAddr的区别就是他不通知
	errcode := GetChangedMonitoredAddr_ErrorCode
	vars := mux.Vars(req)
	addressType := vars["type"]
	appid := vars["appid"]
	owneridstr := vars["ownerid"]
	addr := vars["addr"]
	ownerid, err := strconv.Atoi(owneridstr)
	if err != nil {
		rsp.Write(common.MakeResponseWithErr(errcode, err))
		return
	}
	fragment := common.StartTrace("AddMonitoredAddr db1")
	err = gHandle.mongo.Servicehub[addressType].AddMonitoredAddr(appid, ownerid, addr)
	fragment.StopTrace(0)
	if err != nil {
		rsp.Write(common.MakeResponseWithErr(errcode, err))
		return
	}
	rsp.Write(common.MakeOkRspByData("ok"))
}

func GetMonitoredAddrOwner(rsp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	addressType := vars["type"]
	addr := vars["addr"]

	fragment := common.StartTrace("GetMonitoredAddrOwner db1")
	ownerid := gHandle.mongo.Servicehub[addressType].GetMonitoredAddr(addr)
	fragment.StopTrace(0)

	rsp.Write(common.MakeOkRspByData(ownerid))
}

func GetBalance(rsp http.ResponseWriter, req *http.Request) {
	errcode := GetBalance_ErrorCode
	vars := mux.Vars(req)
	addressType := vars["type"]
	appid := vars["appid"]
	owneridstr := vars["ownerid"]
	ownerid, err := strconv.Atoi(owneridstr)
	if err != nil {
		rsp.Write(common.MakeResponseWithErr(errcode, err))
		return
	}

	fragment := common.StartTrace("GetBalance db1")
	balance, err := gHandle.mongo.Servicehub[addressType].GetBalance(appid, ownerid)
	fragment.StopTrace(0)
	if err != nil {
		rsp.Write(common.MakeResponseWithErr(errcode, err))
		return
	}

	rsp.Write(common.MakeOkRspByData(balance))
}

func GetInventory(rsp http.ResponseWriter, req *http.Request) {
	errcode := GetInventory_ErrorCode
	vars := mux.Vars(req)
	addressType := vars["type"]
	appid := vars["appid"]
	owneridstr := vars["ownerid"]
	startstr := vars["start"]
	endstr := vars["end"]

	ownerid, err := strconv.Atoi(owneridstr)
	if err != nil {
		rsp.Write(common.MakeResponseWithErr(errcode, err))
		return
	}

	start, err := strconv.Atoi(startstr)
	if err != nil {
		rsp.Write(common.MakeResponseWithErr(errcode, err))
		return
	}

	end, err := strconv.Atoi(endstr)
	if err != nil {
		rsp.Write(common.MakeResponseWithErr(errcode, err))
		return
	}

	fragment := common.StartTrace("GetInventory db1")
	addrtxs, err := gHandle.mongo.Servicehub[addressType].GetInventory(appid, ownerid, start, end)
	fragment.StopTrace(0)
	if err != nil {
		rsp.Write(common.MakeResponseWithErr(errcode, err))
		return
	}

	rsp.Write(common.MakeOkRspByData(addrtxs))
}

func GetMonitoredAddrUtxo(rsp http.ResponseWriter, req *http.Request) {
	errcode := GetMonitoredAddrUtxo_ErrorCode
	vars := mux.Vars(req)
	addressType := vars["type"]
	appid := vars["appid"]
	owneridstr := vars["ownerid"]
	valuestr := vars["value"]

	ownerid, err := strconv.Atoi(owneridstr)
	if err != nil {
		rsp.Write(common.MakeResponseWithErr(errcode, err))
		return
	}

	value, err := strconv.ParseInt(valuestr, 10, 64)
	if err != nil {
		rsp.Write(common.MakeResponseWithErr(errcode, err))
		return
	}

	fragment := common.StartTrace("GetMonitoredAddrUtxo db1")
	_, addrs, err := gHandle.mongo.Servicehub[addressType].GetMonitoredAddrUtxo(appid, ownerid, value)
	fragment.StopTrace(0)
	if err != nil {
		rsp.Write(common.MakeResponseWithErr(errcode, err))
		return
	}

	rsp.Write(common.MakeOkRspByData(addrs))
}

func OnNewTx(rsp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	txhashstr := vars["txhashstr"]
	index := vars["index"]
	isinput := vars["isinput"]
	addr := vars["addr"]
	value := vars["value"]
	glog.Infof("OnNewTx %s %d %d %s %d", txhashstr, index, isinput, addr, value)

	rsp.Write(common.MakeOkRspByData(""))
}

func OnNewBlock(rsp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	height := vars["height"]
	hash := vars["hash"]
	glog.Infof("OnNewBlock %d %s", height, hash)

	rsp.Write(common.MakeOkRspByData(""))
}

func (p *GlobalHandle) GetTxResult(bitType string, txid string) (*btcjson.TxRawResult, error) {
	cli := gHandle.GetRpcCli(bitType)
	txhash, err := chainhash.NewHashFromStr(txid)
	if err != nil {
		return nil, err
	}
	return cli.GetRawTransactionVerbose(txhash)
}

func GetTxInfo(rsp http.ResponseWriter, req *http.Request) {
	errcode := GetTxInfo_ErrorCode
	vars := mux.Vars(req)
	txresult, err := gHandle.GetTxResult(vars["type"], vars["txid"])
	if err != nil {
		rsp.Write(common.MakeResponseWithErr(errcode, err))
		return
	}
	rsp.Write(common.MakeOkRspByData(txresult))
}

type Vin struct {
	Coinbase  string             `json:"coinbase"`
	Txid      string             `json:"txid"`
	Vout      uint32             `json:"vout"`
	ScriptSig *btcjson.ScriptSig `json:"scriptSig"`
	Sequence  uint32             `json:"sequence"`
	Witness   []string           `json:"txinwitness"`
	Addr      string             `json:"addr,omitempty"`
	Value     int64              `json:"value,omitempty"`
}

type TxRawResult struct {
	Hex           string         `json:"hex"`
	Txid          string         `json:"txid"`
	Hash          string         `json:"hash"`
	Size          int32          `json:"size"`
	Vsize         int32          `json:"vsize"`
	Version       int32          `json:"version"`
	LockTime      uint32         `json:"locktime"`
	Vin           []Vin          `json:"vin"`
	Vout          []btcjson.Vout `json:"vout"`
	BlockHash     string         `json:"blockhash"`
	BlockHeight   int64          `json:"blockheight"`
	Confirmations uint64         `json:"confirmations"`
	Time          int64          `json:"time"`
	Blocktime     int64          `json:"blocktime"`
	Fee           int64          `json:"fee"`
	ValueOut      int64          `json:"valueOut"`
	ValueIn       int64          `json:"valueIn"`
	CoinDay       int64          `json:"coinday"`
}

func CopyFromBtc(dst *btcjson.TxRawResult) (rst TxRawResult) {
	b, err := json.Marshal(dst)
	if err != nil {
		glog.Error(err)
	}
	rst = TxRawResult{}
	err = json.Unmarshal(b, &rst)
	if err != nil {
		glog.Error(err)
	}
	return rst
}

func GetBrowserTxInfo(rsp http.ResponseWriter, req *http.Request) {
	errcode := GetTxInfo_ErrorCode
	vars := mux.Vars(req)
	txBtcResult, err := gHandle.GetTxResult(vars["type"], vars["txid"])
	if err != nil {
		rsp.Write(common.MakeResponseWithErr(errcode, err))
		return
	}
	mh, err := gHandle.MongoHandle(vars["type"])
	if err != nil {
		rsp.Write(common.MakeResponseWithErr(errcode-1, err))
		return
	}
	txresult := CopyFromBtc(txBtcResult)
	totalInput := int64(0)
	for i, vi := range txresult.Vin {
		tmpTx, err := mh.GetTx(vi.Txid)
		if err != nil {
			txresult.Vin[i].Addr = "Get Tx failed!" + vi.Txid
			txresult.Vin[i].Value = 0
		} else {
			if len(tmpTx.VoutAddrValue) <= int(vi.Vout) {
				txresult.Vin[i].Addr = fmt.Sprintf("Get Tx Vout failed!len:%d, Sequence:%d.", len(tmpTx.VoutAddrValue), vi.Vout)
				txresult.Vin[i].Value = 0
			} else {
				tmpOut := tmpTx.VoutAddrValue[vi.Vout]
				list := strings.Split(tmpOut, "@")
				if len(list) != 2 {
					txresult.Vin[i].Addr = tmpOut
				} else {
					txresult.Vin[i].Addr = list[0]
					txresult.Vin[i].Value, err = strconv.ParseInt(list[1], 10, 64)
					if err != nil {
						glog.Error("illegal int64," + list[1])
						txresult.Vin[i].Value = 0
					}
				}
			}
		}
		totalInput += txresult.Vin[i].Value
	}
	totalOutput := int64(0)
	for _, vo := range txresult.Vout {
		amount, _ := btcutil.NewAmount(vo.Value)
		totalOutput += int64(amount)
	}

	txresult.Fee = totalInput - totalOutput
	txresult.ValueIn = totalInput
	txresult.ValueOut = totalOutput
	txinfo, err := mh.GetTx(txresult.Hash)
	txresult.BlockHeight = txinfo.Height
	txresult.CoinDay = txinfo.CoinDay
	blockrst, err := mh.GetBlockInfoByHeight(txinfo.Height)
	if err != nil {
		rsp.Write(common.MakeResponseWithErr(errcode-2, err))
		return
	}
	txresult.BlockHash = blockrst.Hash
	rsp.Write(common.MakeOkRspByData(txresult))
}

func GetBlockChainInfo(rsp http.ResponseWriter, req *http.Request) {
	errcode := GetChainInfo_ErrorCode
	vars := mux.Vars(req)
	cli := gHandle.GetRpcCli(vars["type"])
	blockinfo, err := cli.GetBlockChainInfo()
	if err != nil {
		rsp.Write(common.MakeResponseWithErr(errcode, err))
		return
	}
	rsp.Write(common.MakeOkRspByData(blockinfo))
}

func GetBlockByHeight(rsp http.ResponseWriter, req *http.Request) {
	errcode := GetBlockByHeight_ErrorCode
	vars := mux.Vars(req)
	height, err := strconv.ParseInt(vars["height"], 10, 64)
	if err != nil {
		rsp.Write(common.MakeResponseWithErr(errcode, err))
		return
	}
	mh, err := gHandle.MongoHandle(vars["type"])
	if err != nil {
		rsp.Write(common.MakeResponseWithErr(errcode-1, err))
		return
	}
	blockHeader, err := mh.GetBlockInfoByHeight(height)
	if err != nil {
		rsp.Write(common.MakeResponseWithErr(errcode-2, err))
		return
	}
	data := StructToMapDemo(blockHeader)
	txinfo, err := mh.GetTx(blockHeader.CoinBaseTxid)
	data["reward"] = txinfo.Fee
	data["coinbase"] = string(blockHeader.CoinBaseInfo)
	rsp.Write(common.MakeOkRspByData(data))
}

func GetBlockByHash(rsp http.ResponseWriter, req *http.Request) {
	errcode := GetBlockByHash_ErrorCode
	vars := mux.Vars(req)
	mh, err := gHandle.MongoHandle(vars["type"])
	if err != nil {
		rsp.Write(common.MakeResponseWithErr(errcode-1, err))
		return
	}
	blockHeader, err := mh.GetBlockInfoByHash(vars["hash"])
	if err != nil {
		rsp.Write(common.MakeResponseWithErr(errcode-2, err))
		return
	}
	data := StructToMapDemo(blockHeader)
	txinfo, err := mh.GetTx(blockHeader.CoinBaseTxid)
	data["reward"] = txinfo.Fee
	data["coinbase"] = string(blockHeader.CoinBaseInfo)
	rsp.Write(common.MakeOkRspByData(data))
}

func GetBestBlock(rsp http.ResponseWriter, req *http.Request) {
	errcode := GetBestBlock_ErrorCode
	vars := mux.Vars(req)
	start, err := strconv.Atoi(vars["limit"])
	if err != nil {
		rsp.Write(common.MakeResponseWithErr(errcode, err))
		return
	}
	limit, err := strconv.Atoi(vars["limit"])
	if err != nil {
		rsp.Write(common.MakeResponseWithErr(errcode-1, err))
		return
	}
	mh, err := gHandle.MongoHandle(vars["type"])
	if err != nil {
		rsp.Write(common.MakeResponseWithErr(errcode-2, err))
		return
	}
	items, err := mh.GetBestBlockHeader(start, limit)
	if err != nil {
		rsp.Write(common.MakeResponseWithErr(errcode-3, err))
		return
	}
	rsp.Write(common.MakeOkRspByData(items))
}

type AddrMount struct {
	Address string `json:"address"`
	Value   int64  `json:"value"`
}

type TxBrowser struct {
	Hash          string      `json:"hash"`
	Height        int64       `json:"height"`
	Vin           []AddrMount `json:"vin"`
	Vout          []AddrMount `json:"vout"`
	Fee           int64       `json:"fee"`
	CoinDay       int64       `json:"coinday"`
	Confirmations int64       `json:"confirmations"`
	Time          int64       `json:"time"`
}

type BlockBrowserTxs struct {
	CoinbaseTxid string      `json:"CoinbaseTxid"`
	Coinbase     []byte      `json:"Coinbase"`
	Txs          []TxBrowser `json:"txs"`
}

func MongoTxFormat2BrowserByOne(item *common.TxBson) TxBrowser {
	rt := TxBrowser{item.Hash, item.Height,
		make([]AddrMount, len(item.VinAddrValue)),
		make([]AddrMount, len(item.VoutAddrValue)),
		item.Fee, item.CoinDay, item.Confirmations, item.Time}
	for i, it := range item.VinAddrValue {
		list := strings.Split(it, "@")
		if len(list) != 2 {
			rt.Vin[i] = AddrMount{"coinbase", 0}
		} else {
			amount, _ := strconv.ParseInt(list[1], 10, 64)
			rt.Vin[i] = AddrMount{list[0], amount}
		}
	}

	for i, it := range item.VoutAddrValue {
		list := strings.Split(it, "@")
		if len(list) != 2 {
			rt.Vout[i] = AddrMount{"coinbase", 0}
		} else {
			amount, _ := strconv.ParseInt(list[1], 10, 64)
			rt.Vout[i] = AddrMount{list[0], amount}
		}
	}
	return rt
}

func MongoTxFormat2Browser(items []*common.TxBson) []TxBrowser {
	results := make([]TxBrowser, len(items), len(items))
	for k, item := range items {
		results[k] = MongoTxFormat2BrowserByOne(item)
	}
	return results
}

func GetBlockBrowserTxsByHash(rsp http.ResponseWriter, req *http.Request) {
	errcode := GetBlockBrowserTxsv_ErrorCode
	vars := mux.Vars(req)
	start, err := strconv.Atoi(vars["start"])
	if err != nil {
		rsp.Write(common.MakeResponseWithErr(errcode-1, err))
		return
	}
	end, err := strconv.Atoi(vars["end"])
	if err != nil {
		rsp.Write(common.MakeResponseWithErr(errcode-2, err))
		return
	}
	mh, err := gHandle.MongoHandle(vars["type"])
	if err != nil {
		rsp.Write(common.MakeResponseWithErr(errcode-3, err))
		return
	}
	blockinfo, err := mh.GetBlockInfoByHash(vars["hash"])
	if err != nil {
		rsp.Write(common.MakeResponseWithErr(errcode-4, err))
		return
	}

	items, err := mh.GetBlockTxs(blockinfo.Hash, start, end)
	if err != nil {
		rsp.Write(common.MakeResponseWithErr(errcode-5, err))
		return
	}
	fr := BlockBrowserTxs{blockinfo.CoinBaseTxid, blockinfo.CoinBaseInfo, MongoTxFormat2Browser(items)}
	rsp.Write(common.MakeOkRspByData(fr))
}

func GetAddressBrowserTxs(rsp http.ResponseWriter, req *http.Request) {
	errcode := GetAddressBrowserTxs_ErrorCode
	vars := mux.Vars(req)
	start, err := strconv.Atoi(vars["start"])
	if err != nil {
		rsp.Write(common.MakeResponseWithErr(errcode-1, err))
		return
	}
	end, err := strconv.Atoi(vars["end"])
	if err != nil {
		rsp.Write(common.MakeResponseWithErr(errcode-2, err))
		return
	}
	mh, err := gHandle.MongoHandle(vars["type"])
	if err != nil {
		rsp.Write(common.MakeResponseWithErr(errcode-3, err))
		return
	}

	items, err := mh.GetAddressTxs(vars["address"], start, end)
	if err != nil {
		rsp.Write(common.MakeResponseWithErr(errcode-4, err))
		return
	}
	rsp.Write(common.MakeOkRspByData(MongoTxFormat2Browser(items)))
}

func NotifyMonitoredAddr(coinType string, appid string, ownerid int, addr string) {
	url := fmt.Sprintf("http://%s/api/monitored_address/add/type/%s/appid/%s/ownerid/%d/addr/%s",
		gHandle.conf.SyncCenterIp, coinType, appid, ownerid, addr)
	_, err := http.Get(url)
	if err != nil {
		glog.Infof("NotifyMonitoredAddr %s fault", addr)
	}
}

func CheckOwneridBalance(rsp http.ResponseWriter, req *http.Request) {
	errcode := GetOwnerBalance_ErrorCode
	vars := mux.Vars(req)
	coinType := vars["type"]
	hub, ok := gHandle.mongo.Servicehub[coinType]
	if !ok {
		rsp.Write(common.MakeResponseWithErr(errcode, errors.New("coin type not sopport")))
		return
	}
	appid := vars["appid"]
	owneridStr := vars["ownerid"]
	ownerid, err := strconv.Atoi(owneridStr)
	if err != nil {
		rsp.Write(common.MakeResponseWithErr(errcode, err))
		return
	}
	go func() {
		balancePairs, count := hub.CheckOwnerBalance(appid, ownerid)
		result := make(map[string]interface{})
		result["bitdbIp"] = glbBitdbIp
		result["total"] = count
		result["different"] = len(balancePairs)
		result["balancePairs"] = balancePairs
		result["ownerid"] = ownerid
		contentByte, err := json.Marshal(result)
		if err != nil {
			glog.Infof("CheckBalance Marshal:%s", err)
		}
		err = go_util.SmsNotify(go_util.Sms_bitdb_incorrect, string(contentByte))
		if err != nil {
			glog.Infof("CheckBalance SmsNotify:%s", err)
		}
	}()
	rsp.Write(common.MakeOkRspByData("OK"))
}

func GetOwneridCount(rsp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	coinType := vars["type"]
	appid := vars["appid"]
	count := gHandle.mongo.Servicehub[coinType].CountOwnerNubmer(appid)
	rsp.Write(common.MakeOkRspByData(count))
}
