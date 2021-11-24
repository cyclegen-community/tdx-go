package v1

import (
	"github.com/cyclegen-community/tdx-go/proto"
	"github.com/cyclegen-community/tdx-go/utils"
)

// 请求包结构
type SetupCmd1Request struct {
	Cmd []byte `struc:"[13]byte" json:"cmd"`
}

// 请求包序列化输出
func (req *SetupCmd1Request) Marshal() ([]byte, error) {
	return proto.DefaultMarshal(req)
}

// 响应包结构
type SetupCmd1Response struct {
	Unknown []byte `json:"unknown"`
}

func (resp *SetupCmd1Response) Unmarshal(data []byte) error {
	resp.Unknown = data
	return nil
}

// 创建SetupCmd1请求包
func NewSetupCmd1Request() (*SetupCmd1Request, error) {
	request := &SetupCmd1Request{
		Cmd: utils.HexString2Bytes("0c 02 18 93 00 01 03 00 03 00 0d 00 01"),
	}
	return request, nil
}

func NewSetupCmd1() (*SetupCmd1Request, *SetupCmd1Response, error) {
	var response SetupCmd1Response
	var request, err = NewSetupCmd1Request()
	return request, &response, err
}
