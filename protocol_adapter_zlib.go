package wgnet

////////////////////////////////////////////////////
// Time        : 2016/7/4 15:59
// Author      : wilsonloo21@163.com
// File        : protocol_adapter_zlib.go
// Software    : PyCharm
// Description : 使用zlib加密的协议解析适配器
////////////////////////////////////////////////////

import (
	"asio"
	"fmt"
	"errors"
)

type ZlibProtocolProcessor struct {
}

// 进行读取数据
func (this *ZlibProtocolProcessor)HandleRecv(session asio.Session)(asio.Message, error) {
	// 消息头
	// 如果读取消息失败，消息要归还给消息池
	msg := NewLenLeadingMessage()
	read_len, err := session.ReadLenFixedData(msg.Header, uint32(len(msg.Header)))
	if err != nil {
		FreeLenLeadingMessage(msg)
		session.SetConnected(false)

		fmt.Printf("[WgNet] failed to Recv data, session marked Disconnected, with error: %s\n", err.Error());
		return nil, err
	}

	if uint16(read_len) != PACKET_HEADER_LEN {
		FreeLenLeadingMessage(msg);
		return nil, errors.New("recv error")
	}

	// 获取消息数据的长度
	packet_len := msg.PacketLen()
	if packet_len == 0 {
		return msg, nil
	}

	if packet_len > MAX_PACKET_DATA_LEN {
		FreeLenLeadingMessage(msg)
		return nil, errors.New("packet length excceed limit.")
	}

	// 创建消息体
	msg.PreparePacket()

	// 阻塞式写满packet数据
	read_len, err = session.ReadLenFixedData(msg.Data[0:], uint32(len(msg.Data)))

	// 检测错误
	if err != nil {
		FreeLenLeadingMessage(msg);
		return nil, err
	}

	// 必须整整一个消息
	if read_len != uint32(msg.PacketLen()) {
		FreeLenLeadingMessage(msg)
		return nil, errors.New("packet len reading error.")
	}

	// todo 进行zip解压

	return msg, nil
}

// 进行读取数据
func (this *ZlibProtocolProcessor)HandleSend(session asio.Session, msg asio.Message) error {
	// 检测连接状态
	if !session.IsConnected() {
		// todo log
		return errors.New(fmt.Sprintf("remote disconnect,send msg fail"))
	}

	var err error

	ll_msg := msg.(*LenLeadingMessage)

	// todo 将源消息进行压缩
	// zip_msg := NewLenLeadingMessage()


	// todo 设置标记，表明已使用zip压缩
	if false {
		ll_msg.AddFlag(0x4000000)
	}

	// 先发送消息头
	if _, err = session.Write(ll_msg.Header); err != nil {
		return err
	}

	// 再发送消息体
	if _, err = session.Write(ll_msg.Data); err != nil {
		return err
	}

	return nil
}