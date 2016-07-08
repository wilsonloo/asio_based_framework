package wgnet

////////////////////////////////////////////////////
// Time        : 2016/6/17 11:14
// Author      : wilsonloo21@163.com
// File        : protocol_adapter_len_leading.go
// Software    : PyCharm
// Description : 长度引导的协议适配器
////////////////////////////////////////////////////

import (
	"asio"
	"fmt"
	"errors"
)

type LenLeadingProtocolProcessor struct {
}

// 进行读取数据
func (this *LenLeadingProtocolProcessor)HandleRecv(session asio.Session)(asio.Message, error) {

	//写消息头
	//如果读取消息失败，消息要归还给消息池
	msg := NewLenLeadingMessage()
	// read_len, err := session.Conn.Read(msg.Header)
	read_len, err := session.ReadLenFixedData(msg.Header, uint32(len(msg.Header)))
	if err != nil {
		FreeLenLeadingMessage(msg)
		session.SetConnected(false)

		fmt.Printf("[WgNet] failed to Recv data, session marked Disconnected, with error: %s \n", err.Error())
		return nil, err
	}
	fmt.Println("ReadLenFixedData:", msg.Dump())

	if uint16(read_len) != PACKET_HEADER_LEN {
		FreeLenLeadingMessage(msg)
		return nil, errors.New("recv error")
	}

	/*
		if err = msg.CheckFormat(); err != nil {
			GxMisc.Error("XXXX %s remote[%s:%s] format err: %d", GxStatic.CmdString[msg.GetCmd()], session.M, session.Remote, msg.GetLen())
			GxMessage.FreeMessage(msg)
			return nil, err
		}*/

	// 获取消息数据的长度
	packet_len := msg.PacketLen()
	if packet_len == 0 {
		return msg, nil
	}
	fmt.Println("packet len is ", packet_len)

	if packet_len > MAX_PACKET_DATA_LEN {
		FreeLenLeadingMessage(msg)
		return nil, errors.New("packet length error.")
	}

	// 创建消息体
	msg.PreparePacket()

	// 阻塞式写满packet数据
	read_len, err = session.ReadLenFixedData(msg.Data[0:], uint32(len(msg.Data)))
	fmt.Println("ReadLenFixedData data:", msg.Dump())

	// 检测错误
	if err != nil {
		/* if err != io.EOF {
			return nil, err
		}*/
		FreeLenLeadingMessage(msg)

		fmt.Println("failed to read data:", err.Error())
		return nil, err
	}

	// 必须整整一个消息
	if read_len != uint32(msg.PacketLen()) {
		fmt.Printf("failed to read data, expecting %d got %d", msg.PacketLen(), read_len)
		FreeLenLeadingMessage(msg)
		return nil, errors.New("packet len reading error.")
	}

	return msg, nil
}

// 进行读取数据
func (this *LenLeadingProtocolProcessor)HandleSend(session asio.Session, msg asio.Message) error {
	// 检测连接状态
	if !session.IsConnected() {
		// todo log
		return errors.New(fmt.Sprintf("remote disconnect,send msg fail"))
	}

	ll_msg := msg.(*LenLeadingMessage)

	var err error
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