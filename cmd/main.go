package main

import (
	"github.com/cyclegen/tdx-go/config"
	"github.com/cyclegen/tdx-go/core"
	"github.com/cyclegen/tdx-go/proto"
	"github.com/cyclegen/tdx-go/proto/v1"
	"log"
)

func init() {
	log.SetFlags(log.Lshortfile | log.Ldate)
}
func main() {
	quotesSrv := config.GetBestStockQuotesServer()
	//quotesSrvAddr := "106.120.74.86:7711" // quotesSrv.Addr()
	log.Println("正在连接到最优行情服务器: ", quotesSrv.Addr())
	T(quotesSrv.IP, quotesSrv.Port)
	//T("106.120.74.86", 7709)
}
func T(ip string, port int) {
	cli := core.NewClient(ip, port)

	// CMD信令 1
	testProto(cli, func() (req proto.Marshaler, resp proto.Unmarshaler, err error) {
		req, resp, err = v1.NewSetupCmd1()
		return
	})
	// CMD信令 2
	testProto(cli, func() (req proto.Marshaler, resp proto.Unmarshaler, err error) {
		req, resp, err = v1.NewSetupCmd2()
		return
	})
	// CMD信令 3
	testProto(cli, func() (req proto.Marshaler, resp proto.Unmarshaler, err error) {
		req, resp, err = v1.NewSetupCmd3()
		return
	})
	// 查询股票数量
	testProto(cli, func() (req proto.Marshaler, resp proto.Unmarshaler, err error) {
		req, resp, err = v1.NewGetSecurityCount(v1.MarketShangHai)
		return
	})

	testProto(cli, func() (req proto.Marshaler, resp proto.Unmarshaler, err error) {
		req, resp, err = v1.NewGetSecurityList(v1.MarketShangHai, 255)
		return
	})
	//testProto(cli, func() (req proto.Marshaler, resp proto.Unmarshaler, err error) {
	//	req, resp, err = v1.NewGetSecurityQuotes()
	//	return
	//})
}

func testProto(cli *core.Client, factory proto.Factory) {
	req, resp, err := factory()
	if err != nil {
		log.Fatal(err)
	}
	err = cli.Do(req, resp)
	if err != nil {
		log.Println(err)
	}
	log.Println(resp)
}
