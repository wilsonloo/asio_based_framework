package wgnet

////////////////////////////////////////////////////
// Time        : 2016/7/4 17:40
// Author      : wilsonloo21@163.com
// File        : packet_base.go
// Software    : PyCharm
// Description : 消息基类
////////////////////////////////////////////////////

type PacketBase struct {
	BaseFlag uint16 // 0xA839
	BaseSum uint16 // 校验码
	BaseIndex uint16 // 编号
	First uint8 // 一级指令
	Second uint8 // 二级指令
}