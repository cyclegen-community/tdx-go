package v1

import (
	"github.com/cyclegen-community/tdx-go/proto"
	"github.com/cyclegen-community/tdx-go/utils"
)

type Market int

const (
	MarketShenZhen Market = iota
	MarketShangHai        = 1
)

// 请求包结构
type GetSecurityCountRequest struct {
	Unknown1 []byte `struc:"[12]byte"`
	Market   Market `struc:"uint16,little";json:"market"`
	Unknown2 []byte `struc:"[4]byte"`
}

// 请求包序列化输出
func (req *GetSecurityCountRequest) Marshal() ([]byte, error) {
	return proto.DefaultMarshal(req)
}

// 响应包结构
type GetSecurityCountResponse struct {
	Count uint `struc:"uint16,little";json:"count"`
}

func (resp *GetSecurityCountResponse) Unmarshal(data []byte) error {
	return proto.DefaultUnmarshal(data, resp)
}

// todo: 检测market是否为合法值
func NewGetSecurityCountRequest(market Market) (*GetSecurityCountRequest, error) {
	request := &GetSecurityCountRequest{
		Unknown1: utils.HexString2Bytes("0c 0c 18 6c 00 01 08 00 08 00 4e 04"),
		Market:   market,
		Unknown2: utils.HexString2Bytes("75 c7 33 01"),
	}
	return request, nil
}

func NewGetSecurityCount(market Market) (*GetSecurityCountRequest, *GetSecurityCountResponse, error) {
	var response GetSecurityCountResponse
	var request, err = NewGetSecurityCountRequest(market)
	return request, &response, err
}
