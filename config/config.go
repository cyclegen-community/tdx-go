package config

import (
	"encoding/json"
	"github.com/cyclegen-community/tdx-go/utils"
	"github.com/sparrc/go-ping"
	"io/ioutil"
	"log"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Server struct {
	Name string `json:"name"`
	IP   string `json:"ip"`
	Port int    `json:"port"`
}

func (srv *Server) Addr() string {
	return strings.Join([]string{srv.IP, strconv.Itoa(srv.Port)}, ":")
}

const (
	_StockQuotesServerConfigFile = "stock_ip.json"
)

// StockQuotesServer 股票行情线路信息
type StockQuotesServer []Server

// GetStockQuotesServer 获取股票行情线路列表
func GetStockQuotesServer() StockQuotesServer {
	var instance StockQuotesServer
	raw, err := ioutil.ReadFile("config/" + _StockQuotesServerConfigFile)
	if err != nil {
		log.Fatalf("读取"+_StockQuotesServerConfigFile+"失败, 错误详情: %v", err.Error())
	}
	err = json.Unmarshal(raw, &instance)
	if err != nil {
		log.Fatalf("解析"+_StockQuotesServerConfigFile+"失败, 错误详情: %v", err.Error())
	}
	return instance
}

// GetBestStockQuotesServer 获取最优股票行情线路
func GetBestStockQuotesServer() Server {
	srvs := GetStockQuotesServer()
	results := sync.Map{}
	sortableSrvs := utils.SortableMapList{}
	wg := sync.WaitGroup{}
	for idx := range srvs {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			if srvs[id].IP == "" {
				return
			}
			pinger, err := ping.NewPinger(srvs[id].IP)
			if err != nil {
				log.Println(err)
				return
			}
			pinger.SetPrivileged(true)
			pinger.Count = 5
			pinger.Timeout = time.Second
			pinger.Run() // blocks until finished

			stats := pinger.Statistics() // get send/receive/rtt stats

			var avgRtt int64
			// 丢包率不能高于50%，否则设置平均时延为负数
			if stats.PacketLoss > 0.5 {
				avgRtt = -1
			} else {
				avgRtt = stats.AvgRtt.Nanoseconds()
			}
			results.Store(srvs[id], avgRtt)
		}(idx)
	}
	wg.Wait()
	results.Range(func(key, value interface{}) bool {
		srv := key.(Server)
		avgRtt := value.(int64)
		//log.Printf("%v: %d ns", srv.Addr(), avgRtt)
		if avgRtt > 0 {
			sortableSrvs = append(sortableSrvs, utils.SortableMap{
				srv,
				avgRtt,
			})
		}
		return true
	})

	sort.Sort(sortableSrvs)
	if len(sortableSrvs) > 0 {
		srv := sortableSrvs[0].Key.(Server)
		return srv
	} else {
		panic("所有服务器均无法连通！")
	}
}
