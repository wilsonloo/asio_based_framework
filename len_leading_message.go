package wgnet

/**
作者:wilsonloo
模块：tcp连接的 消息
说明：消息采用 消息头 + 数据 的组织方式，其中消息头是一个数组，其长度可自定义 MESSAGE_HEADER_LEN
创建时间：2016-5-29
**/

import (
	"fmt"
	proto "github.com/golang/protobuf/proto"
)

// 定义
const (
	PACKET_HEADER_LEN   = 4    // 消息头的长度
	MAX_PACKET_DATA_LEN = 1024 // 最大消息长度
	PACKET_MASK = 0x0001FFFF // 消息掩码
	FLAGS_MASK = 0xFFFE0000 // 控制标记掩码
)

// todo 此部分可修改
// 消息头声明（注：该结构体只是用以说明，不会被使用）
type _PacketHeader struct {
	len uint16 // 长度
	cmd uint16 // 命令号
}

// todo 此部分可修改
// 消息头声明（注：该结构体只是用以说明，不会被使用）
type _PacketData struct {
	holder uint32 // 占位符
}

// 消息定义
type LenLeadingMessage struct {
	Header []byte // 消息头，
	Data   []byte // 实际消息
}

// 获取整体消息大小
func (msg *LenLeadingMessage) MessageTotalSize() uint32{
	return uint32(PACKET_HEADER_LEN) + uint32(msg.PacketLen())
}

// 添加标记
func (msg *LenLeadingMessage) AddFlag(flags uint32)  {
	new_flags := msg.GetFlags()
	new_flags = new_flags | flags

	msg.Header[3] = byte(new_flags >> 8)
	msg.Header[2] = byte(new_flags)
}

// 获取消息头信息
func (msg *LenLeadingMessage)GetHeaderInfo() ( /*data*/ []byte, /*data len*/ uint32) {
	return msg.Header, PACKET_HEADER_LEN
}

// 获取消息体信息
func (msg *LenLeadingMessage) GetBodyInfo()  ( /*data*/ []byte, /*data len*/ uint32) {
	return msg.Data, uint32(msg.PacketLen())
}

// todo 此部分可修改
func GetPacketLen(header []byte) uint16 {
	val := GetUint32(header)
	return uint16(val & PACKET_MASK)
}

// 获取消息长度
func (msg *LenLeadingMessage) PacketLen() uint16 {
	return GetPacketLen(msg.Header)
}

// 获取控制标记
func (msg *LenLeadingMessage) GetFlags() uint32 {
	val := GetUint32(msg.Header)
	return uint32(val & FLAGS_MASK)
}

// todo 此部分可修改
// 设置消息长度
func (msg *LenLeadingMessage) SetPacketLen(len uint16) {
	SetUint16(msg.Header[0:], len)
}

// todo 此部分可修改
// 重置消息
func (msg *LenLeadingMessage) ResetPacket() {
	// todo 需要放回池内
	msg.Data = nil

	msg.Header[0] = 0
	msg.Header[1] = 0
}

func GetUint32(buf []byte) uint32 {
	val := uint32(buf[3]) << 24 |
		uint32(buf[2]) << 16 |
		uint32(buf[1]) << 8 |
		uint32(buf[0])
fmt.Println("GetUint32:", buf, val)
	return val
}

func SetUint32(buf []byte, val uint32) {
	buf[3] = byte(val >> 24)
	buf[2] = byte(val >> 16)
	buf[1] = byte(val >> 8)
	buf[0] = byte(val)
}

func SetUint16(buf []byte, val uint16) {
	buf[1] = byte(val >> 8)
	buf[0] = byte(val)
}

func SetUint8(buf []byte, val uint8) {
	buf[0] = byte(val)
}

func (msg *LenLeadingMessage) InitData() {
	// todo 需要从池内获取
	msg.Data = make([]byte, msg.PacketLen())
}

func(msg *LenLeadingMessage) Dump() []byte {
	ret := make([]byte, msg.MessageTotalSize())
	copy(ret[0:], msg.Header[:])
	copy(ret[PACKET_HEADER_LEN:], msg.Data[:])

	return ret
}

//UnpackagePbmsg 解包protobuf消息
func (msg *LenLeadingMessage) Unpackage2Pbmsg(pb proto.Message) error {
	return proto.Unmarshal(msg.Data, pb)
}

//Package 打包原生字符串
// @param raw_len 返回原始数据的长度
func (msg *LenLeadingMessage) Package(cmd uint16, buff []byte) (packeted_len uint32, err error) {
	size := len(buff)
	if size == 0 {
		return 0, nil
	}

	msg.ResetPacket()
	msg.SetPacketLen(uint16(size))
	msg.InitData()

	// 先写
	copy(msg.Data[:], buff)
	return uint32(size), nil
}

// 按照消息长度初始化 消息体
func (msg *LenLeadingMessage) PreparePacket() {
	// todo 优化到从缓冲池读取数据
	msg.Data = make([]byte, msg.PacketLen())
}

func MakeHeader() []byte {
	return make([]byte, PACKET_HEADER_LEN)
}

/* 创建一个消息
@param data_size 预设的消息长度
*/
func NewLenLeadingMessage() *LenLeadingMessage {
	// todo 优化：采用 message pool 的方式

	msg := new(LenLeadingMessage)

	// 固定消息头长度
	msg.Header = MakeHeader()

	return msg
}

func FreeLenLeadingMessage(msg *LenLeadingMessage) {
	// todo 进行回收到消息池
}