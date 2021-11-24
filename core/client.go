package core

import (
	"errors"
	"github.com/cyclegen-community/tdx-go/proto"
	"github.com/cyclegen-community/tdx-go/utils"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

type Client struct {
	conn          net.Conn
	Host          string
	Port          int
	Timeout       time.Duration
	MaxRetryTimes int
	RetryDuration time.Duration
}

// https://www.sohamkamani.com/golang/options-pattern/
// https://lingchao.xin/post/functional-options-pattern-in-go.html
// NewBaseClient 创建BaseClient实例
func NewClient(host string, port int) *Client {
	addr := strings.Join([]string{host, strconv.Itoa(port)}, ":")

	conn, err := net.Dial("tcp", addr) // net.DialTimeout()

	if err != nil {
		log.Fatalln(err)
	}
	return &Client{
		conn:          conn,
		Host:          host,
		Port:          port,
		MaxRetryTimes: 5,
		Timeout:       time.Second,
		RetryDuration: time.Millisecond * 200,
	}
}

func (cli *Client) Do(request proto.Marshaler, response proto.Unmarshaler) error {
	// 序列化请求
	req, err := request.Marshal()
	if err != nil {
		return err
	}
	// 发送请求
	retryTimes := 0
SEND:
	n, err := cli.conn.Write(req)
	// 重试
	if n < len(req) {
		retryTimes += 1
		if retryTimes <= cli.MaxRetryTimes {
			log.Printf("第%d次重试\n", retryTimes)
			goto SEND
		} else {
			return errors.New("数据未完整发送")
		}
	}
	if err != nil {
		return err
	}
	// 解析响应包头
	var header proto.PacketHeader
	// 读取包头 大小为16字节
	// 单次获取的字列流
	headerLength := 0x10
	headerBytes := make([]byte, headerLength)
	// 调用socket获取字节流并保存到data中
	headerBytes, err = cli.receive(headerLength)
	if err != nil {
		//log.Println(err)
		return err
	}
	err = header.Unmarshal(headerBytes)
	if err != nil {
		return err
	}
	// 根据获取响应体结构
	// 调用socket获取字节流并保存到data中
	bodyBytes, err := cli.receive(header.ZipSize)
	if err != nil {
		return err
	}
	// zlib解压缩
	if header.Compressed() {
		bodyBytes, err = utils.ZlibUnCompress(bodyBytes)
	}
	// 反序列化为响应体结构
	err = response.Unmarshal(bodyBytes)
	if err != nil {
		return err
	}
	return nil
}
func (cli *Client) receive(length int) (data []byte, err error) {
	var (
		receivedSize int
	)
READ:
	tmp := make([]byte, length)
	// 调用socket获取字节流并保存到data中
	receivedSize, err = cli.conn.Read(tmp)
	// socket错误,可能为EOF
	if err != nil {
		return nil, err
	}
	// 数据添加到总输出,由于tmp申请内存时使用了length的长度，
	// 所以直接全部复制到data中会使得未完全传输的部分被填充为0导致数据获取不完整，
	// 故使用tmp[:receivedSize]
	data = append(data, tmp[:receivedSize]...)
	// 数据读满就可以返回了
	if len(data) == length {
		return
	}
	// 读取小于标准尺寸，说明到文件尾或者读取出现了问题没读满，可以返回了
	if receivedSize < length {
		goto READ
	}
	return
}
func (cli *Client) Close() error {
	return cli.conn.Close()
}
